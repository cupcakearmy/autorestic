package internal

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/cupcakearmy/autorestic/internal/colors"
	"github.com/cupcakearmy/autorestic/internal/flags"
)

type BackendRest struct {
	User     string `mapstructure:"user,omitempty" yaml:"user,omitempty"`
	Password string `mapstructure:"password,omitempty" yaml:"password,omitempty"`
}

type Backend struct {
	name       string
	Type       string            `mapstructure:"type,omitempty" yaml:"type,omitempty"`
	Path       string            `mapstructure:"path,omitempty" yaml:"path,omitempty"`
	Key        string            `mapstructure:"key,omitempty" yaml:"key,omitempty"`
	RequireKey bool              `mapstructure:"requireKey,omitempty" yaml:"requireKey,omitempty"`
	Init       bool              `mapstructure:"init,omitempty" yaml:"init,omitempty"`
	Env        map[string]string `mapstructure:"env,omitempty" yaml:"env,omitempty"`
	Rest       BackendRest       `mapstructure:"rest,omitempty" yaml:"rest,omitempty"`
	Options    Options           `mapstructure:"options,omitempty" yaml:"options,omitempty"`
}

func GetBackend(name string) (Backend, bool) {
	b, ok := GetConfig().Backends[name]
	b.name = name
	return b, ok
}

func (b Backend) generateRepo() (string, error) {
	switch b.Type {
	case "local":
		return GetPathRelativeToConfig(b.Path)
	case "rest":
		parsed, err := url.Parse(os.ExpandEnv(b.Path))
		if err != nil {
			return "", err
		}
		if b.Rest.User != "" {
			if b.Rest.Password == "" {
				parsed.User = url.User(b.Rest.User)
			} else {
				parsed.User = url.UserPassword(b.Rest.User, b.Rest.Password)
			}
		}
		return fmt.Sprintf("%s:%s", b.Type, parsed.String()), nil
	case "b2", "azure", "gs", "s3", "sftp", "rclone":
		return fmt.Sprintf("%s:%s", b.Type, b.Path), nil
	default:
		return "", fmt.Errorf("backend type \"%s\" is invalid", b.Type)
	}
}

var nonAlphaRegex = regexp.MustCompile("[^A-Za-z0-9]")

func (b Backend) getEnv() (map[string]string, error) {
	env := make(map[string]string)
	// Key
	if b.Key != "" {
		env["RESTIC_PASSWORD"] = b.Key
	}

	// From config file
	repo, err := b.generateRepo()
	env["RESTIC_REPOSITORY"] = repo
	for key, value := range b.Env {
		env[strings.ToUpper(key)] = value
	}

	// From Envfile and passed as env
	nameForEnv := strings.ToUpper(b.name)
	nameForEnv = nonAlphaRegex.ReplaceAllString(nameForEnv, "_")
	var prefix = "AUTORESTIC_" + nameForEnv + "_"
	for _, variable := range os.Environ() {
		var splitted = strings.SplitN(variable, "=", 2)
		if strings.HasPrefix(splitted[0], prefix) {
			env[strings.TrimPrefix(splitted[0], prefix)] = splitted[1]
		}
	}

	return env, err
}

func generateRandomKey() string {
	b := make([]byte, 64)
	rand.Read(b)
	key := base64.StdEncoding.EncodeToString(b)
	key = strings.ReplaceAll(key, "=", "")
	key = strings.ReplaceAll(key, "+", "")
	key = strings.ReplaceAll(key, "/", "")
	return key
}

func (b Backend) validate() error {
	if b.Type == "" {
		return fmt.Errorf(`Backend "%s" has no "type"`, b.name)
	}
	if b.Path == "" {
		return fmt.Errorf(`Backend "%s" has no "path"`, b.name)
	}
	if b.Key == "" {
		// Check if key is set in environment
		env, _ := b.getEnv()
		if _, found := env["RESTIC_PASSWORD"]; !found {
			if b.RequireKey {
				return fmt.Errorf("backend %s requires a key but none was provided", b.name)
			}
			// No key set in config file or env => generate random key and save file
			key := generateRandomKey()
			b.Key = key
			c := GetConfig()
			tmp := c.Backends[b.name]
			tmp.Key = key
			c.Backends[b.name] = tmp
			if err := c.SaveConfig(); err != nil {
				return err
			}
		}
	}
	env, err := b.getEnv()
	if err != nil {
		return err
	}
	options := ExecuteOptions{Envs: env, Silent: true}

	err = b.EnsureInit()
	if err != nil {
		return err
	}

	cmd := []string{"check"}
	cmd = append(cmd, combineBackendOptions("check", b)...)
	_, _, err = ExecuteResticCommand(options, cmd...)
	return err
}

// EnsureInit initializes the backend if it is not already initialized
func (b Backend) EnsureInit() error {
	env, err := b.getEnv()
	if err != nil {
		return err
	}
	options := ExecuteOptions{Envs: env, Silent: true}

	checkInitCmd := []string{"cat", "config"}
	checkInitCmd = append(checkInitCmd, combineBackendOptions("cat", b)...)
	_, _, err = ExecuteResticCommand(options, checkInitCmd...)

	// Note that `restic` has a special exit code (10) to indicate that the
	// repository does not exist. This exit code was introduced in `restic@0.17.0`
	// on 2024-07-26. We're not using it here because this is a too recent and
	// people on older versions of `restic` won't have this feature work  correctly.
	// See: https://restic.readthedocs.io/en/latest/075_scripting.html#exit-codes
	if err != nil {
		colors.Body.Printf("Initializing backend \"%s\"...\n", b.name)
		initCmd := []string{"init"}
		initCmd = append(initCmd, combineBackendOptions("init", b)...)
		_, _, err := ExecuteResticCommand(options, initCmd...)
		return err
	}

	return err
}

func (b Backend) Exec(args []string) error {
	env, err := b.getEnv()
	if err != nil {
		return err
	}
	options := ExecuteOptions{Envs: env}
	args = append(args, combineBackendOptions("exec", b)...)
	_, out, err := ExecuteResticCommand(options, args...)
	if err != nil {
		colors.Error.Println(out)
		return err
	}
	return nil
}

func (b Backend) ExecDocker(l Location, args []string) (int, string, error) {
	env, err := b.getEnv()
	if err != nil {
		return -1, "", err
	}
	volume := l.From[0]
	options := ExecuteOptions{
		Command: "docker",
		Envs:    env,
	}
	dir := "/data"
	args = append([]string{"restic"}, args...)
	docker := []string{
		"run", "--rm",
		"--entrypoint", "ash",
		"--workdir", dir,
		"--volume", volume + ":" + dir,
	}
	// Use of docker host, not the container host
	if hostname, err := os.Hostname(); err == nil {
		docker = append(docker, "--hostname", hostname)
	}
	switch b.Type {
	case "local":
		actual := env["RESTIC_REPOSITORY"]
		docker = append(docker, "--volume", actual+":"+"/repo")
		env["RESTIC_REPOSITORY"] = "/repo"
	case "b2":
	case "s3":
	case "azure":
	case "gs":
	case "rest":
		// No additional setup needed
	case "rclone":
		// Read host rclone config and mount it into the container
		code, configFile, err := ExecuteCommand(ExecuteOptions{Command: "rclone"}, "config", "file")
		if err != nil {
			return code, "", err
		}
		splitted := strings.Split(strings.TrimSpace(configFile), "\n")
		configFilePath := splitted[len(splitted)-1]
		docker = append(docker, "--volume", configFilePath+":"+"/root/.config/rclone/rclone.conf:ro")
	default:
		return -1, "", fmt.Errorf("Backend type \"%s\" is not supported as volume endpoint", b.Type)
	}
	for key, value := range env {
		docker = append(docker, "--env", key+"="+value)
	}

	docker = append(docker, flags.DOCKER_IMAGE, "-c", strings.Join(args, " "))
	return ExecuteCommand(options, docker...)
}
