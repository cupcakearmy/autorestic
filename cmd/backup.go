package cmd

import (
	"github.com/cupcakearmy/autorestic/internal"
	"github.com/cupcakearmy/autorestic/internal/lock"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create backups for given locations",
	Run: func(cmd *cobra.Command, args []string) {
		err := lock.Lock()
		CheckErr(err)
		defer lock.Unlock()

		CheckErr(internal.CheckConfig())

		selected, err := internal.GetAllOrSelected(cmd, false)
		CheckErr(err)
		for _, name := range selected {
			location, _ := internal.GetLocation(name)
			location.Backup(false)
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	internal.AddFlagsToCommand(backupCmd, false)
}
