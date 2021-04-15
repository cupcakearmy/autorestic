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

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore backup for location",
	Run: func(cmd *cobra.Command, args []string) {
		err := lock.Lock()
		CheckErr(err)
		defer lock.Unlock()

		location, _ := cmd.Flags().GetString("location")
		l, ok := internal.GetLocation(location)
		if !ok {
			CheckErr(fmt.Errorf("invalid location \"%s\"", location))
		}
		target, _ := cmd.Flags().GetString("to")
		from, _ := cmd.Flags().GetString("from")
		force, _ := cmd.Flags().GetBool("force")
		err = l.Restore(target, from, force)
		CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().BoolP("force", "f", false, "Force, target folder will be overwritten")
	restoreCmd.Flags().String("from", "", "Which backend to use")
	restoreCmd.Flags().String("to", "", "Where to restore the data")
	// restoreCmd.MarkFlagRequired("to")
	restoreCmd.Flags().StringP("location", "l", "", "Location to be restored")
	restoreCmd.MarkFlagRequired("location")
}
