package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecuteCommandWithMock(t *testing.T) {
	// Backup original
	originalExecutor := DefaultExecutor
	defer func() { DefaultExecutor = originalExecutor }()

	// Inject mock
	mock := &MockExecutor{
		ExecuteFunc: func(options ExecuteOptions, args ...string) (int, string, error) {
			assert.Equal(t, "docker", options.Command)
			return 0, "mock output", nil
		},
	}
	DefaultExecutor = mock

	code, out, err := ExecuteCommand(ExecuteOptions{Command: "docker"}, "info")

	assert.NoError(t, err)
	assert.Equal(t, 0, code)
	assert.Equal(t, "mock output", out)
}
