package cmd

import (
	"github.com/cseitz-forks/autorestic/internal/bins"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall restic and autorestic",
	Run: func(cmd *cobra.Command, args []string) {
		restic, _ := cmd.Flags().GetBool("restic")
		bins.Uninstall(restic)
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().Bool("restic", false, "also uninstall restic.")
}
