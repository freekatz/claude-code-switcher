package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/katz/ccs/internal/claude"
	"github.com/katz/ccs/internal/config"
	"github.com/katz/ccs/internal/i18n"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:     "use [alias]",
	Aliases: []string{"u"},
	Short:   "Switch to a provider / 切换提供商 (alias: u)",
	Long:    `Switch to a specified provider and update Claude Code settings.`,
	Run:     runUse,
}

func runUse(cmd *cobra.Command, args []string) {
	t := i18n.T()

	cfg, err := config.Load()
	if err != nil {
		color.Red(t.ErrLoadConfig, err)
		return
	}

	if len(cfg.Providers) == 0 {
		color.Yellow(t.MsgNoProviders)
		return
	}

	var alias string
	if len(args) > 0 {
		alias = args[0]
	} else {
		// Let user select provider
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
			Message: t.PromptSelectProvider + ":",
			Options: options,
		}
		if err := survey.AskOne(prompt, &selected); err != nil {
			return
		}
		alias = cfg.Providers[selected].Alias
	}

	provider, err := cfg.GetProvider(alias)
	if err != nil {
		color.Red(t.ErrProviderNotFound, alias)
		return
	}

	// Load Claude Code settings
	settings, err := claude.LoadSettings()
	if err != nil {
		color.Red(t.ErrUpdateSettings, err)
		return
	}

	// Clear existing provider settings first to ensure clean state
	settings.ClearProviderSettings()

	// Apply new provider settings
	settings.ApplyProvider(provider)

	// Save Claude Code settings
	if err := settings.Save(); err != nil {
		color.Red(t.ErrUpdateSettings, err)
		return
	}

	// Update current provider in CCS config
	cfg.CurrentProvider = alias
	if err := cfg.Save(); err != nil {
		color.Red(t.ErrSaveConfig, err)
		return
	}

	color.Green(t.MsgProviderSwitched, provider.Name)
}
