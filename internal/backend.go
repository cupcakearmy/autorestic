package internal

import (
	"fmt"
)

type Backend struct {
	Type string            `mapstructure:"type"`
	Path string            `mapstructure:"path"`
	Key  string            `mapstructure:"key"`
	Env  map[string]string `mapstructure:"env"`
}

func (b Backend) generateRepo() (string, error) {
	switch b.Type {
	case "local":
		return GetPathRelativeToConfig(b.Path), nil
	case "b2", "azure", "gs", "s3", "sftp", "rest":
		return fmt.Sprintf("%s:%s", b.Type, b.Path), nil
	default:
		return "", fmt.Errorf("backend type \"%s\" is invalid", b.Type)
	}
}

func (b Backend) getEnv() map[string]string {
	env := make(map[string]string)
	env["RESTIC_PASSWORD"] = b.Key
	repo, err := b.generateRepo()
	if err != nil {
		panic(err)
	}
	env["RESTIC_REPOSITORY"] = repo
	return env
}

func (b Backend) validate() error {
	options := ExecuteOptions{Envs: b.getEnv()}
	// Check if already initialized
	_, err := ExecuteResticCommand(options, "snapshots")
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
	options := ExecuteOptions{Envs: b.getEnv()}
	out, err := ExecuteResticCommand(options, args...)
	fmt.Println(out)
	return err
}
