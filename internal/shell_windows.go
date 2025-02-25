//go:build windows

package internal

import (
	"os/exec"
	"syscall"
)

const shellName string = "PowerShell"
const shellArgs string = "-Command"

func execCommand(command string, args ...string) *exec.Cmd {
	cmd := exec.Command(command, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}
