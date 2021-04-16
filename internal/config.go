package internal

import (
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/cupcakearmy/autorestic/internal/colors"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const VERSION = "1.0.0"

var CI bool = false
var VERBOSE bool = false

type Config struct {
	Locations map[string]Location `mapstructure:"locations"`
	Backends  map[string]Backend  `mapstructure:"backends"`
}

var once sync.Once
var config *Config

func GetConfig() *Config {
	if config == nil {
		once.Do(func() {
			if err := viper.ReadInConfig(); err == nil {
				colors.Faint.Println("Using config file:", viper.ConfigFileUsed())
			} else {
				return
			}

			config = &Config{}
			if err := viper.UnmarshalExact(config); err != nil {
				panic(err)
			}
		})
	}
	return config
}

func GetPathRelativeToConfig(p string) (string, error) {
	if path.IsAbs(p) {
		return p, nil
	} else if strings.HasPrefix(p, "~") {
		home, err := homedir.Dir()
		return path.Join(home, strings.TrimPrefix(p, "~")), err
	} else {
		return path.Join(path.Dir(viper.ConfigFileUsed()), p), nil
	}
}

func CheckConfig() error {
	c := GetConfig()
	if c == nil {
		return fmt.Errorf("config could not be loaded/found")
	}
	if !CheckIfResticIsCallable() {
		return fmt.Errorf(`restic was not found. Install either with "autorestic install" or manually`)
	}
	for name, backend := range c.Backends {
		backend.name = name
		if err := backend.validate(); err != nil {
			return err
		}
	}
	for name, location := range c.Locations {
		location.name = name
		if err := location.validate(c); err != nil {
			return err
		}
	}
	return nil
}

func GetAllOrSelected(cmd *cobra.Command, backends bool) ([]string, error) {
	var list []string
	if backends {
		for name := range config.Backends {
			list = append(list, name)
		}
	} else {
		for name := range config.Locations {
			list = append(list, name)
		}
	}

	all, _ := cmd.Flags().GetBool("all")
	if all {
		return list, nil
	}

	var selected []string
	if backends {
		selected, _ = cmd.Flags().GetStringSlice("backend")
	} else {
		selected, _ = cmd.Flags().GetStringSlice("location")
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
			if backends {
				return nil, fmt.Errorf("invalid backend \"%s\"", s)
			} else {
				return nil, fmt.Errorf("invalid location \"%s\"", s)
			}
		}
	}

	if len(selected) == 0 {
		return selected, fmt.Errorf("nothing selected, aborting")
	}
	return selected, nil
}

func AddFlagsToCommand(cmd *cobra.Command, backend bool) {
	cmd.PersistentFlags().BoolP("all", "a", false, "Backup all locations")
	if backend {
		cmd.PersistentFlags().StringSliceP("backend", "b", []string{}, "backends")
	} else {
		cmd.PersistentFlags().StringSliceP("location", "l", []string{}, "Locations")
	}
}
