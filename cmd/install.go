package cmd

import (
	"github.com/cseitz-forks/autorestic/internal/bins"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install restic if missing",
	Run: func(cmd *cobra.Command, args []string) {
		err := bins.InstallRestic()
		CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
