package internal

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func CheckIfCommandIsCallable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func CheckIfResticIsCallable() bool {
	return CheckIfCommandIsCallable("restic")
}

type ExecuteOptions struct {
	Command string
	Envs    map[string]string
	Dir     string
}

func ExecuteCommand(options ExecuteOptions, args ...string) (string, error) {
	cmd := exec.Command(options.Command, args...)
	env := os.Environ()
	for k, v := range options.Envs {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = env
	cmd.Dir = options.Dir

	var out bytes.Buffer
	var error bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &error
	err := cmd.Run()
	if err != nil {
		return error.String(), err
	}
	return out.String(), nil
}

func ExecuteResticCommand(options ExecuteOptions, args ...string) (string, error) {
	options.Command = "restic"
	return ExecuteCommand(options, args...)
}
