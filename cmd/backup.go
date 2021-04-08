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
	"fmt"

	"github.com/cupcakearmy/autorestic/internal"
	"github.com/cupcakearmy/autorestic/internal/lock"
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		config := internal.GetConfig()
		{
			err := config.CheckConfig()
			cobra.CheckErr(err)
		}
		{
			err := lock.Lock()
			cobra.CheckErr(err)
		}
		defer lock.Unlock()
		{
			backup(internal.GetAllOrLocation(cmd, false), config)
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.PersistentFlags().StringSliceP("location", "l", []string{}, "Locations")
	backupCmd.PersistentFlags().BoolP("all", "a", false, "Backup all locations")
}

func backup(locations []string, config *internal.Config) {
	for _, name := range locations {
		location, ok := config.Locations[name]
		if !ok {
			fmt.Println(fmt.Errorf("location `%s` does not exist", name))
		} else {
			fmt.Printf("Backing up: `%s`", name)
			location.Backup()
		}
	}
}
