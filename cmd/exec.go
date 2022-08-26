package cmd

import (
	"fmt"

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

		var errors []error
		for _, name := range selected {
			colors.PrimaryPrint("  Executing on \"%s\"  ", name)
			backend, _ := internal.GetBackend(name)
			err := backend.Exec(args)
			if err != nil {
				errors = append(errors, err)
			}
		}

		if len(errors) > 0 {
			for _, err := range errors {
				colors.Error.Printf("%s\n\n", err)
			}

			CheckErr(fmt.Errorf("%d errors were found", len(errors)))
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	internal.AddFlagsToCommand(execCmd, true)
}
