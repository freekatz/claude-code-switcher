package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/katz/ccs/internal/config"
	"github.com/katz/ccs/internal/i18n"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all providers / 列出所有提供商 (alias: ls)",
	Long:    `List all configured API providers with their details.`,
	Run:     runList,
}

func runList(cmd *cobra.Command, args []string) {
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

	fmt.Println(t.MsgProviderList)
	fmt.Println()

	for _, p := range cfg.Providers {
		isCurrent := p.Alias == cfg.CurrentProvider

		// Provider header
		if isCurrent {
			color.New(color.FgGreen, color.Bold).Printf("* %s", p.Name)
			color.Green(" (%s) %s", p.Alias, t.LabelCurrent)
		} else {
			color.New(color.FgWhite, color.Bold).Printf("  %s", p.Name)
			color.White(" (%s)", p.Alias)
		}
		fmt.Println()

		// Details
		printField(t.LabelBaseURL, p.BaseURL, isCurrent)
		printField(t.LabelModel, p.Model, isCurrent)

		// Only show different models
		if p.SmallModel != p.Model {
			printField(t.LabelSmallModel, p.SmallModel, isCurrent)
		}
		if p.SonnetModel != p.Model {
			printField(t.LabelSonnetModel, p.SonnetModel, isCurrent)
		}
		if p.OpusModel != p.Model {
			printField(t.LabelOpusModel, p.OpusModel, isCurrent)
		}
		if p.HaikuModel != p.Model {
			printField(t.LabelHaikuModel, p.HaikuModel, isCurrent)
		}

		fmt.Println()
	}
}

func printField(label, value string, isCurrent bool) {
	indent := "    "
	if isCurrent {
		color.Green("%s%s: %s", indent, label, value)
	} else {
		fmt.Printf("%s%s: %s\n", indent, label, value)
	}
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return strings.Repeat("*", len(key))
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}
