package cmd

import (
	"github.com/cupcakearmy/autorestic/internal/bins"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade autorestic and restic",
	Run: func(cmd *cobra.Command, args []string) {
		noRestic, _ := cmd.Flags().GetBool("no-restic")
		err := bins.Upgrade(!noRestic)
		CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().Bool("no-restic", false, "also update restic")
}
