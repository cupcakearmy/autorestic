package cmd

import (
	"github.com/cseitz-forks/autorestic/internal/bins"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade autorestic and restic",
	Run: func(cmd *cobra.Command, args []string) {
		restic, _ := cmd.Flags().GetBool("restic")
		err := bins.Upgrade(restic)
		CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().Bool("restic", true, "also update restic")
}
