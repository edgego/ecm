package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const exampleStr = `
  To load completions:

  Bash:

  $ source <(ecm completion bash)

  # To load completions for each session, execute once:
  Linux:
    $ ecm completion bash > /etc/bash_completion.d/ecm
  MacOS:
    $ ecm completion bash > /usr/local/etc/bash_completion.d/ecm

  Zsh:

  # If shell completion is not already enabled in your environment you will need
  # to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ ecm completion zsh > "${fpath[1]}/_ecm"

  # You will need to start a new shell for this setup to take effect.

  Fish:

  $ ecm completion fish | source

  # To load completions for each session, execute once:
  $ ecm completion fish > ~/.config/fish/completions/ecm.fish
`

// CompletionCommand used for command completion.
func CompletionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "completion [bash|zsh|fish|powershell]",
		Short:                 "Generate completion script",
		Example:               exampleStr,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				_ = cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				_ = cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				_ = cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				_ = cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
		},
	}

	return cmd
}
