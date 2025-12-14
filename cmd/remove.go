package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/katz/ccs/internal/config"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove [alias]",
	Aliases: []string{"rm"},
	Short:   "Remove a provider (alias: rm)",
	Run:     runRemove,
}

func runRemove(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		color.Red("Failed to load config: %v", err)
		return
	}

	if len(cfg.Providers) == 0 {
		color.Yellow("No providers configured")
		return
	}

	var alias string
	if len(args) > 0 {
		alias = args[0]
	} else {
		options := make([]string, len(cfg.Providers))
		for i, p := range cfg.Providers {
			options[i] = fmt.Sprintf("%s (%s)", p.Name, p.Alias)
		}

		var selected int
		prompt := &survey.Select{
			Message: "Select provider:",
			Options: options,
		}
		if err := survey.AskOne(prompt, &selected); err != nil {
			return
		}
		alias = cfg.Providers[selected].Alias
	}

	provider, err := cfg.GetProvider(alias)
	if err != nil {
		color.Red("Provider '%s' not found", alias)
		return
	}

	var confirm bool
	confirmPrompt := &survey.Confirm{
		Message: fmt.Sprintf("Delete '%s'?", provider.Name),
		Default: false,
	}
	if err := survey.AskOne(confirmPrompt, &confirm); err != nil {
		return
	}

	if !confirm {
		return
	}

	if err := cfg.RemoveProvider(alias); err != nil {
		color.Red("Provider '%s' not found", alias)
		return
	}

	if err := cfg.Save(); err != nil {
		color.Red("Failed to save: %v", err)
		return
	}

	color.Green("Provider '%s' removed", provider.Name)
}
