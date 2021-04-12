package internal

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cupcakearmy/autorestic/internal/colors"
	"github.com/spf13/viper"
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
	if b.Name == "" {
		return fmt.Errorf(`Backend has no "name"`)
	}
	if b.Type == "" {
		return fmt.Errorf(`Backend "%s" has no "type"`, b.Name)
	}
	if b.Path == "" {
		return fmt.Errorf(`Backend "%s" has no "path"`, b.Name)
	}
	if b.Key == "" {
		key := generateRandomKey()
		b.Key = key
		c := GetConfig()
		for i, backend := range c.Backends {
			if backend.Name == b.Name {
				c.Backends[i].Key = key
				break
			}
		}
		file := viper.ConfigFileUsed()
		if err := CopyFile(file, file+".old"); err != nil {
			return err
		}
		colors.Secondary.Println("Saved a backup copy of your file next the the original.")
		viper.Set("backends", c.Backends)
		viper.WriteConfig()
	}
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
		colors.Faint.Println(out)
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
	if VERBOSE {
		colors.Faint.Println(out)
	}
	return err
}
