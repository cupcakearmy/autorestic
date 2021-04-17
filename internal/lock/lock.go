package lock

import (
	"errors"
	"path"
	"sync"

	"github.com/spf13/viper"
)

var lock *viper.Viper
var file string
var once sync.Once

func getLock() *viper.Viper {
	if lock == nil {

		once.Do(func() {
			lock = viper.New()
			lock.SetDefault("running", false)
			p := path.Dir(viper.ConfigFileUsed())
			file = path.Join(p, ".autorestic.lock.yml")
			lock.SetConfigFile(file)
			lock.SetConfigType("yml")
			lock.ReadInConfig()
		})
	}
	return lock
}

func setLock(locked bool) error {
	lock := getLock()
	if locked {
		running := lock.GetBool("running")
		if running {
			panic(errors.New("an instance is already running"))
		}
	}
	lock.Set("running", locked)
	if err := lock.WriteConfigAs(file); err != nil {
		return err
	}
	return nil
}

func GetCron(location string) int64 {
	lock := getLock()
	return lock.GetInt64("cron." + location)
}

func SetCron(location string, value int64) {
	lock.Set("cron."+location, value)
	lock.WriteConfigAs(file)
}

func Lock() error {
	return setLock(true)
}

func Unlock() error {
	return setLock(false)
}
