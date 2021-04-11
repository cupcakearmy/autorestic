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

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute arbitrary native restic commands for given backends",
	Run: func(cmd *cobra.Command, args []string) {
		err := lock.Lock()
		CheckErr(err)
		defer lock.Unlock()

		config := internal.GetConfig()
		err = config.CheckConfig()
		CheckErr(err)

		selected, err := internal.GetAllOrSelected(cmd, true)
		CheckErr(err)
		for _, name := range selected {
			fmt.Println(name)
			backend, _ := internal.GetBackend(name)
			backend.Exec(args)
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	internal.AddFlagsToCommand(execCmd, true)
}
