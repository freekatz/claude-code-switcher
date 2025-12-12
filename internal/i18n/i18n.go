package i18n

// Language represents supported languages
type Language string

const (
	English Language = "en"
	Chinese Language = "zh"
)

// Messages contains all translatable strings
type Messages struct {
	// General
	AppName        string
	AppDescription string
	Version        string

	// Commands
	CmdList       string
	CmdListDesc   string
	CmdAdd        string
	CmdAddDesc    string
	CmdEdit       string
	CmdEditDesc   string
	CmdUse        string
	CmdUseDesc    string
	CmdRemove     string
	CmdRemoveDesc string
	CmdCurrent    string
	CmdCurrentDesc string
	CmdLang       string
	CmdLangDesc   string

	// Prompts
	PromptName       string
	PromptAlias      string
	PromptBaseURL    string
	PromptAPIKey     string
	PromptModel      string
	PromptSmallModel string
	PromptSonnetModel string
	PromptOpusModel  string
	PromptHaikuModel string
	PromptTimeout    string
	PromptConfirmDelete string
	PromptSelectProvider string
	PromptSelectField string

	// Messages
	MsgNoProviders      string
	MsgProviderAdded    string
	MsgProviderUpdated  string
	MsgProviderRemoved  string
	MsgProviderSwitched string
	MsgCurrentProvider  string
	MsgNoCurrentProvider string
	MsgLangSwitched     string
	MsgProviderList     string

	// Errors
	ErrProviderNotFound string
	ErrProviderExists   string
	ErrInvalidAlias     string
	ErrLoadConfig       string
	ErrSaveConfig       string
	ErrUpdateSettings   string

	// Labels
	LabelName      string
	LabelAlias     string
	LabelBaseURL   string
	LabelAPIKey    string
	LabelModel     string
	LabelSmallModel string
	LabelSonnetModel string
	LabelOpusModel string
	LabelHaikuModel string
	LabelTimeout   string
	LabelCurrent   string
	LabelAll       string

	// Field names for edit
	FieldName       string
	FieldAlias      string
	FieldBaseURL    string
	FieldAPIKey     string
	FieldModel      string
	FieldSmallModel string
	FieldSonnetModel string
	FieldOpusModel  string
	FieldHaikuModel string
	FieldTimeout    string
}

var currentLang Language = English

var messages = map[Language]Messages{
	English: {
		AppName:        "ccs",
		AppDescription: "Claude Code Switcher - Manage multiple Claude Code API providers",
		Version:        "Version",

		CmdList:       "list",
		CmdListDesc:   "List all providers",
		CmdAdd:        "add",
		CmdAddDesc:    "Add a new provider",
		CmdEdit:       "edit",
		CmdEditDesc:   "Edit a provider",
		CmdUse:        "use",
		CmdUseDesc:    "Switch to a provider",
		CmdRemove:     "remove",
		CmdRemoveDesc: "Remove a provider",
		CmdCurrent:    "current",
		CmdCurrentDesc: "Show current configuration",
		CmdLang:       "lang",
		CmdLangDesc:   "Switch language (en/zh)",

		PromptName:          "Provider name",
		PromptAlias:         "Provider alias (short name)",
		PromptBaseURL:       "API Base URL",
		PromptAPIKey:        "API Key",
		PromptModel:         "Main model name",
		PromptSmallModel:    "Small/fast model (leave empty to use main model)",
		PromptSonnetModel:   "Sonnet model (leave empty to use main model)",
		PromptOpusModel:     "Opus model (leave empty to use main model)",
		PromptHaikuModel:    "Haiku model (leave empty to use main model)",
		PromptTimeout:       "API timeout in ms (default: 300000)",
		PromptConfirmDelete: "Are you sure you want to delete provider '%s'?",
		PromptSelectProvider: "Select a provider",
		PromptSelectField:   "Select field to edit",

		MsgNoProviders:      "No providers configured. Use 'ccs add' to add one.",
		MsgProviderAdded:    "Provider '%s' added successfully.",
		MsgProviderUpdated:  "Provider '%s' updated successfully.",
		MsgProviderRemoved:  "Provider '%s' removed successfully.",
		MsgProviderSwitched: "Switched to provider '%s'.",
		MsgCurrentProvider:  "Current provider: %s",
		MsgNoCurrentProvider: "No provider is currently active.",
		MsgLangSwitched:     "Language switched to %s.",
		MsgProviderList:     "Configured providers:",

		ErrProviderNotFound: "Provider '%s' not found.",
		ErrProviderExists:   "Provider with alias '%s' already exists.",
		ErrInvalidAlias:     "Invalid provider alias.",
		ErrLoadConfig:       "Failed to load configuration: %v",
		ErrSaveConfig:       "Failed to save configuration: %v",
		ErrUpdateSettings:   "Failed to update Claude Code settings: %v",

		LabelName:       "Name",
		LabelAlias:      "Alias",
		LabelBaseURL:    "Base URL",
		LabelAPIKey:     "API Key",
		LabelModel:      "Model",
		LabelSmallModel: "Small Model",
		LabelSonnetModel: "Sonnet Model",
		LabelOpusModel:  "Opus Model",
		LabelHaikuModel: "Haiku Model",
		LabelTimeout:    "Timeout (ms)",
		LabelCurrent:    "(current)",
		LabelAll:        "All fields",

		FieldName:       "Name",
		FieldAlias:      "Alias",
		FieldBaseURL:    "Base URL",
		FieldAPIKey:     "API Key",
		FieldModel:      "Main Model",
		FieldSmallModel: "Small Model",
		FieldSonnetModel: "Sonnet Model",
		FieldOpusModel:  "Opus Model",
		FieldHaikuModel: "Haiku Model",
		FieldTimeout:    "Timeout",
	},
	Chinese: {
		AppName:        "ccs",
		AppDescription: "Claude Code Switcher - 管理多个 Claude Code API 提供商",
		Version:        "版本",

		CmdList:       "list",
		CmdListDesc:   "列出所有提供商",
		CmdAdd:        "add",
		CmdAddDesc:    "添加新提供商",
		CmdEdit:       "edit",
		CmdEditDesc:   "编辑提供商配置",
		CmdUse:        "use",
		CmdUseDesc:    "切换到指定提供商",
		CmdRemove:     "remove",
		CmdRemoveDesc: "删除提供商",
		CmdCurrent:    "current",
		CmdCurrentDesc: "显示当前配置",
		CmdLang:       "lang",
		CmdLangDesc:   "切换语言 (en/zh)",

		PromptName:          "提供商名称",
		PromptAlias:         "提供商缩写",
		PromptBaseURL:       "API Base URL",
		PromptAPIKey:        "API Key",
		PromptModel:         "主模型名称",
		PromptSmallModel:    "小模型 (留空使用主模型)",
		PromptSonnetModel:   "Sonnet 模型 (留空使用主模型)",
		PromptOpusModel:     "Opus 模型 (留空使用主模型)",
		PromptHaikuModel:    "Haiku 模型 (留空使用主模型)",
		PromptTimeout:       "API 超时时间(毫秒，默认: 300000)",
		PromptConfirmDelete: "确定要删除提供商 '%s' 吗？",
		PromptSelectProvider: "选择提供商",
		PromptSelectField:   "选择要编辑的字段",

		MsgNoProviders:      "暂无配置的提供商。使用 'ccs add' 添加一个。",
		MsgProviderAdded:    "提供商 '%s' 添加成功。",
		MsgProviderUpdated:  "提供商 '%s' 更新成功。",
		MsgProviderRemoved:  "提供商 '%s' 删除成功。",
		MsgProviderSwitched: "已切换到提供商 '%s'。",
		MsgCurrentProvider:  "当前提供商: %s",
		MsgNoCurrentProvider: "当前没有激活的提供商。",
		MsgLangSwitched:     "语言已切换为 %s。",
		MsgProviderList:     "已配置的提供商:",

		ErrProviderNotFound: "未找到提供商 '%s'。",
		ErrProviderExists:   "缩写为 '%s' 的提供商已存在。",
		ErrInvalidAlias:     "无效的提供商缩写。",
		ErrLoadConfig:       "加载配置失败: %v",
		ErrSaveConfig:       "保存配置失败: %v",
		ErrUpdateSettings:   "更新 Claude Code 设置失败: %v",

		LabelName:       "名称",
		LabelAlias:      "缩写",
		LabelBaseURL:    "Base URL",
		LabelAPIKey:     "API Key",
		LabelModel:      "模型",
		LabelSmallModel: "小模型",
		LabelSonnetModel: "Sonnet 模型",
		LabelOpusModel:  "Opus 模型",
		LabelHaikuModel: "Haiku 模型",
		LabelTimeout:    "超时(ms)",
		LabelCurrent:    "(当前)",
		LabelAll:        "所有字段",

		FieldName:       "名称",
		FieldAlias:      "缩写",
		FieldBaseURL:    "Base URL",
		FieldAPIKey:     "API Key",
		FieldModel:      "主模型",
		FieldSmallModel: "小模型",
		FieldSonnetModel: "Sonnet 模型",
		FieldOpusModel:  "Opus 模型",
		FieldHaikuModel: "Haiku 模型",
		FieldTimeout:    "超时时间",
	},
}

// SetLanguage sets the current language
func SetLanguage(lang Language) {
	if lang == English || lang == Chinese {
		currentLang = lang
	}
}

// GetLanguage returns the current language
func GetLanguage() Language {
	return currentLang
}

// T returns the messages for the current language
func T() Messages {
	return messages[currentLang]
}
