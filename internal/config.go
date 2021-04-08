package internal

import (
	"fmt"
	"log"
	"path"
	"strings"
	"sync"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Locations map[string]Location `mapstructure:"locations"`
	Backends  map[string]Backend  `mapstructure:"backends"`
}

var once sync.Once
var config *Config

func GetConfig() *Config {
	if config == nil {
		once.Do(func() {
			config = &Config{}
			if err := viper.UnmarshalExact(config); err != nil {
				log.Fatal("Nope ", err)
			}
		})
	}
	return config
}

func GetPathRelativeToConfig(p string) string {
	if path.IsAbs(p) {
		return p
	} else if strings.HasPrefix(p, "~") {
		home, err := homedir.Dir()
		if err != nil {
			panic(err)
		}
		return path.Join(home, strings.TrimPrefix(p, "~"))
	} else {
		return path.Join(path.Dir(viper.ConfigFileUsed()), p)
	}
}

func (c Config) CheckConfig() error {
	for name, backend := range c.Backends {
		if err := backend.validate(); err != nil {
			return fmt.Errorf("backend \"%s\": %s", name, err)
		}
	}
	for name, location := range c.Locations {
		if err := location.validate(c); err != nil {
			return fmt.Errorf("location \"%s\": %s", name, err)
		}
	}
	return nil
}

func GetAllOrLocation(cmd *cobra.Command, backends bool) []string {
	var list []string
	if backends {
		for key := range config.Backends {
			list = append(list, key)
		}
	} else {
		for key := range config.Locations {
			list = append(list, key)
		}
	}
	all, _ := cmd.Flags().GetBool("all")
	if all {
		return list
	} else {
		var selected []string
		if backends {
			tmp, _ := cmd.Flags().GetStringSlice("backend")
			selected = tmp
		} else {
			tmp, _ := cmd.Flags().GetStringSlice("location")
			selected = tmp
		}
		for _, s := range selected {
			found := false
			for _, l := range list {
				if l == s {
					found = true
					break
				}
			}
			if !found {
				panic("invalid key")
			}
		}
		return selected
	}
}
