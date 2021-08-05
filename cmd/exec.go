package cmd

import (
	"github.com/cupcakearmy/autorestic/internal"
	"github.com/cupcakearmy/autorestic/internal/colors"
	"github.com/cupcakearmy/autorestic/internal/lock"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute arbitrary native restic commands for given backends",
	Run: func(cmd *cobra.Command, args []string) {
		internal.GetConfig()
		err := lock.Lock()
		CheckErr(err)
		defer lock.Unlock()

		selected, err := internal.GetAllOrSelected(cmd, true)
		CheckErr(err)
		for _, name := range selected {
			colors.PrimaryPrint("  Executing on \"%s\"  ", name)
			backend, _ := internal.GetBackend(name)
			backend.Exec(args)
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	internal.AddFlagsToCommand(execCmd, true)
}
