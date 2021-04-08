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
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		config := internal.GetConfig()
		if err := config.CheckConfig(); err != nil {
			panic(err)
		}
		exec(internal.GetAllOrLocation(cmd, true), config, args)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.PersistentFlags().StringSliceP("backend", "b", []string{}, "backends")
	execCmd.PersistentFlags().BoolP("all", "a", false, "Exec in all backends")
}

func exec(backends []string, config *internal.Config, args []string) {
	for _, name := range backends {
		fmt.Println(name)
		backend := config.Backends[name]
		backend.Exec(args)
	}
}
