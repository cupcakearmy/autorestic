package cmd

import (
	"github.com/cupcakearmy/autorestic/internal"
	"github.com/cupcakearmy/autorestic/internal/colors"
	"github.com/cupcakearmy/autorestic/internal/lock"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if everything is setup",
	Run: func(cmd *cobra.Command, args []string) {
		err := lock.Lock()
		CheckErr(err)
		defer lock.Unlock()

		CheckErr(internal.CheckConfig())

		colors.Success.Println("Everything is fine.")
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
