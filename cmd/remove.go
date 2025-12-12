package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/katz/ccs/internal/config"
	"github.com/katz/ccs/internal/i18n"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove [alias]",
	Aliases: []string{"rm"},
	Short:   "Remove a provider / 删除提供商 (alias: rm)",
	Long:    `Remove a provider from the configuration.`,
	Run:     runRemove,
}

func runRemove(cmd *cobra.Command, args []string) {
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
			options[i] = fmt.Sprintf("%s (%s)", p.Name, p.Alias)
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

	// Confirm deletion
	var confirm bool
	confirmPrompt := &survey.Confirm{
		Message: fmt.Sprintf(t.PromptConfirmDelete, provider.Name),
		Default: false,
	}
	if err := survey.AskOne(confirmPrompt, &confirm); err != nil {
		return
	}

	if !confirm {
		return
	}

	// Remove provider
	if err := cfg.RemoveProvider(alias); err != nil {
		color.Red(t.ErrProviderNotFound, alias)
		return
	}

	if err := cfg.Save(); err != nil {
		color.Red(t.ErrSaveConfig, err)
		return
	}

	color.Green(t.MsgProviderRemoved, provider.Name)
}
