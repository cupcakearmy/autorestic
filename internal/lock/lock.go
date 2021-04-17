package lock

import (
	"os"
	"path"
	"sync"

	"github.com/cupcakearmy/autorestic/internal/colors"
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
			colors.Error.Println("an instance is already running. exiting")
			os.Exit(1)
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
