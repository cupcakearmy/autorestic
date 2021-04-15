/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

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
	Long:    "Documentation: https://autorestic.vercel.app",
}

func Execute() {
	CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.autorestic.yml or ./.autorestic.yml)")
	rootCmd.PersistentFlags().BoolVar(&internal.CI, "ci", false, "CI mode disabled interactive mode and colors and enables verbosity")
	rootCmd.PersistentFlags().BoolVar(&internal.VERBOSE, "verbose", false, "verbose mode")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if ci, _ := rootCmd.Flags().GetBool("ci"); ci {
		colors.DisableColors(true)
		internal.VERBOSE = true
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		CheckErr(err)

		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigName(".autorestic")
	}
	viper.AutomaticEnv()
	internal.GetConfig()
}
