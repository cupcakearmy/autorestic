package cmd

import (
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

const xdgConfigHome = "XDG_CONFIG_HOME"

func assertContains(t *testing.T, array []string, element string) {
	if !slices.Contains(array, element) {
		t.Errorf("Expected %s to be contained in %s", element, array)
	}
}

func TestConfigResolving(t *testing.T) {
	t.Run("~/.config/autorestic is used if XDG_CONFIG_HOME is not set", func(t *testing.T) {
		// Override env using testing so that env gets restored after test
		t.Setenv(xdgConfigHome, "")
		_ = os.Unsetenv("XDG_CONFIG_HOME")
		configPaths := getConfigPaths()
		homeDir, _ := homedir.Dir()
		expectedConfigPath := filepath.Join(homeDir, ".config/autorestic")
		assertContains(t, configPaths, expectedConfigPath)
	})

	t.Run("XDG_CONFIG_HOME is respected if set", func(t *testing.T) {
		t.Setenv(xdgConfigHome, "/foo/bar")

		configPaths := getConfigPaths()
		assertContains(t, configPaths, filepath.Join("/", "foo", "bar", "autorestic"))
	})
}
