package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cseitz-forks/autorestic/internal"
	"github.com/cseitz-forks/autorestic/internal/colors"
	"github.com/cseitz-forks/autorestic/internal/flags"
	"github.com/cseitz-forks/autorestic/internal/lock"
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
	rootCmd.PersistentFlags().BoolVar(&flags.CI, "ci", false, "CI mode disabled interactive mode and colors and enables verbosity")
	rootCmd.PersistentFlags().BoolVarP(&flags.VERBOSE, "verbose", "v", false, "verbose mode")
	rootCmd.PersistentFlags().StringVar(&flags.RESTIC_BIN, "restic-bin", "restic", "specify custom restic binary")
	rootCmd.PersistentFlags().StringVar(&flags.DOCKER_IMAGE, "docker-image", "cseitz-forks/autorestic:"+internal.VERSION, "specify a custom docker image")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if ci, _ := rootCmd.Flags().GetBool("ci"); ci {
		colors.DisableColors(true)
		flags.VERBOSE = true
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.AutomaticEnv()
		if viper.ConfigFileUsed() == "" {
			colors.Error.Printf("cannot read config file %s\n", cfgFile)
			os.Exit(1)
		}
	} else {
		configPaths := getConfigPaths()
		for _, cfgPath := range configPaths {
			viper.AddConfigPath(cfgPath)
		}
		if flags.VERBOSE {
			colors.Faint.Printf("Using config paths: %s\n", strings.Join(configPaths, " "))
		}
		cfgFileName := ".autorestic"
		viper.SetConfigName(cfgFileName)
		viper.AutomaticEnv()
	}
}

func getConfigPaths() []string {
	result := []string{"."}
	if home, err := homedir.Dir(); err == nil {
		result = append(result, home)
	}

	{
		xdgConfigHome, found := os.LookupEnv("XDG_CONFIG_HOME")
		if !found {
			if home, err := homedir.Dir(); err == nil {
				xdgConfigHome = filepath.Join(home, ".config")
			}
		}
		xdgConfig := filepath.Join(xdgConfigHome, "autorestic")
		result = append(result, xdgConfig)
	}
	return result
}
