package integration_test

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestVersion(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "--version")
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Contains(t, string(out), "autorestic")
}
