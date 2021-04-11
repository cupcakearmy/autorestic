package internal

import (
	"fmt"
)

type Backend struct {
	Name string            `mapstructure:"name"`
	Type string            `mapstructure:"type"`
	Path string            `mapstructure:"path"`
	Key  string            `mapstructure:"key"`
	Env  map[string]string `mapstructure:"env"`
}

func GetBackend(name string) (Backend, bool) {
	c := GetConfig()
	for _, b := range c.Backends {
		if b.Name == name {
			return b, true
		}
	}
	return Backend{}, false
}

func (b Backend) generateRepo() (string, error) {
	switch b.Type {
	case "local":
		return GetPathRelativeToConfig(b.Path)
	case "b2", "azure", "gs", "s3", "sftp", "rest":
		return fmt.Sprintf("%s:%s", b.Type, b.Path), nil
	default:
		return "", fmt.Errorf("backend type \"%s\" is invalid", b.Type)
	}
}

func (b Backend) getEnv() (map[string]string, error) {
	env := make(map[string]string)
	env["RESTIC_PASSWORD"] = b.Key
	repo, err := b.generateRepo()
	env["RESTIC_REPOSITORY"] = repo
	return env, err
}

func (b Backend) validate() error {
	env, err := b.getEnv()
	if err != nil {
		return err
	}
	options := ExecuteOptions{Envs: env}
	// Check if already initialized
	_, err = ExecuteResticCommand(options, "snapshots")
	if err == nil {
		return nil
	} else {
		// If not initialize
		out, err := ExecuteResticCommand(options, "init")
		fmt.Println(out)
		return err
	}
}

func (b Backend) Exec(args []string) error {
	env, err := b.getEnv()
	if err != nil {
		return err
	}
	options := ExecuteOptions{Envs: env}
	out, err := ExecuteResticCommand(options, args...)
	fmt.Println(out)
	return err
}
