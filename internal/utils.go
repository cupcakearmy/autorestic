package internal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/cupcakearmy/autorestic/internal/colors"
)

var RESTIC_BIN string

func CheckIfCommandIsCallable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func CheckIfResticIsCallable() bool {
	return CheckIfCommandIsCallable(RESTIC_BIN)
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

	if VERBOSE {
		colors.Faint.Printf("> Executing: %s\n", cmd)
	}

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
	options.Command = RESTIC_BIN
	var c = GetConfig()
	var optionsAsString = getOptions(c.Global, "")
	args = append(optionsAsString, args...)
	return ExecuteCommand(options, args...)
}

func CopyFile(from, to string) error {
	original, err := os.Open(from)
	if err != nil {
		return nil
	}
	defer original.Close()

	new, err := os.Create(to)
	if err != nil {
		return nil
	}
	defer new.Close()

	if _, err := io.Copy(new, original); err != nil {
		return err
	}
	return nil
}

func CheckIfVolumeExists(volume string) bool {
	_, err := ExecuteCommand(ExecuteOptions{Command: "docker"}, "volume", "inspect", volume)
	return err == nil
}
