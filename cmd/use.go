package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/katz/ccs/internal/claude"
	"github.com/katz/ccs/internal/config"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:     "use [alias]",
	Aliases: []string{"u"},
	Short:   "Switch to a provider (alias: u)",
	Run:     runUse,
}

func runUse(cmd *cobra.Command, args []string) {
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
			marker := "  "
			if p.Alias == cfg.CurrentProvider {
				marker = "* "
			}
			options[i] = fmt.Sprintf("%s%s (%s)", marker, p.Name, p.Alias)
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

	settings, err := claude.LoadSettings()
	if err != nil {
		color.Red("Failed to load Claude settings: %v", err)
		return
	}

	settings.ClearProviderSettings()
	settings.ApplyProvider(provider)

	if err := settings.Save(); err != nil {
		color.Red("Failed to update Claude settings: %v", err)
		return
	}

	cfg.CurrentProvider = alias
	if err := cfg.Save(); err != nil {
		color.Red("Failed to save config: %v", err)
		return
	}

	color.Green("Switched to '%s'", provider.Name)
}
