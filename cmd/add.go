package cmd

import (
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/katz/ccs/internal/config"
	"github.com/katz/ccs/internal/i18n"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add a new provider / 添加新提供商 (alias: a)",
	Long:    `Add a new API provider configuration interactively.`,
	Run:     runAdd,
}

func runAdd(cmd *cobra.Command, args []string) {
	t := i18n.T()

	cfg, err := config.Load()
	if err != nil {
		color.Red(t.ErrLoadConfig, err)
		return
	}

	var provider config.Provider

	// Required fields
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: t.PromptName + ":"},
			Validate: survey.Required,
		},
		{
			Name:     "alias",
			Prompt:   &survey.Input{Message: t.PromptAlias + ":"},
			Validate: survey.Required,
		},
		{
			Name:     "baseurl",
			Prompt:   &survey.Input{Message: t.PromptBaseURL + ":"},
			Validate: survey.Required,
		},
		{
			Name:     "apikey",
			Prompt:   &survey.Password{Message: t.PromptAPIKey + ":"},
			Validate: survey.Required,
		},
		{
			Name:   "model",
			Prompt: &survey.Input{Message: t.PromptModel + ":"},
		},
	}

	answers := struct {
		Name    string
		Alias   string
		BaseURL string
		APIKey  string
		Model   string
	}{}

	if err := survey.Ask(questions, &answers); err != nil {
		return
	}

	provider.Name = answers.Name
	provider.Alias = answers.Alias
	provider.BaseURL = answers.BaseURL
	provider.APIKey = answers.APIKey
	provider.Model = answers.Model

	// Optional model fields
	var smallModel, sonnetModel, opusModel, haikuModel string
	var timeoutStr string

	survey.AskOne(&survey.Input{Message: t.PromptSmallModel + ":"}, &smallModel)
	survey.AskOne(&survey.Input{Message: t.PromptSonnetModel + ":"}, &sonnetModel)
	survey.AskOne(&survey.Input{Message: t.PromptOpusModel + ":"}, &opusModel)
	survey.AskOne(&survey.Input{Message: t.PromptHaikuModel + ":"}, &haikuModel)
	survey.AskOne(&survey.Input{Message: t.PromptTimeout + ":", Default: "300000"}, &timeoutStr)

	provider.SmallModel = smallModel
	provider.SonnetModel = sonnetModel
	provider.OpusModel = opusModel
	provider.HaikuModel = haikuModel

	if timeout, err := strconv.Atoi(timeoutStr); err == nil {
		provider.Timeout = timeout
	} else {
		provider.Timeout = 300000
	}

	// Add provider (FillDefaults is called inside AddProvider)
	if err := cfg.AddProvider(provider); err != nil {
		if err == config.ErrProviderExists {
			color.Red(t.ErrProviderExists, provider.Alias)
		} else {
			color.Red(t.ErrSaveConfig, err)
		}
		return
	}

	// Save config
	if err := cfg.Save(); err != nil {
		color.Red(t.ErrSaveConfig, err)
		return
	}

	color.Green(t.MsgProviderAdded, provider.Name)
}
