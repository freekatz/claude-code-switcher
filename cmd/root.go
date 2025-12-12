package cmd

import (
	"fmt"
	"os"

	"github.com/katz/ccs/internal/config"
	"github.com/katz/ccs/internal/i18n"
	"github.com/spf13/cobra"
)

// version can be set via ldflags during build
var version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "ccs",
	Short: "Claude Code Switcher",
	Long:  `Claude Code Switcher - Manage multiple Claude Code API providers`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load config and set language before any command runs
		cfg, err := config.Load()
		if err == nil && cfg.Language != "" {
			i18n.SetLanguage(i18n.Language(cfg.Language))
		}
	},
}

func init() {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(fmt.Sprintf("{{.Name}} version %s\n", version))

	// Add all subcommands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(langCmd)
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
