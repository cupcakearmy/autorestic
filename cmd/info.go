package cmd

import (
	"github.com/cseitz-forks/autorestic/internal"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show info about the config",
	Run: func(cmd *cobra.Command, args []string) {
		internal.GetConfig().Describe()
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
