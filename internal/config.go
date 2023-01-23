package internal

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cupcakearmy/autorestic/internal/colors"
	"github.com/cupcakearmy/autorestic/internal/flags"
	"github.com/cupcakearmy/autorestic/internal/lock"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const VERSION = "1.7.5"

type OptionMap map[string][]interface{}
type Options map[string]OptionMap

type Config struct {
	Version   string              `mapstructure:"version"`
	Extras    interface{}         `mapstructure:"extras"`
	Locations map[string]Location `mapstructure:"locations"`
	Backends  map[string]Backend  `mapstructure:"backends"`
	Global    Options             `mapstructure:"global"`
}

var once sync.Once
var config *Config

func exitConfig(err error, msg string) {
	if err != nil {
		colors.Error.Println(err)
	}
	if msg != "" {
		colors.Error.Println(msg)
	}
	lock.Unlock()
	os.Exit(1)
}

func GetDockerVolumesLocations() (map[string]Location, error) {
	var (
		cli              *client.Client
		err              error
		volumeListOKBody volume.VolumeListOKBody
		ctx              context.Context
	)

	volumesList := make(map[string]Location)

	if flags.DOCKER_HOST != "" {
		cli, err = client.NewClientWithOpts(client.WithHost(flags.DOCKER_HOST))
	} else {
		cli, err = client.NewClientWithOpts()
	}
	if err != nil {
		return volumesList, fmt.Errorf("Could not create client: %s", err.Error())
	}
	defer cli.Close()

	ctx = context.Background()
	volumeListOKBody, err = cli.VolumeList(ctx, filters.NewArgs())
	if err != nil {
		return volumesList, fmt.Errorf("Could not list volumes: %s", err.Error())
	}

	for _, volume := range volumeListOKBody.Volumes {
		val, ok := volume.Labels["autorestic.enable"]
		if ok && (val == "1" || val == "true") {
			volumeTo, ok := volume.Labels["autorestic.to"]
			if !ok {
				return volumesList, fmt.Errorf("Autorestic is enabled but there is no \"autorestic.to\" label for volume %s", volume.Name)
			}
			volumeLocation := Location{
				name: volume.Name,
				From: []string{volume.Name},
				Type: "volume",
				To: strings.Split(volumeTo, ";"),
			}
			for label, value := range volume.Labels {
				if label == "autorestic.cron" {
					volumeLocation.Cron = value
				} else if label == "autorestic.forget" {
					volumeLocation.ForgetOption = LocationForgetOption(value)
				} else if label == "autorestic.hooks.dir" {
					volumeLocation.Hooks.Dir = value
				} else if label == "autorestic.hooks.before" {
					volumeLocation.Hooks.Before = []string{value}
				} else if label == "autorestic.hooks.after" {
					volumeLocation.Hooks.After = []string{value}
				} else if label == "autorestic.hooks.success" {
					volumeLocation.Hooks.Success = []string{value}
				} else if label == "autorestic.hooks.failure" {
					volumeLocation.Hooks.Failure = []string{value}
				} else if strings.HasPrefix(label, "autorestic.options.") {
					label = strings.TrimPrefix(label, "autorestic.options.")
					target := strings.Split(label, ".")
					if len(target) != 2 {
						fmt.Errorf("Failed to parse label \"%s\" for volume %s", volume.Name)
					}

					values := strings.Split(value, ";")
					interfaceValues := make([]interface{}, 0)
					for _, v := range values {
						interfaceValues = append(interfaceValues, interface{}(v))
					}

					if volumeLocation.Options == nil {
						volumeLocation.Options = make(Options)
					}
					volumeLocation.Options[target[0]] = make(OptionMap)
					volumeLocation.Options[target[0]][target[1]] = interfaceValues
				} else if strings.HasPrefix(label, "autorestic.copy.") {
					fromBackend := strings.TrimPrefix(label, "autorestic.copy.")
					if volumeLocation.CopyOption == nil {
						volumeLocation.CopyOption = make(LocationCopy)
					}
					volumeLocation.CopyOption[fromBackend] = strings.Split(value, ";")
				}
			}
			volumesList[volume.Name] = volumeLocation
		}
	}

	return volumesList, nil
}

func GetConfig() *Config {
	if config == nil {
		once.Do(func() {
			if err := viper.ReadInConfig(); err == nil {
				absConfig, _ := filepath.Abs(viper.ConfigFileUsed())
				if !flags.CRON_LEAN {
					colors.Faint.Println("Using config: \t", absConfig)
				}
				// Load env file
				envFile := filepath.Join(filepath.Dir(absConfig), ".autorestic.env")
				err = godotenv.Load(envFile)
				if err == nil && !flags.CRON_LEAN {
					colors.Faint.Println("Using env:\t", envFile)
				}
			} else {
				text := err.Error()
				if strings.Contains(text, "no such file or directory") {
					cfgFileName := ".autorestic"
					colors.Error.Println(
						fmt.Sprintf(
							"cannot find configuration file '%s.yml' or '%s.yaml'.",
							cfgFileName, cfgFileName))
				} else {
					colors.Error.Println("could not load config file\n" + text)
				}
				os.Exit(1)
			}

			var versionConfig interface{}
			viper.UnmarshalKey("version", &versionConfig)
			if versionConfig == nil {
				exitConfig(nil, "no version specified in config file. please see docs on how to migrate")
			}
			version, ok := versionConfig.(int)
			if !ok {
				exitConfig(nil, "version specified in config file is not an int")
			} else {
				// Check for version
				if version != 2 {
					exitConfig(nil, "unsupported config version number. please check the docs for migration\nhttps://autorestic.vercel.app/migration/")
				}
			}

			config = &Config{}
			if err := viper.UnmarshalExact(config); err != nil {
				exitConfig(err, "Could not parse config file!")
			}

			if flags.DOCKER_DISCOVERY {
				if config.Locations == nil {
					config.Locations = make(map[string]Location)
				}
				volumeLocations, err := GetDockerVolumesLocations()
				if err != nil {
					exitConfig(err, "failed to load locations from docker")
				}
				for name, location := range volumeLocations {
					config.Locations[name] = location
				}
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

func (c *Config) Describe() {
	// Locations
	for name, l := range c.Locations {
		var tmp string
		colors.PrimaryPrint(`Location: "%s"`, name)

		tmp = ""
		for _, path := range l.From {
			tmp += fmt.Sprintf("\t%s %s\n", colors.Success.Sprint("←"), path)
		}
		colors.PrintDescription("From", tmp)

		tmp = ""
		for _, to := range l.To {
			tmp += fmt.Sprintf("\t%s %s\n", colors.Success.Sprint("→"), to)
		}
		colors.PrintDescription("To", tmp)

		if l.Cron != "" {
			colors.PrintDescription("Cron", l.Cron)
		}

		tmp = ""
		hooks := map[string][]string{
			"Before":  l.Hooks.Before,
			"After":   l.Hooks.After,
			"Failure": l.Hooks.Failure,
			"Success": l.Hooks.Success,
		}
		for hook, commands := range hooks {
			if len(commands) > 0 {
				tmp += "\n\t" + hook
				for _, cmd := range commands {
					tmp += colors.Faint.Sprintf("\n\t  ▶ %s", cmd)
				}
			}
		}
		if tmp != "" {
			colors.PrintDescription("Hooks", tmp)
		}

		if len(l.Options) > 0 {
			tmp = ""
			for t, options := range l.Options {
				tmp += "\n\t" + t
				for option, values := range options {
					for _, value := range values {
						tmp += colors.Faint.Sprintf("\n\t  ✧ --%s=%s", option, value)
					}
				}
			}
			colors.PrintDescription("Options", tmp)
		}
	}

	// Backends
	for name, b := range c.Backends {
		colors.PrimaryPrint("Backend: \"%s\"", name)
		colors.PrintDescription("Type", b.Type)
		colors.PrintDescription("Path", b.Path)

		if len(b.Env) > 0 {
			tmp := ""
			for option, value := range b.Env {
				tmp += fmt.Sprintf("\n\t%s %s %s", colors.Success.Sprint("✧"), strings.ToUpper(option), colors.Faint.Sprint(value))
			}
			colors.PrintDescription("Env", tmp)
		}
	}
}

func CheckConfig() error {
	c := GetConfig()
	if c == nil {
		return fmt.Errorf("config could not be loaded/found")
	}
	if !CheckIfResticIsCallable() {
		return fmt.Errorf(`%s was not found. Install either with "autorestic install" or manually`, flags.RESTIC_BIN)
	}
	for name, backend := range c.Backends {
		backend.name = name
		if err := backend.validate(); err != nil {
			return err
		}
	}
	for name, location := range c.Locations {
		location.name = name
		if err := location.validate(); err != nil {
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
		var splitted = strings.Split(s, "@")
		for _, l := range list {
			if l == splitted[0] {
				goto found
			}
		}
		if backends {
			return nil, fmt.Errorf("invalid backend \"%s\"", s)
		} else {
			return nil, fmt.Errorf("invalid location \"%s\"", s)
		}
	found:
	}

	if len(selected) == 0 {
		return selected, fmt.Errorf("nothing selected, aborting")
	}
	return selected, nil
}

func AddFlagsToCommand(cmd *cobra.Command, backend bool) {
	var usage string
	if backend {
		usage = "all backends"
	} else {
		usage = "all locations"
	}
	cmd.PersistentFlags().BoolP("all", "a", false, usage)
	if backend {
		cmd.PersistentFlags().StringSliceP("backend", "b", []string{}, "select backends")
	} else {
		cmd.PersistentFlags().StringSliceP("location", "l", []string{}, "select locations")
	}
}

func (c *Config) SaveConfig() error {
	file := viper.ConfigFileUsed()
	if err := CopyFile(file, file+".old"); err != nil {
		return err
	}
	colors.Secondary.Println("Saved a backup copy of your file next to the original.")

	viper.Set("backends", c.Backends)
	viper.Set("locations", c.Locations)

	return viper.WriteConfig()
}

func optionToString(option string) string {
	if !strings.HasPrefix(option, "-") {
		return "--" + option
	}
	return option
}

func appendOptionsToSlice(str *[]string, options OptionMap) {
	for key, values := range options {
		for _, value := range values {
			// Bool
			asBool, ok := value.(bool)
			if ok && asBool {
				*str = append(*str, optionToString(key))
				continue
			}
			*str = append(*str, optionToString(key), fmt.Sprint(value))
		}
	}
}

func getOptions(options Options, keys []string) []string {
	var selected []string
	for _, key := range keys {
		appendOptionsToSlice(&selected, options[key])
	}
	return selected
}

func combineBackendOptions(key string, b Backend) []string {
	// Priority: backend > global
	var options []string
	gFlags := getOptions(GetConfig().Global, []string{key})
	bFlags := getOptions(b.Options, []string{"all", key})
	options = append(options, gFlags...)
	options = append(options, bFlags...)
	return options
}

func combineAllOptions(key string, l Location, b Backend) []string {
	// Priority: location > backend > global
	var options []string
	gFlags := getOptions(GetConfig().Global, []string{key})
	bFlags := getOptions(b.Options, []string{"all", key})
	lFlags := getOptions(l.Options, []string{"all", key})
	options = append(options, gFlags...)
	options = append(options, bFlags...)
	options = append(options, lFlags...)
	return options
}
