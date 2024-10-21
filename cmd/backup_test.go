package cmd

import (
	"os"
	"path"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func runCmd(t *testing.T, args ...string) error {
	t.Helper()

	viper.Reset()
	rootCmd.SetArgs(args)

	err := rootCmd.Execute()
	return err
}

func TestBackupCmd(t *testing.T) {
	workDir := t.TempDir()

	// Prepare content to be backed up
	locationDir := path.Join(workDir, "my-location")
	err := os.Mkdir(locationDir, 0750)
	assert.Nil(t, err)
	err = os.WriteFile(path.Join(locationDir, "back-me-up.txt"), []byte("hello world"), 0640)
	assert.Nil(t, err)

	// Write config file
	config, err := yaml.Marshal(map[string]interface{}{
		"version": 2,
		"locations": map[string]map[string]interface{}{
			"my-location": {
				"type": "local",
				"from": []string{locationDir},
				"to":   []string{"test"},
			},
		},
		"backends": map[string]map[string]interface{}{
			"test": {
				"type": "local",
				"path": path.Join(workDir, "test-backend"),
				"key":  "supersecret",
			},
		},
	})
	assert.Nil(t, err)
	configPath := path.Join(workDir, ".autorestic.yml")
	err = os.WriteFile(configPath, config, 0640)
	assert.Nil(t, err)

	// Init repo (not initialized by default)
	err = runCmd(t, "exec", "--ci", "-a", "-c", configPath, "init")
	assert.Nil(t, err)

	// Do the backup
	err = runCmd(t, "backup", "--ci", "-a", "-c", configPath)
	assert.Nil(t, err)

	// Restore in a separate dir
	restoreDir := path.Join(workDir, "restore")
	err = runCmd(t, "restore", "--ci", "-c", configPath, "-l", "my-location", "--to", restoreDir)
	assert.Nil(t, err)

	// Check restored file
	restoredContent, err := os.ReadFile(path.Join(restoreDir, locationDir, "back-me-up.txt"))
	assert.Nil(t, err)
	assert.Equal(t, "hello world", string(restoredContent))
}
