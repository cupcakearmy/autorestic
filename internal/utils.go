package internal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/cupcakearmy/autorestic/internal/colors"
	"github.com/cupcakearmy/autorestic/internal/flags"
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

func ExecuteCommand(options ExecuteOptions, args ...string) (int, string, error) {
	cmd := exec.Command(options.Command, args...)
	env := os.Environ()
	for k, v := range options.Envs {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = env
	cmd.Dir = options.Dir

	if flags.VERBOSE {
		colors.Faint.Printf("> Executing: %s\n", cmd)
	}

	var out bytes.Buffer
	var error bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &error
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), error.String(), err
		} else {
			return -1, error.String(), err
		}
	}
	return 0, out.String(), nil
}

func ExecuteResticCommand(options ExecuteOptions, args ...string) (int, string, error) {
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
	_, _, err := ExecuteCommand(ExecuteOptions{Command: "docker"}, "volume", "inspect", volume)
	return err == nil
}

func ArrayContains[T comparable](arr []T, needle T) bool {
	for _, item := range arr {
		if item == needle {
			return true
		}
	}
	return false
}
