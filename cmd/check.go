package cmd

import (
	"github.com/cupcakearmy/autorestic/internal"
	"github.com/cupcakearmy/autorestic/internal/colors"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if everything is setup",
	Run: func(cmd *cobra.Command, args []string) {
		internal.GetConfig()
		err := internal.Lock()
		CheckErr(err)
		defer internal.Unlock()

		CheckErr(internal.CheckConfig())

		colors.Success.Println("Everything is fine.")
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
