package cmd

import (
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/katz/ccs/internal/claude"
	"github.com/katz/ccs/internal/config"
	"github.com/katz/ccs/internal/i18n"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:     "edit [alias]",
	Aliases: []string{"e"},
	Short:   "Edit a provider / 编辑提供商配置 (alias: e)",
	Long:    `Edit an existing provider configuration.`,
	Run:     runEdit,
}

func runEdit(cmd *cobra.Command, args []string) {
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

	// Select field to edit
	fields := []string{
		t.LabelAll,
		t.FieldName,
		t.FieldAlias,
		t.FieldBaseURL,
		t.FieldAPIKey,
		t.FieldModel,
		t.FieldSmallModel,
		t.FieldSonnetModel,
		t.FieldOpusModel,
		t.FieldHaikuModel,
		t.FieldTimeout,
	}

	var selectedField int
	fieldPrompt := &survey.Select{
		Message: t.PromptSelectField + ":",
		Options: fields,
	}
	if err := survey.AskOne(fieldPrompt, &selectedField); err != nil {
		return
	}

	// Edit selected field(s)
	updated := *provider

	if selectedField == 0 {
		// Edit all fields
		editAllFields(&updated, t)
	} else {
		// Edit specific field
		editField(&updated, selectedField, t)
	}

	// Check if this is the current provider
	isCurrentProvider := (cfg.CurrentProvider == alias) || (cfg.CurrentProvider == updated.Alias)

	// Update provider
	if err := cfg.UpdateProvider(alias, updated); err != nil {
		if err == config.ErrProviderExists {
			color.Red(t.ErrProviderExists, updated.Alias)
		} else {
			color.Red(t.ErrSaveConfig, err)
		}
		return
	}

	if err := cfg.Save(); err != nil {
		color.Red(t.ErrSaveConfig, err)
		return
	}

	// If this is the current provider, update Claude Code settings immediately
	if isCurrentProvider {
		if err := updateClaudeSettings(&updated, t); err != nil {
			color.Yellow("Warning: Provider updated but failed to update Claude Code settings: %v", err)
		}
	}

	color.Green(t.MsgProviderUpdated, updated.Name)
}

func editAllFields(p *config.Provider, t i18n.Messages) {
	survey.AskOne(&survey.Input{Message: t.PromptName + ":", Default: p.Name}, &p.Name)
	survey.AskOne(&survey.Input{Message: t.PromptAlias + ":", Default: p.Alias}, &p.Alias)
	survey.AskOne(&survey.Input{Message: t.PromptBaseURL + ":", Default: p.BaseURL}, &p.BaseURL)

	var apiKey string
	survey.AskOne(&survey.Password{Message: t.PromptAPIKey + " (leave empty to keep current):"}, &apiKey)
	if apiKey != "" {
		p.APIKey = apiKey
	}

	survey.AskOne(&survey.Input{Message: t.PromptModel + " (clear to empty):", Default: p.Model}, &p.Model)
	survey.AskOne(&survey.Input{Message: t.PromptSmallModel + " (clear to empty):", Default: p.SmallModel}, &p.SmallModel)
	survey.AskOne(&survey.Input{Message: t.PromptSonnetModel + " (clear to empty):", Default: p.SonnetModel}, &p.SonnetModel)
	survey.AskOne(&survey.Input{Message: t.PromptOpusModel + " (clear to empty):", Default: p.OpusModel}, &p.OpusModel)
	survey.AskOne(&survey.Input{Message: t.PromptHaikuModel + " (clear to empty):", Default: p.HaikuModel}, &p.HaikuModel)

	var timeoutStr string
	survey.AskOne(&survey.Input{Message: t.PromptTimeout + ":", Default: strconv.Itoa(p.Timeout)}, &timeoutStr)
	if timeout, err := strconv.Atoi(timeoutStr); err == nil {
		p.Timeout = timeout
	}
}

func editField(p *config.Provider, fieldIndex int, t i18n.Messages) {
	switch fieldIndex {
	case 1: // Name
		survey.AskOne(&survey.Input{Message: t.PromptName + ":", Default: p.Name}, &p.Name)
	case 2: // Alias
		survey.AskOne(&survey.Input{Message: t.PromptAlias + ":", Default: p.Alias}, &p.Alias)
	case 3: // BaseURL
		survey.AskOne(&survey.Input{Message: t.PromptBaseURL + ":", Default: p.BaseURL}, &p.BaseURL)
	case 4: // APIKey
		var apiKey string
		survey.AskOne(&survey.Password{Message: t.PromptAPIKey + ":"}, &apiKey)
		if apiKey != "" {
			p.APIKey = apiKey
		}
	case 5: // Model
		survey.AskOne(&survey.Input{Message: t.PromptModel + " (clear to empty):", Default: p.Model}, &p.Model)
	case 6: // SmallModel
		survey.AskOne(&survey.Input{Message: t.PromptSmallModel + " (clear to empty):", Default: p.SmallModel}, &p.SmallModel)
	case 7: // SonnetModel
		survey.AskOne(&survey.Input{Message: t.PromptSonnetModel + " (clear to empty):", Default: p.SonnetModel}, &p.SonnetModel)
	case 8: // OpusModel
		survey.AskOne(&survey.Input{Message: t.PromptOpusModel + " (clear to empty):", Default: p.OpusModel}, &p.OpusModel)
	case 9: // HaikuModel
		survey.AskOne(&survey.Input{Message: t.PromptHaikuModel + " (clear to empty):", Default: p.HaikuModel}, &p.HaikuModel)
	case 10: // Timeout
		var timeoutStr string
		survey.AskOne(&survey.Input{Message: t.PromptTimeout + ":", Default: strconv.Itoa(p.Timeout)}, &timeoutStr)
		if timeout, err := strconv.Atoi(timeoutStr); err == nil {
			p.Timeout = timeout
		}
	}
}

// updateClaudeSettings updates the Claude Code settings with the provider config
func updateClaudeSettings(p *config.Provider, t i18n.Messages) error {
	settings, err := claude.LoadSettings()
	if err != nil {
		return err
	}

	// Clear existing provider settings first to ensure clean state
	settings.ClearProviderSettings()

	// Apply new provider settings
	settings.ApplyProvider(p)

	// Save settings
	return settings.Save()
}
