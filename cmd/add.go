package cmd

import (
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/katz/ccs/internal/config"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add a new provider (alias: a)",
	Run:     runAdd,
}

func runAdd(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		color.Red("Failed to load config: %v", err)
		return
	}

	var provider config.Provider

	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name:"},
			Validate: survey.Required,
		},
		{
			Name:     "alias",
			Prompt:   &survey.Input{Message: "Alias:"},
			Validate: survey.Required,
		},
		{
			Name:     "baseurl",
			Prompt:   &survey.Input{Message: "Base URL:"},
			Validate: survey.Required,
		},
		{
			Name:     "apikey",
			Prompt:   &survey.Password{Message: "API Key:"},
			Validate: survey.Required,
		},
		{
			Name:   "model",
			Prompt: &survey.Input{Message: "Model (empty for Claude default):"},
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

	// Optional fields
	var smallModel, sonnetModel, opusModel, haikuModel, timeoutStr string

	survey.AskOne(&survey.Input{Message: "Small model (empty=main):"}, &smallModel)
	survey.AskOne(&survey.Input{Message: "Sonnet model (empty=main):"}, &sonnetModel)
	survey.AskOne(&survey.Input{Message: "Opus model (empty=main):"}, &opusModel)
	survey.AskOne(&survey.Input{Message: "Haiku model (empty=main):"}, &haikuModel)
	survey.AskOne(&survey.Input{Message: "Timeout ms:", Default: "300000"}, &timeoutStr)

	provider.SmallModel = smallModel
	provider.SonnetModel = sonnetModel
	provider.OpusModel = opusModel
	provider.HaikuModel = haikuModel

	if timeout, err := strconv.Atoi(timeoutStr); err == nil {
		provider.Timeout = timeout
	} else {
		provider.Timeout = 300000
	}

	if err := cfg.AddProvider(provider); err != nil {
		if err == config.ErrProviderExists {
			color.Red("Provider '%s' already exists", provider.Alias)
		} else {
			color.Red("Failed to add provider: %v", err)
		}
		return
	}

	if err := cfg.Save(); err != nil {
		color.Red("Failed to save config: %v", err)
		return
	}

	color.Green("Provider '%s' added", provider.Name)
}
