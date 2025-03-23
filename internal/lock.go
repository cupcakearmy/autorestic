package internal

import (
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/cupcakearmy/autorestic/internal/colors"
	"github.com/cupcakearmy/autorestic/internal/flags"
	"github.com/spf13/viper"
)

var lock *viper.Viper
var file string
var lockOnce sync.Once

const (
	RUNNING = "running"
)

// getLockfilePath returns the path to the lockfile. The path for the lockfile
// can be sources from multiple places If flags.LOCKFILE is set, its value is
// used; if the config has the `lockfile` option set, its value is used;
// otherwise the path is generated relative to the config file.
func getLockfilePath() string {
	if flags.LOCKFILE != "" {
		abs, err := filepath.Abs(flags.LOCKFILE)
		if err != nil {
			return flags.LOCKFILE
		}
		return abs
	}

	if lockfile := GetConfig().Lockfile; lockfile != "" {
		abs, err := filepath.Abs(lockfile)
		if err != nil {
			return lockfile
		}
		return abs
	}

	p := viper.ConfigFileUsed()
	if p == "" {
		colors.Error.Println("cannot lock before reading config location")
		os.Exit(1)
	}
	return path.Join(path.Dir(p), ".autorestic.lock.yml")
}

func getLock() *viper.Viper {
	if lock == nil {
		lockOnce.Do(func() {
			lock = viper.New()
			lock.SetDefault("running", false)
			file = getLockfilePath()
			if !flags.CRON_LEAN {
				colors.Faint.Println("Using lock:\t", file)
			}
			lock.SetConfigFile(file)
			lock.SetConfigType("yml")
			lock.ReadInConfig()
		})
	}
	return lock
}

func setLockValue(key string, value interface{}) (*viper.Viper, error) {
	lock := getLock()

	if key == RUNNING {
		value := value.(bool)
		if value && lock.GetBool(key) {
			colors.Error.Println("an instance is already running. exiting")
			os.Exit(1)
		}
	}

	lock.Set(key, value)
	if err := lock.WriteConfigAs(file); err != nil {
		return nil, err
	}
	return lock, nil
}

func GetCron(location string) int64 {
	return getLock().GetInt64("cron." + location)
}

func SetCron(location string, value int64) {
	setLockValue("cron."+location, value)
}

func Lock() error {
	_, err := setLockValue(RUNNING, true)
	return err
}

func Unlock() error {
	_, err := setLockValue(RUNNING, false)
	return err
}
