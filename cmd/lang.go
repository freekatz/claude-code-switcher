package cmd

import (
	"github.com/fatih/color"
	"github.com/katz/ccs/internal/config"
	"github.com/katz/ccs/internal/i18n"
	"github.com/spf13/cobra"
)

var langCmd = &cobra.Command{
	Use:   "lang [en|zh]",
	Short: "Switch language / 切换语言",
	Long:  `Switch the display language between English (en) and Chinese (zh).`,
	Args:  cobra.MaximumNArgs(1),
	Run:   runLang,
}

func runLang(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		color.Red("Failed to load config: %v", err)
		return
	}

	if len(args) == 0 {
		// Toggle language
		if cfg.Language == "zh" {
			cfg.Language = "en"
		} else {
			cfg.Language = "zh"
		}
	} else {
		lang := args[0]
		if lang != "en" && lang != "zh" {
			color.Red("Invalid language. Use 'en' or 'zh'.")
			return
		}
		cfg.Language = lang
	}

	if err := cfg.Save(); err != nil {
		color.Red("Failed to save config: %v", err)
		return
	}

	// Update current language
	i18n.SetLanguage(i18n.Language(cfg.Language))
	t := i18n.T()

	langName := "English"
	if cfg.Language == "zh" {
		langName = "中文"
	}

	color.Green(t.MsgLangSwitched, langName)
}
