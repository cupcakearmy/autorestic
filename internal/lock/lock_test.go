package lock

import (
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/spf13/viper"
)

var testDirectory = "autorestic_test_tmp"

// All tests must share the same lock file as it is only initialized once
func setup(t *testing.T) {
	d, err := os.MkdirTemp("", testDirectory)
	if err != nil {
		log.Fatalf("error creating temp dir: %v", err)
		return
	}
	// set config file location
	viper.SetConfigFile(d + "/.autorestic.yml")

	t.Cleanup(func() {
		os.RemoveAll(d)
		viper.Reset()
	})
}

func TestLock(t *testing.T) {
	setup(t)

	t.Run("getLock", func(t *testing.T) {
		result := getLock().GetBool(RUNNING)

		if result {
			t.Errorf("got %v, want %v", result, false)
		}
	})

	t.Run("lock", func(t *testing.T) {
		lock, err := setLockValue(RUNNING, true)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		result := lock.GetBool(RUNNING)
		if !result {
			t.Errorf("got %v, want %v", result, true)
		}
	})

	t.Run("unlock", func(t *testing.T) {
		lock, err := setLockValue(RUNNING, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		result := lock.GetBool(RUNNING)
		if result {
			t.Errorf("got %v, want %v", result, false)
		}
	})

	t.Run("set cron", func(t *testing.T) {
		expected := int64(5)
		SetCron("foo", expected)

		result, err := strconv.ParseInt(getLock().GetString("cron.foo"), 10, 64)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != expected {
			t.Errorf("got %d, want %d", result, expected)
		}
	})

	t.Run("get cron", func(t *testing.T) {
		expected := int64(5)
		result := GetCron("foo")

		if result != expected {
			t.Errorf("got %d, want %d", result, expected)
		}
	})
}
