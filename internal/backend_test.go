package internal

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/cupcakearmy/autorestic/internal/flags"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGenerateRepo(t *testing.T) {

	t.Run("empty backend", func(t *testing.T) {
		b := Backend{
			name: "empty backend",
			Type: "",
		}
		_, err := b.generateRepo()
		if err == nil {
			t.Errorf("Error expected for empty backend type")
		}
	})

	t.Run("local backend", func(t *testing.T) {
		b := Backend{
			name: "local backend",
			Type: "local",
			Path: "/foo/bar",
		}
		result, err := b.generateRepo()
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}
		assertEqual(t, result, "/foo/bar")
	})

	t.Run("local backend with homedir prefix", func(t *testing.T) {
		b := Backend{
			name: "local backend",
			Type: "local",
			Path: "~/foo/bar",
		}
		result, err := b.generateRepo()
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}
		assertEqual(t, result, fmt.Sprintf("%s/foo/bar", os.Getenv("HOME")))
	})

	t.Run("local backend with config file", func(t *testing.T) {
		// config file path should always be present from initConfig
		viper.SetConfigFile("/tmp/.autorestic.yml")
		defer viper.Reset()

		b := Backend{
			name: "local backend",
			Type: "local",
		}
		result, err := b.generateRepo()
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}
		assertEqual(t, result, "/tmp")
	})

	t.Run("rest backend with valid path", func(t *testing.T) {
		b := Backend{
			name: "rest backend",
			Type: "rest",
			Path: "http://localhost:8000/foo",
		}
		result, err := b.generateRepo()
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}
		assertEqual(t, result, "rest:http://localhost:8000/foo")
	})

	t.Run("rest backend with user", func(t *testing.T) {
		b := Backend{
			name: "rest backend",
			Type: "rest",
			Path: "http://localhost:8000/foo",
			Rest: BackendRest{
				User:     "user",
				Password: "",
			},
		}
		result, err := b.generateRepo()
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}
		assertEqual(t, result, "rest:http://user@localhost:8000/foo")
	})

	t.Run("rest backend with user and password", func(t *testing.T) {
		b := Backend{
			name: "rest backend",
			Type: "rest",
			Path: "http://localhost:8000/foo",
			Rest: BackendRest{
				User:     "user",
				Password: "pass",
			},
		}
		result, err := b.generateRepo()
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}
		assertEqual(t, result, "rest:http://user:pass@localhost:8000/foo")
	})

	backendTests := []struct {
		name    string
		backend Backend
		want    string
	}{
		{name: "b2 backend", backend: Backend{name: "b2", Type: "b2", Path: "foo"}, want: "b2:foo"},
		{name: "azure backend", backend: Backend{name: "azure", Type: "azure", Path: "foo"}, want: "azure:foo"},
		{name: "gs backend", backend: Backend{name: "gs", Type: "gs", Path: "foo"}, want: "gs:foo"},
		{name: "s3 backend", backend: Backend{name: "s3", Type: "s3", Path: "foo"}, want: "s3:foo"},
		{name: "sftp backend", backend: Backend{name: "sftp", Type: "sftp", Path: "foo"}, want: "sftp:foo"},
		{name: "rclone backend", backend: Backend{name: "rclone", Type: "rclone", Path: "foo"}, want: "rclone:foo"},
	}

	for _, tt := range backendTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.backend.generateRepo()
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestGetEnv(t *testing.T) {
	t.Run("env in key field", func(t *testing.T) {
		b := Backend{
			name: "",
			Type: "local",
			Path: "/foo/bar",
			Key:  "secret123",
		}
		result, err := b.getEnv()
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}
		assertEqual(t, result["RESTIC_REPOSITORY"], "/foo/bar")
		assertEqual(t, result["RESTIC_PASSWORD"], "secret123")
	})

	t.Run("env in config file", func(t *testing.T) {
		b := Backend{
			name: "",
			Type: "local",
			Path: "/foo/bar",
			Env: map[string]string{
				"B2_ACCOUNT_ID":  "foo123",
				"B2_ACCOUNT_KEY": "foo456",
			},
		}
		result, err := b.getEnv()
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}
		assertEqual(t, result["RESTIC_REPOSITORY"], "/foo/bar")
		assertEqual(t, result["RESTIC_PASSWORD"], "")
		assertEqual(t, result["B2_ACCOUNT_ID"], "foo123")
		assertEqual(t, result["B2_ACCOUNT_KEY"], "foo456")
	})

	t.Run("env in Envfile or env vars", func(t *testing.T) {
		// generate env variables
		// TODO better way to teardown
		defer os.Unsetenv("AUTORESTIC_FOO_RESTIC_PASSWORD")
		defer os.Unsetenv("AUTORESTIC_FOO_B2_ACCOUNT_ID")
		defer os.Unsetenv("AUTORESTIC_FOO_B2_ACCOUNT_KEY")
		os.Setenv("AUTORESTIC_FOO_RESTIC_PASSWORD", "secret123")
		os.Setenv("AUTORESTIC_FOO_B2_ACCOUNT_ID", "foo123")
		os.Setenv("AUTORESTIC_FOO_B2_ACCOUNT_KEY", "foo456")

		b := Backend{
			name: "foo",
			Type: "local",
			Path: "/foo/bar",
		}
		result, err := b.getEnv()
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}
		assertEqual(t, result["RESTIC_REPOSITORY"], "/foo/bar")
		assertEqual(t, result["RESTIC_PASSWORD"], "secret123")
		assertEqual(t, result["B2_ACCOUNT_ID"], "foo123")
		assertEqual(t, result["B2_ACCOUNT_KEY"], "foo456")
	})

	for _, char := range "@-_:/" {
		t.Run(fmt.Sprintf("env var with special char (%c)", char), func(t *testing.T) {
			// generate env variables
			// TODO better way to teardown
			defer os.Unsetenv("AUTORESTIC_FOO_BAR_RESTIC_PASSWORD")
			defer os.Unsetenv("AUTORESTIC_FOO_BAR_B2_ACCOUNT_ID")
			defer os.Unsetenv("AUTORESTIC_FOO_BAR_B2_ACCOUNT_KEY")
			os.Setenv("AUTORESTIC_FOO_BAR_RESTIC_PASSWORD", "secret123")
			os.Setenv("AUTORESTIC_FOO_BAR_B2_ACCOUNT_ID", "foo123")
			os.Setenv("AUTORESTIC_FOO_BAR_B2_ACCOUNT_KEY", "foo456")

			b := Backend{
				name: fmt.Sprintf("foo%cbar", char),
				Type: "local",
				Path: "/foo/bar",
			}
			result, err := b.getEnv()
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}
			assertEqual(t, result["RESTIC_REPOSITORY"], "/foo/bar")
			assertEqual(t, result["RESTIC_PASSWORD"], "secret123")
			assertEqual(t, result["B2_ACCOUNT_ID"], "foo123")
			assertEqual(t, result["B2_ACCOUNT_KEY"], "foo456")
		})
	}
}

func TestValidate(t *testing.T) {
	t.Run("no type given", func(t *testing.T) {
		b := Backend{
			name: "foo",
			Type: "",
			Path: "/foo/bar",
		}
		err := b.validate()
		if err == nil {
			t.Error("expected to get error")
		}
		assertEqual(t, err.Error(), "Backend \"foo\" has no \"type\"")
	})

	t.Run("no path given", func(t *testing.T) {
		b := Backend{
			name: "foo",
			Type: "local",
			Path: "",
		}
		err := b.validate()
		if err == nil {
			t.Error("expected to get error")
		}
		assertEqual(t, err.Error(), "Backend \"foo\" has no \"path\"")
	})

	t.Run("require key with no key", func(t *testing.T) {
		b := Backend{
			name:       "foo",
			Type:       "local",
			Path:       "~/foo/bar",
			RequireKey: true,
		}
		err := b.validate()
		fmt.Printf("error: %v\n", err)
		assert.EqualError(t, err, "backend foo requires a key but none was provided")
	})
}

func TestValidateInitsRepo(t *testing.T) {
	// This is normally initialized by the cobra commands but they don't run in
	// this test so we do it ourselves.
	flags.RESTIC_BIN = "restic"

	workDir := t.TempDir()

	b := Backend{
		name: "test",
		Type: "local",
		Path: path.Join(workDir, "backend"),
		Key:  "supersecret",
	}

	config = &Config{Backends: map[string]Backend{"test": b}}
	defer func() { config = nil }()

	// Check should fail because the repo doesn't exist
	err := b.Exec([]string{"check"})
	assert.Error(t, err)

	err = b.validate()
	assert.NoError(t, err)

	// Check should pass now
	err = b.Exec([]string{"check"})
	assert.NoError(t, err)
}

func TestEnsureInit(t *testing.T) {
	// This is normally initialized by the cobra commands but they don't run in
	// this test so we do it ourselves.
	flags.RESTIC_BIN = "restic"

	workDir := t.TempDir()

	b := Backend{
		name: "test",
		Type: "local",
		Path: path.Join(workDir, "backend"),
		Key:  "supersecret",
	}

	config = &Config{Backends: map[string]Backend{"test": b}}
	defer func() { config = nil }()

	// Check should fail because the repo doesn't exist
	err := b.Exec([]string{"check"})
	assert.Error(t, err)

	err = b.EnsureInit()
	assert.NoError(t, err)

	// Check should pass now
	err = b.Exec([]string{"check"})
	assert.NoError(t, err)

	// Run again to make sure it's idempotent
	err = b.EnsureInit()
	assert.NoError(t, err)

	// Check should still pass
	err = b.Exec([]string{"check"})
	assert.NoError(t, err)
}
