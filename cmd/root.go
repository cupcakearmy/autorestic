package cmd

import (
	"os"
	"path/filepath"

	"github.com/cupcakearmy/autorestic/internal"
	"github.com/cupcakearmy/autorestic/internal/colors"
	"github.com/cupcakearmy/autorestic/internal/lock"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func CheckErr(err error) {
	if err != nil {
		colors.Error.Fprintln(os.Stderr, "Error:", err)
		lock.Unlock()
		os.Exit(1)
	}
}

var cfgFile string

var rootCmd = &cobra.Command{
	Version: internal.VERSION,
	Use:     "autorestic",
	Short:   "CLI Wrapper for restic",
	Long:    "Documentation:\thttps://autorestic.vercel.app\nSupport:\thttps://discord.gg/wS7RpYTYd2",
}

func Execute() {
	CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.autorestic.yml or ./.autorestic.yml)")
	rootCmd.PersistentFlags().BoolVar(&internal.CI, "ci", false, "CI mode disabled interactive mode and colors and enables verbosity")
	rootCmd.PersistentFlags().BoolVarP(&internal.VERBOSE, "verbose", "v", false, "verbose mode")
	rootCmd.PersistentFlags().StringVar(&internal.RESTIC_BIN, "restic-bin", "restic", "specify custom restic binary")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if ci, _ := rootCmd.Flags().GetBool("ci"); ci {
		colors.DisableColors(true)
		internal.VERBOSE = true
	}

	if cfgFile != "" {
		if internal.VERBOSE {
			colors.Faint.Printf("> Using config file: %s\n", cfgFile)
		}
		viper.SetConfigFile(cfgFile)
		viper.AutomaticEnv()
		if viper.ConfigFileUsed() == "" {
			colors.Error.Println("cannot read config file %s\n", cfgFile)
			os.Exit(1)
		}
	} else {
		configPaths := []string{"."}

		// Home
		if home, err := homedir.Dir(); err != nil {
			configPaths = append(configPaths, home)
		}

		// XDG_CONFIG_HOME
		{
			prefix, found := os.LookupEnv("XDG_CONFIG_HOME")
			if !found {
				if home, err := homedir.Dir(); err != nil {
					prefix = filepath.Join(home, ".config")
				}
			}
			xdgConfig := filepath.Join(prefix, "autorestic")
			configPaths = append(configPaths, xdgConfig)
		}
		for _, cfgPath := range configPaths {
			if internal.VERBOSE {
				colors.Faint.Printf("> Adding config path: '%s'\n", cfgPath)
			}
			viper.AddConfigPath(cfgPath)
		}
		cfgFileName := ".autorestic"
		viper.SetConfigName(cfgFileName)
		viper.AutomaticEnv()
	}
}
