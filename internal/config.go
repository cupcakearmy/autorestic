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
	Locations []Location `mapstructure:"locations"`
	Backends  []Backend  `mapstructure:"backends"`
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

func (c *Config) CheckConfig() error {
	if c == nil {
		return fmt.Errorf("config could not be loaded/found")
	}
	if !CheckIfResticIsCallable() {
		return fmt.Errorf(`restic was not found. Install either with "autorestic install" or manually`)
	}
	found := map[string]bool{}
	for _, backend := range c.Backends {
		if err := backend.validate(); err != nil {
			return err
		}
		if _, ok := found[backend.Name]; ok {
			return fmt.Errorf(`duplicate name for backends "%s"`, backend.Name)
		} else {
			found[backend.Name] = true
		}
	}
	found = map[string]bool{}
	for _, location := range c.Locations {
		if err := location.validate(c); err != nil {
			return err
		}
		if _, ok := found[location.Name]; ok {
			return fmt.Errorf(`duplicate name for locations "%s"`, location.Name)
		} else {
			found[location.Name] = true
		}
	}
	return nil
}

func GetAllOrSelected(cmd *cobra.Command, backends bool) ([]string, error) {
	var list []string
	if backends {
		for _, b := range config.Backends {
			list = append(list, b.Name)
		}
	} else {
		for _, l := range config.Locations {
			list = append(list, l.Name)
		}
	}
	all, _ := cmd.Flags().GetBool("all")
	if all {
		return list, nil
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
}

func AddFlagsToCommand(cmd *cobra.Command, backend bool) {
	cmd.PersistentFlags().BoolP("all", "a", false, "Backup all locations")
	if backend {
		cmd.PersistentFlags().StringSliceP("backend", "b", []string{}, "backends")
	} else {
		cmd.PersistentFlags().StringSliceP("location", "l", []string{}, "Locations")
	}
}
