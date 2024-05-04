package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cseitz-forks/autorestic/internal"
	"github.com/cseitz-forks/autorestic/internal/colors"
	"github.com/cseitz-forks/autorestic/internal/lock"
	"github.com/spf13/cobra"
)

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock autorestic only if you are sure that no other instance is running",
	Long: `Unlock autorestic only if you are sure that no other instance is running.
To check you can run "ps aux | grep autorestic".`,
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

// isAutoresticRunning checks if autorestic is running
// and returns true if it is.
// It also prints the processes to stdout.
func isAutoresticRunning() bool {
	cmd := exec.Command("sh", "-c", "ps aux | grep autorestic")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false
	}

	lines := strings.Split(out.String(), "\n")
	autoresticProcesses := []string{}
	currentPid := fmt.Sprint(os.Getpid())

	for _, line := range lines {
		if strings.Contains(line, "autorestic") && !strings.Contains(line, "grep autorestic") && !strings.Contains(line, currentPid) {
			autoresticProcesses = append(autoresticProcesses, line)
		}
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
