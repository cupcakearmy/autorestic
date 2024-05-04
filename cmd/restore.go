package cmd

import (
	"fmt"

	"github.com/cseitz-forks/autorestic/internal"
	"github.com/cseitz-forks/autorestic/internal/lock"
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

		// Get optional flags
		optional := []string{}
		for _, flag := range []string{"include", "exclude", "iinclude", "iexclude"} {
			values, err := cmd.Flags().GetStringSlice(flag)
			if err == nil {
				for _, value := range values {
					optional = append(optional, "--"+flag, value)
				}
			}
		}

		err = l.Restore(target, from, force, snapshot, optional)
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

	// Passed on flags
	restoreCmd.Flags().StringSliceP("include", "i", []string{}, "Include a pattern")
	restoreCmd.Flags().StringSliceP("exclude", "e", []string{}, "Exclude a pattern")
	restoreCmd.Flags().StringSlice("iinclude", []string{}, "Include a pattern, case insensitive")
	restoreCmd.Flags().StringSlice("iexclude", []string{}, "Exclude a pattern, case insensitive")
}
