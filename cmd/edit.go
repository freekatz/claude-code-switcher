package cmd

import (
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/katz/ccs/internal/claude"
	"github.com/katz/ccs/internal/config"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:     "edit [alias]",
	Aliases: []string{"e"},
	Short:   "Edit a provider (alias: e)",
	Run:     runEdit,
}

func runEdit(cmd *cobra.Command, args []string) {
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

	fields := []string{
		"All fields",
		"Name",
		"Alias",
		"Base URL",
		"API Key",
		"Model",
		"Small model",
		"Sonnet model",
		"Opus model",
		"Haiku model",
		"Timeout",
	}

	var selectedField int
	fieldPrompt := &survey.Select{
		Message: "Select field:",
		Options: fields,
	}
	if err := survey.AskOne(fieldPrompt, &selectedField); err != nil {
		return
	}

	updated := *provider

	if selectedField == 0 {
		editAllFields(&updated)
	} else {
		editField(&updated, selectedField)
	}

	isCurrentProvider := (cfg.CurrentProvider == alias) || (cfg.CurrentProvider == updated.Alias)

	if err := cfg.UpdateProvider(alias, updated); err != nil {
		if err == config.ErrProviderExists {
			color.Red("Provider '%s' already exists", updated.Alias)
		} else {
			color.Red("Failed to update: %v", err)
		}
		return
	}

	if err := cfg.Save(); err != nil {
		color.Red("Failed to save: %v", err)
		return
	}

	if isCurrentProvider {
		if err := updateClaudeSettings(&updated); err != nil {
			color.Yellow("Warning: Failed to update Claude settings: %v", err)
		}
	}

	color.Green("Provider '%s' updated", updated.Name)
}

func editAllFields(p *config.Provider) {
	survey.AskOne(&survey.Input{Message: "Name:", Default: p.Name}, &p.Name)
	survey.AskOne(&survey.Input{Message: "Alias:", Default: p.Alias}, &p.Alias)
	survey.AskOne(&survey.Input{Message: "Base URL:", Default: p.BaseURL}, &p.BaseURL)

	var apiKey string
	survey.AskOne(&survey.Password{Message: "API Key (empty=keep):"}, &apiKey)
	if apiKey != "" {
		p.APIKey = apiKey
	}

	survey.AskOne(&survey.Input{Message: "Model:", Default: p.Model}, &p.Model)
	survey.AskOne(&survey.Input{Message: "Small model:", Default: p.SmallModel}, &p.SmallModel)
	survey.AskOne(&survey.Input{Message: "Sonnet model:", Default: p.SonnetModel}, &p.SonnetModel)
	survey.AskOne(&survey.Input{Message: "Opus model:", Default: p.OpusModel}, &p.OpusModel)
	survey.AskOne(&survey.Input{Message: "Haiku model:", Default: p.HaikuModel}, &p.HaikuModel)

	var timeoutStr string
	survey.AskOne(&survey.Input{Message: "Timeout ms:", Default: strconv.Itoa(p.Timeout)}, &timeoutStr)
	if timeout, err := strconv.Atoi(timeoutStr); err == nil {
		p.Timeout = timeout
	}
}

func editField(p *config.Provider, fieldIndex int) {
	switch fieldIndex {
	case 1:
		survey.AskOne(&survey.Input{Message: "Name:", Default: p.Name}, &p.Name)
	case 2:
		survey.AskOne(&survey.Input{Message: "Alias:", Default: p.Alias}, &p.Alias)
	case 3:
		survey.AskOne(&survey.Input{Message: "Base URL:", Default: p.BaseURL}, &p.BaseURL)
	case 4:
		var apiKey string
		survey.AskOne(&survey.Password{Message: "API Key:"}, &apiKey)
		if apiKey != "" {
			p.APIKey = apiKey
		}
	case 5:
		survey.AskOne(&survey.Input{Message: "Model:", Default: p.Model}, &p.Model)
	case 6:
		survey.AskOne(&survey.Input{Message: "Small model:", Default: p.SmallModel}, &p.SmallModel)
	case 7:
		survey.AskOne(&survey.Input{Message: "Sonnet model:", Default: p.SonnetModel}, &p.SonnetModel)
	case 8:
		survey.AskOne(&survey.Input{Message: "Opus model:", Default: p.OpusModel}, &p.OpusModel)
	case 9:
		survey.AskOne(&survey.Input{Message: "Haiku model:", Default: p.HaikuModel}, &p.HaikuModel)
	case 10:
		var timeoutStr string
		survey.AskOne(&survey.Input{Message: "Timeout ms:", Default: strconv.Itoa(p.Timeout)}, &timeoutStr)
		if timeout, err := strconv.Atoi(timeoutStr); err == nil {
			p.Timeout = timeout
		}
	}
}

func updateClaudeSettings(p *config.Provider) error {
	settings, err := claude.LoadSettings()
	if err != nil {
		return err
	}
	settings.ClearProviderSettings()
	settings.ApplyProvider(p)
	return settings.Save()
}
