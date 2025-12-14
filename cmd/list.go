package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/katz/ccs/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list [alias]",
	Aliases: []string{"ls"},
	Short:   "List providers or show provider details (alias: ls)",
	Run:     runList,
}

func runList(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		color.Red("Failed to load config: %v", err)
		return
	}

	if len(cfg.Providers) == 0 {
		color.Yellow("No providers configured. Use 'ccs add' to add one.")
		return
	}

	// If alias provided, show details for that provider
	if len(args) > 0 {
		showProviderDetail(cfg, args[0])
		return
	}

	// List all providers (alias only)
	for _, p := range cfg.Providers {
		isCurrent := p.Alias == cfg.CurrentProvider
		if isCurrent {
			color.Green("* %s", p.Alias)
		} else {
			fmt.Printf("  %s\n", p.Alias)
		}
	}
}

func showProviderDetail(cfg *config.Config, alias string) {
	p, err := cfg.GetProvider(alias)
	if err != nil {
		color.Red("Provider '%s' not found", alias)
		return
	}

	isCurrent := p.Alias == cfg.CurrentProvider

	// Header
	if isCurrent {
		color.Green("%s (%s) [current]", p.Name, p.Alias)
	} else {
		fmt.Printf("%s (%s)\n", p.Name, p.Alias)
	}

	// Details
	printDetail("URL", p.BaseURL, isCurrent)
	printDetail("Models", buildModelLine(*p), isCurrent)
	printDetail("Timeout", fmt.Sprintf("%dms", p.Timeout), isCurrent)
}

func printDetail(label, value string, isCurrent bool) {
	if isCurrent {
		color.Green("  %s: %s", label, value)
	} else {
		fmt.Printf("  %s: %s\n", label, value)
	}
}

// buildModelLine creates a compact model display string
func buildModelLine(p config.Provider) string {
	if p.Model == "" {
		return "default"
	}

	allSame := p.SmallModel == p.Model &&
		p.SonnetModel == p.Model &&
		p.OpusModel == p.Model &&
		p.HaikuModel == p.Model

	if allSame {
		return p.Model
	}

	var parts []string
	parts = append(parts, fmt.Sprintf("main:%s", p.Model))

	if p.SonnetModel != p.Model {
		parts = append(parts, fmt.Sprintf("sonnet:%s", p.SonnetModel))
	}
	if p.OpusModel != p.Model {
		parts = append(parts, fmt.Sprintf("opus:%s", p.OpusModel))
	}
	if p.HaikuModel != p.Model {
		parts = append(parts, fmt.Sprintf("haiku:%s", p.HaikuModel))
	}
	if p.SmallModel != p.Model {
		parts = append(parts, fmt.Sprintf("small:%s", p.SmallModel))
	}

	return strings.Join(parts, ", ")
}
