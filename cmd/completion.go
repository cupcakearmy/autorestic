package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(autorestic completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ autorestic completion bash > /etc/bash_completion.d/autorestic
  # macOS:
  $ autorestic completion bash > /usr/local/etc/bash_completion.d/autorestic

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ autorestic completion zsh > "${fpath[1]}/_autorestic"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ autorestic completion fish | source

  # To load completions for each session, execute once:
  $ autorestic completion fish > ~/.config/fish/completions/autorestic.fish

PowerShell:

  PS> autorestic completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> autorestic completion powershell > autorestic.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
