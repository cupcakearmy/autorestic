package lock

import (
	"os"
	"path"
	"sync"

	"github.com/cupcakearmy/autorestic/internal/colors"
	"github.com/cupcakearmy/autorestic/internal/flags"
	"github.com/spf13/viper"
)

var lock *viper.Viper
var file string
var once sync.Once

const (
	RUNNING = "running"
)

func getLock() *viper.Viper {
	if lock == nil {

		once.Do(func() {
			lock = viper.New()
			lock.SetDefault("running", false)
			p := viper.ConfigFileUsed()
			if p == "" {
				colors.Error.Println("cannot lock before reading config location")
				os.Exit(1)
			}
			file = path.Join(path.Dir(p), ".autorestic.lock.yml")
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
