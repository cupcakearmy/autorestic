package cmd

import (
	"fmt"

	"github.com/cupcakearmy/autorestic/internal"
	"github.com/cupcakearmy/autorestic/internal/lock"
	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore [snapshot id]",
	Short: "Restore backup for location",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		internal.GetConfig()
		err := lock.Lock()
		CheckErr(err)
		defer lock.Unlock()

		location, _ := cmd.Flags().GetString("location")
		l, ok := internal.GetLocation(location)
		if !ok {
			CheckErr(fmt.Errorf("invalid location \"%s\"", location))
		}
		target, _ := cmd.Flags().GetString("to")
		from, _ := cmd.Flags().GetString("from")
		force, _ := cmd.Flags().GetBool("force")
		snapshot := ""
		if len(args) > 0 {
			snapshot = args[0]
		}
		err = l.Restore(target, from, force, snapshot)
		CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().BoolP("force", "f", false, "Force, target folder will be overwritten")
	restoreCmd.Flags().String("from", "", "Which backend to use")
	restoreCmd.Flags().String("to", "", "Where to restore the data")
	restoreCmd.Flags().StringP("location", "l", "", "Location to be restored")
	restoreCmd.MarkFlagRequired("location")
}
