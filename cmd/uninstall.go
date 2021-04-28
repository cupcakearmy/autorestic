package cmd

import (
	"github.com/cupcakearmy/autorestic/internal/bins"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall restic and autorestic",
	Run: func(cmd *cobra.Command, args []string) {
		noRestic, _ := cmd.Flags().GetBool("no-restic")
		bins.Uninstall(!noRestic)
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().Bool("no-restic", false, "do not uninstall restic.")
}
