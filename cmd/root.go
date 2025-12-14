package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

// Build info variables - can be set via ldflags during build
var (
	version   = "dev"
	buildTime = "unknown"
	gitBranch = "unknown"
	gitCommit = "unknown"
)

var name = "Claude Code Switcher"

var rootCmd = &cobra.Command{
	Use:   "ccs",
	Short: name,
	Long:  `Claude Code Switcher - Manage multiple Claude Code API providers`,
}

func init() {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(fmt.Sprintf(`%s version %s
Built:      %s
Git branch: %s
Git commit: %s
Go version: %s
OS/Arch:    %s/%s
`, name, version, buildTime, gitBranch, gitCommit, runtime.Version(), runtime.GOOS, runtime.GOARCH))

	// Hide completion command
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	// Add help command with alias
	helpCmd := &cobra.Command{
		Use:     "help [command]",
		Aliases: []string{"h"},
		Short:   "Help about any command (alias: h)",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				rootCmd.Help()
				return
			}
			for _, c := range rootCmd.Commands() {
				if c.Name() == args[0] || contains(c.Aliases, args[0]) {
					c.Help()
					return
				}
			}
			rootCmd.Help()
		},
	}
	rootCmd.SetHelpCommand(helpCmd)

	// Add all subcommands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(removeCmd)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
