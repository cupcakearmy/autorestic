package cmd

import (
	"github.com/cupcakearmy/autorestic/internal"
	"github.com/cupcakearmy/autorestic/internal/colors"
	"github.com/cupcakearmy/autorestic/internal/lock"
	"github.com/spf13/cobra"
)

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock autorestic only if you are sure that no other instance is running (ps aux | grep autorestic)",
	Long: `Unlock autorestic only if you are sure that no other instance is running.
To check you can run "ps aux | grep autorestic".`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.GetConfig()
		err := lock.Unlock()
		if err != nil {
			colors.Error.Println("Could not unlock:", err)
			return
		}

		colors.Success.Println("Unlock successful")
	},
}

func init() {
	rootCmd.AddCommand(unlockCmd)
}
