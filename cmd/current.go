package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/katz/ccs/internal/claude"
	"github.com/katz/ccs/internal/config"
	"github.com/katz/ccs/internal/i18n"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:     "current",
	Aliases: []string{"c"},
	Short:   "Show current configuration / 显示当前配置 (alias: c)",
	Long:    `Show the current provider configuration and Claude Code settings.`,
	Run:     runCurrent,
}

func runCurrent(cmd *cobra.Command, args []string) {
	t := i18n.T()

	cfg, err := config.Load()
	if err != nil {
		color.Red(t.ErrLoadConfig, err)
		return
	}

	// Show CCS current provider
	if cfg.CurrentProvider == "" {
		color.Yellow(t.MsgNoCurrentProvider)
	} else {
		provider, err := cfg.GetProvider(cfg.CurrentProvider)
		if err != nil {
			color.Yellow(t.MsgNoCurrentProvider)
		} else {
			color.Green(t.MsgCurrentProvider, provider.Name+" ("+provider.Alias+")")
			fmt.Println()

			// Provider details
			color.Cyan("Provider Configuration:")
			fmt.Printf("  %s: %s\n", t.LabelName, provider.Name)
			fmt.Printf("  %s: %s\n", t.LabelAlias, provider.Alias)
			fmt.Printf("  %s: %s\n", t.LabelBaseURL, provider.BaseURL)
			fmt.Printf("  %s: %s\n", t.LabelModel, provider.Model)
			if provider.SmallModel != provider.Model {
				fmt.Printf("  %s: %s\n", t.LabelSmallModel, provider.SmallModel)
			}
			if provider.SonnetModel != provider.Model {
				fmt.Printf("  %s: %s\n", t.LabelSonnetModel, provider.SonnetModel)
			}
			if provider.OpusModel != provider.Model {
				fmt.Printf("  %s: %s\n", t.LabelOpusModel, provider.OpusModel)
			}
			if provider.HaikuModel != provider.Model {
				fmt.Printf("  %s: %s\n", t.LabelHaikuModel, provider.HaikuModel)
			}
			fmt.Printf("  %s: %d\n", t.LabelTimeout, provider.Timeout)
		}
	}

	fmt.Println()

	// Show actual Claude Code settings
	settings, err := claude.LoadSettings()
	if err != nil {
		color.Red(t.ErrUpdateSettings, err)
		return
	}

	color.Cyan("Claude Code Settings (~/.claude/settings.json):")
	envConfig := settings.GetCurrentEnvConfig()

	if len(envConfig) == 0 {
		fmt.Println("  No provider environment configured")
	} else {
		for key, value := range envConfig {
			displayValue := value
			// Mask API key
			if key == "ANTHROPIC_AUTH_TOKEN" && len(value) > 8 {
				displayValue = value[:4] + "****" + value[len(value)-4:]
			}
			fmt.Printf("  %s: %s\n", key, displayValue)
		}
	}
}
