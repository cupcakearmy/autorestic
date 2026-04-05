package internal

import (
	"os"
	"path"
	"sync"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetType(t *testing.T) {

	t.Run("TypeLocal", func(t *testing.T) {
		l := Location{
			Type: "local",
		}
		result, err := l.getType()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		assertEqual(t, result, TypeLocal)
	})

	t.Run("TypeVolume", func(t *testing.T) {
		l := Location{
			Type: "volume",
		}
		result, err := l.getType()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		assertEqual(t, result, TypeVolume)
	})

	t.Run("Empty type", func(t *testing.T) {
		l := Location{
			Type: "",
		}
		result, err := l.getType()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		assertEqual(t, result, TypeLocal)
	})

	t.Run("Invalid type", func(t *testing.T) {
		l := Location{
			Type: "foo",
		}
		_, err := l.getType()
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestBuildTag(t *testing.T) {
	result := buildTag("foo", "bar")
	expected := "ar:foo:bar"
	assertEqual(t, result, expected)
}

func TestGetLocationTags(t *testing.T) {
	l := Location{
		name: "foo",
	}
	result := l.getLocationTags()
	expected := "ar:location:foo"
	assertEqual(t, result, expected)
}

func TestHasBackend(t *testing.T) {
	t.Run("backend present", func(t *testing.T) {
		l := Location{
			name: "foo",
			To:   []string{"foo", "bar"},
		}
		result := l.hasBackend("foo")
		assertEqual(t, result, true)
	})

	t.Run("backend absent", func(t *testing.T) {
		l := Location{
			name: "foo",
			To:   []string{"bar", "baz"},
		}
		result := l.hasBackend("foo")
		assertEqual(t, result, false)
	})
}

func TestBuildRestoreCommand(t *testing.T) {
	l := Location{
		name: "foo",
	}
	result := buildRestoreCommand(l, "to", "snapshot", []string{"options"})
	expected := []string{"restore", "--target", "to", "--tag", "ar:location:foo", "snapshot", "options"}
	assertSliceEqual(t, result, expected)
}

func TestLocationBackupWithMock(t *testing.T) {
	// Backup original
	originalExecutor := DefaultExecutor
	defer func() { DefaultExecutor = originalExecutor }()

	// Inject mock
	mock := &MockExecutor{
		ExecuteResticFunc: func(options ExecuteOptions, args ...string) (int, string, error) {
			assert.Equal(t, "backup", args[0])
			return 0, "success", nil
		},
	}
	DefaultExecutor = mock

	// Setup dummy config
	workDir := t.TempDir()
	configFile := path.Join(workDir, ".autorestic.yml")
	err := os.WriteFile(configFile, []byte("version: 2"), 0644)
	assert.NoError(t, err)

	viper.Reset()
	viper.SetConfigFile(configFile)
	viper.Set("version", 2)
	// Register test-backend
	viper.Set("backends.test-backend.type", "local")
	viper.Set("backends.test-backend.path", workDir)

	config = nil
	once = sync.Once{}

	loc := Location{
		name: "test-location",
		To:   []string{"test-backend"},
		From: []string{"/"},
		Type: "local",
	}

	errs := loc.Backup(false, false, "")
	if len(errs) != 0 {
		t.Errorf("expected no error, got %v", errs)
	}
}
