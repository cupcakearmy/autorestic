package integration_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func runAutorestic(t *testing.T, dir string, configPath string, args ...string) string {
	// Find project root dynamically
	wd, err := os.Getwd()
	assert.NoError(t, err)
	root := filepath.Dir(wd)

	mainGo := filepath.Join(root, "main.go")

	// Get restic path from mise
	resticPath := "/Users/cupcakearmy/.local/share/mise/installs/restic/0.18.1/restic"

	// Convert configPath to absolute path
	absConfigPath, err := filepath.Abs(configPath)
	assert.NoError(t, err)

	cmd := exec.Command("go", append([]string{"run", mainGo, "--restic-bin", resticPath, "-c", absConfigPath}, args...)...)

	// Run from root to find go.mod
	cmd.Dir = root

	// Add mise path to environment and set dummy password
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PATH="+os.Getenv("PATH"))
	cmd.Env = append(cmd.Env, "RESTIC_PASSWORD=password")

	output, err := cmd.CombinedOutput()
	// NOTE: We don't assert NoError here because tests might expect failures
	return string(output)
}

func initRepo(t *testing.T, repoPath string) {
	resticPath := "/Users/cupcakearmy/.local/share/mise/installs/restic/0.18.1/restic"
	cmd := exec.Command(resticPath, "-r", repoPath, "init")
	cmd.Env = append(os.Environ(), "RESTIC_PASSWORD=password")
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, string(output))
}

func TestAutoresticCheck(t *testing.T) {
	tempDir := t.TempDir()

	configContent := `
version: 2
locations:
  my-data:
    from: ` + tempDir + `
    to: local
backends:
  local:
    type: local
    path: ` + filepath.Join(tempDir, "repo") + `
    key: password
`
	configPath := filepath.Join(tempDir, "autorestic.yml")
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	_ = runAutorestic(t, tempDir, configPath, "check")
}

func TestBackupRestore(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Create a source file
	sourceFile := "source.txt"
	err := os.WriteFile(filepath.Join(tempDir, sourceFile), []byte("hello world"), 0644)
	assert.NoError(t, err)

	repoPath := filepath.Join(tempDir, "repo")
	initRepo(t, repoPath)

	configContent := `
version: 2
locations:
  my-data:
    from: 
      - ` + sourceFile + `
    to: local
backends:
  local:
    type: local
    path: ` + repoPath + `
    key: password
`
	configPath := filepath.Join(tempDir, "autorestic.yml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	// 2. Backup
	output := runAutorestic(t, tempDir, configPath, "backup", "-l", "my-data")
	assert.Contains(t, output, "Done")

	// 3. Restore
	restoreDir := filepath.Join(tempDir, "restore")
	err = os.MkdirAll(restoreDir, 0755)
	assert.NoError(t, err)
	output = runAutorestic(t, tempDir, configPath, "restore", "-l", "my-data", "--to", restoreDir)
	t.Logf("Restore output: %s", output)

	// DEBUG: List files in restoreDir and subdirectories
	err = filepath.Walk(restoreDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		t.Logf("Found: %s", path)
		return nil
	})
	assert.NoError(t, err)

	// 4. Verify
	// It might be nested depending on how it was backed up
	// Let's look for source.txt in any subdirectory
	var restoredFile string
	err = filepath.Walk(restoreDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() == "source.txt" {
			restoredFile = path
		}
		return nil
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, restoredFile, "source.txt not found")

	content, err := os.ReadFile(restoredFile)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", string(content))
}

func TestHooks(t *testing.T) {
	tempDir := t.TempDir()

	// Create a dummy file to back up
	sourceFile := "source.txt"
	err := os.WriteFile(filepath.Join(tempDir, sourceFile), []byte("data"), 0644)
	assert.NoError(t, err)

	repoPath := filepath.Join(tempDir, "repo")
	initRepo(t, repoPath)

	configContent := `
version: 2
locations:
  my-data:
    from: 
      - ` + sourceFile + `
    to: local
    hooks:
      before:
        - touch before.txt
      after:
        - touch after.txt
backends:
  local:
    type: local
    path: ` + repoPath + `
    key: password
`
	configPath := filepath.Join(tempDir, "autorestic.yml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Run backup
	runAutorestic(t, tempDir, configPath, "backup", "-l", "my-data")

	// Verify
	assert.FileExists(t, filepath.Join(tempDir, "before.txt"))
	assert.FileExists(t, filepath.Join(tempDir, "after.txt"))
}

func TestCopy(t *testing.T) {
	tempDir := t.TempDir()

	// Create a dummy file to back up
	sourceFile := "source.txt"
	err := os.WriteFile(filepath.Join(tempDir, sourceFile), []byte("data"), 0644)
	assert.NoError(t, err)

	repoPath := filepath.Join(tempDir, "repo")
	initRepo(t, repoPath)

	remoteRepoPath := filepath.Join(tempDir, "remote_repo")
	initRepo(t, remoteRepoPath)

	configContent := `
version: 2
locations:
  my-data:
    from: 
      - ` + sourceFile + `
    to: local
    copy:
      local:
        - remote
backends:
  local:
    type: local
    path: ` + repoPath + `
    key: password
  remote:
    type: local
    path: ` + remoteRepoPath + `
    key: password
`
	configPath := filepath.Join(tempDir, "autorestic.yml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Run backup
	output := runAutorestic(t, tempDir, configPath, "backup", "-l", "my-data")

	// Verify copy in output
	assert.Contains(t, output, "Copying local → remote")
}
