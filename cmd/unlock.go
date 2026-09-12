package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cupcakearmy/autorestic/internal"
	"github.com/cupcakearmy/autorestic/internal/colors"
	"github.com/cupcakearmy/autorestic/internal/lock"
	"github.com/spf13/cobra"
)

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock autorestic only if you are sure that no other instance is running",
	Long: `Unlock autorestic only if you are sure that no other instance is running.
To check you can run "ps aux | grep '[a]utorestic'".`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.GetConfig()

		force, _ := cmd.Flags().GetBool("force")

		if !force && isAutoresticRunning() {
			colors.Error.Print("Another autorestic instance is running. Are you sure you want to unlock? (yes/no): ")
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "yes" {
				colors.Primary.Println("Unlocking aborted.")
				return
			}
		}

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
	unlockCmd.Flags().Bool("force", false, "force unlock")
}

// isAutoresticRunning checks if an autorestic binary is running and returns
// true if it is. It also prints matching processes to stdout.
func isAutoresticRunning() bool {
	out, err := exec.Command("ps", "aux").Output()
	if err != nil {
		return false
	}

	lines := strings.Split(string(out), "\n")
	autoresticProcesses := []string{}
	currentPid := fmt.Sprint(os.Getpid())

	for _, line := range lines {
		if !psLineLooksLikeAutorestic(line, currentPid) {
			continue
		}
		autoresticProcesses = append(autoresticProcesses, line)
	}

	if len(autoresticProcesses) > 0 {
		colors.Faint.Println("Found autorestic processes:")
		for _, proc := range autoresticProcesses {
			colors.Faint.Println(proc)
		}
		return true
	}
	return false
}

// psLineLooksLikeAutorestic reports whether a `ps aux` line is a running
// autorestic binary. Substring matches on unit names like
// `systemctl start autorestic-backup.service` are ignored (see #406).
func psLineLooksLikeAutorestic(line, currentPid string) bool {
	fields := strings.Fields(line)
	// USER PID %CPU %MEM VSZ RSS TTY STAT START TIME COMMAND...
	if len(fields) < 11 {
		return false
	}
	if fields[1] == currentPid {
		return false
	}

	for _, tok := range fields[10:] {
		base := filepath.Base(tok)
		if base == "autorestic" || base == "autorestic.exe" {
			return true
		}
	}
	return false
}
