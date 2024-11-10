package cmd

import (
	"fmt"
	"strings"

	"github.com/cupcakearmy/autorestic/internal"
	"github.com/cupcakearmy/autorestic/internal/colors"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create backups for given locations",
	Run: func(cmd *cobra.Command, args []string) {
		internal.GetConfig()
		err := internal.Lock()
		CheckErr(err)
		defer internal.Unlock()

		selected, err := internal.GetAllOrSelected(cmd, false)
		CheckErr(err)
		errors := 0
		for _, name := range selected {
			var splitted = strings.Split(name, "@")
			var specificBackend = ""
			if len(splitted) > 1 {
				specificBackend = splitted[1]
			}
			location, _ := internal.GetLocation(splitted[0])
			errs := location.Backup(false, specificBackend)
			for _, err := range errs {
				colors.Error.Printf("%s\n\n", err)
				errors++
			}
		}
		if errors > 0 {
			CheckErr(fmt.Errorf("%d errors were found", errors))
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	internal.AddFlagsToCommand(backupCmd, false)
}
