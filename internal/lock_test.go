package internal

import (
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"
	"testing"

	"github.com/cupcakearmy/autorestic/internal/flags"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// All tests must share the same lock file as it is only initialized once
func setup(t *testing.T) {
	t.Helper()
	cleanup := func() {
		flags.LOCKFILE = ""
		config = nil
		once = sync.Once{}
		viper.Reset()
	}

	cleanup()
	d := t.TempDir()
	viper.SetConfigFile(d + "/.autorestic.yml")
	viper.Set("version", 2)
	viper.WriteConfig()

	t.Cleanup(cleanup)
}

func TestGetLockfilePath(t *testing.T) {
	t.Run("user specified", func(t *testing.T) {
		testCases := []struct {
			name     string
			flag     string
			config   string
			expected string
		}{
			{
				name:     "flag and config",
				flag:     "/flag.lock.yml",
				config:   "/config.lock.yml",
				expected: "/flag.lock.yml",
			},
			{
				name:     "flag only",
				flag:     "/flag.lock.yml",
				config:   "",
				expected: "/flag.lock.yml",
			},
			{
				name:     "config only",
				flag:     "",
				config:   "/config.lock.yml",
				expected: "/config.lock.yml",
			},
		}
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				setup(t)
				flags.LOCKFILE = testCase.flag
				if testCase.config != "" {
					viper.Set("lockfile", testCase.config)
					err := viper.WriteConfig()
					assert.NoError(t, err)
				}

				p := getLockfilePath()
				assert.Equal(t, testCase.expected, p)
			})
		}
	})

	t.Run("default", func(t *testing.T) {
		setup(t)

		configPath := viper.ConfigFileUsed()
		expectedLockfile := path.Join(path.Dir(configPath), ".autorestic.lock.yml")

		p := getLockfilePath()
		assert.Equal(t, expectedLockfile, p)
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

	// locking a locked instance exits the instance
	// this trick to capture os.Exit(1) is discussed here:
	// https://talks.golang.org/2014/testing.slide#23
	t.Run("lock twice", func(t *testing.T) {
		if os.Getenv("CRASH") == "1" {
			err := Lock()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			// should fail
			Lock()
		}

		cmd := exec.Command(os.Args[0], "-test.run=TestLock/lock_twice")
		cmd.Env = append(os.Environ(), "CRASH=1")
		err := cmd.Run()

		err, ok := err.(*exec.ExitError)
		if !ok {
			t.Error("unexpected error")
		}
		expected := "exit status 1"
		if err.Error() != expected {
			t.Errorf("got %q, want %q", err.Error(), expected)
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
