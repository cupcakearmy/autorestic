//go:build !windows

package internal

import "os/exec"

const shellName string = "bash"
const shellArgs string = "-c"

func execCommand(cmd string, args ...string) *exec.Cmd {
	return exec.Command(cmd, args...)
}
