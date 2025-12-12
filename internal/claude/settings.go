package claude

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/katz/ccs/internal/config"
)

// Settings represents the Claude Code settings.json structure
type Settings struct {
	Model                   string                 `json:"model,omitempty"`
	AlwaysThinkingEnabled   *bool                  `json:"alwaysThinkingEnabled,omitempty"`
	Env                     map[string]interface{} `json:"env,omitempty"`
	Permissions             map[string]interface{} `json:"permissions,omitempty"`
	AllowedDirectories      []string               `json:"allowedDirectories,omitempty"`
	McpServers              map[string]interface{} `json:"mcpServers,omitempty"`
	AdditionalFields        map[string]interface{} `json:"-"`
}

// stripJSONComments removes comments from JSON content (JSONC format)
// and handles trailing commas
func stripJSONComments(data []byte) []byte {
	content := string(data)
	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		originalLine := line

		// Remove # style comments (commonly used in Claude settings)
		if idx := strings.Index(line, "#"); idx >= 0 {
			// Check if # is inside a string by counting quotes before it
			beforeHash := line[:idx]
			quoteCount := strings.Count(beforeHash, `"`) - strings.Count(beforeHash, `\"`)
			if quoteCount%2 == 0 {
				// # is not inside a string, it's a comment
				line = beforeHash
			}
		}

		// Remove // style comments
		if idx := strings.Index(line, "//"); idx >= 0 {
			beforeSlash := line[:idx]
			quoteCount := strings.Count(beforeSlash, `"`) - strings.Count(beforeSlash, `\"`)
			if quoteCount%2 == 0 {
				line = beforeSlash
			}
		}

		// Skip lines that are now empty or whitespace-only (were comment-only lines)
		trimmed := strings.TrimSpace(line)
		if trimmed == "" && strings.TrimSpace(originalLine) != "" {
			continue
		}

		result = append(result, line)
	}

	joined := strings.Join(result, "\n")

	// Remove trailing commas before } or ]
	trailingCommaRe := regexp.MustCompile(`,(\s*[\}\]])`)
	joined = trailingCommaRe.ReplaceAllString(joined, "$1")

	return []byte(joined)
}

// GetSettingsPath returns the Claude Code settings.json path
func GetSettingsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claude", "settings.json"), nil
}

// LoadSettings loads the Claude Code settings.json
func LoadSettings() (*Settings, error) {
	path, err := GetSettingsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty settings if file doesn't exist
			return &Settings{
				Env: make(map[string]interface{}),
			}, nil
		}
		return nil, err
	}

	// Strip comments from JSONC format
	cleanData := stripJSONComments(data)

	// First, unmarshal to a generic map to preserve all fields
	var rawMap map[string]interface{}
	if err := json.Unmarshal(cleanData, &rawMap); err != nil {
		return nil, err
	}

	// Then unmarshal to our struct
	var settings Settings
	if err := json.Unmarshal(cleanData, &settings); err != nil {
		return nil, err
	}

	// Store additional fields that we don't explicitly handle
	settings.AdditionalFields = make(map[string]interface{})
	knownFields := map[string]bool{
		"model": true, "alwaysThinkingEnabled": true, "env": true,
		"permissions": true, "allowedDirectories": true, "mcpServers": true,
	}
	for key, value := range rawMap {
		if !knownFields[key] {
			settings.AdditionalFields[key] = value
		}
	}

	if settings.Env == nil {
		settings.Env = make(map[string]interface{})
	}

	return &settings, nil
}

// SaveSettings saves the Claude Code settings.json
func (s *Settings) Save() error {
	path, err := GetSettingsPath()
	if err != nil {
		return err
	}

	// Create .claude directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Build the output map with proper field ordering
	output := make(map[string]interface{})

	// Add known fields in order
	if s.Model != "" {
		output["model"] = s.Model
	}
	if s.AlwaysThinkingEnabled != nil {
		output["alwaysThinkingEnabled"] = *s.AlwaysThinkingEnabled
	}
	if len(s.Env) > 0 {
		output["env"] = s.Env
	}
	if len(s.Permissions) > 0 {
		output["permissions"] = s.Permissions
	}
	if len(s.AllowedDirectories) > 0 {
		output["allowedDirectories"] = s.AllowedDirectories
	}
	if len(s.McpServers) > 0 {
		output["mcpServers"] = s.McpServers
	}

	// Add any additional fields that were in the original file
	for key, value := range s.AdditionalFields {
		output[key] = value
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ApplyProvider applies a provider configuration to the settings
func (s *Settings) ApplyProvider(p *config.Provider) {
	if s.Env == nil {
		s.Env = make(map[string]interface{})
	}

	// Set the provider configuration (always set these)
	s.Env["ANTHROPIC_BASE_URL"] = p.BaseURL
	s.Env["ANTHROPIC_AUTH_TOKEN"] = p.APIKey

	// Handle timeout: delete first, then set only if > 0
	delete(s.Env, "API_TIMEOUT_MS")
	if p.Timeout > 0 {
		s.Env["API_TIMEOUT_MS"] = strconv.Itoa(p.Timeout)
	}

	// Handle models: delete first, then set only if not empty
	// This ensures that empty model fields don't appear in settings.json
	// and that previously set models are removed when changed to empty
	modelFields := map[string]string{
		"ANTHROPIC_MODEL":                p.Model,
		"ANTHROPIC_SMALL_FAST_MODEL":     p.SmallModel,
		"ANTHROPIC_DEFAULT_SONNET_MODEL": p.SonnetModel,
		"ANTHROPIC_DEFAULT_OPUS_MODEL":   p.OpusModel,
		"ANTHROPIC_DEFAULT_HAIKU_MODEL":  p.HaikuModel,
	}

	for key, value := range modelFields {
		delete(s.Env, key)
		if value != "" {
			s.Env[key] = value
		}
	}

	// Disable non-essential traffic for third-party providers
	s.Env["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"] = 1
}

// ClearProviderSettings removes provider-related settings from env
func (s *Settings) ClearProviderSettings() {
	delete(s.Env, "ANTHROPIC_BASE_URL")
	delete(s.Env, "ANTHROPIC_AUTH_TOKEN")
	delete(s.Env, "API_TIMEOUT_MS")
	delete(s.Env, "ANTHROPIC_MODEL")
	delete(s.Env, "ANTHROPIC_SMALL_FAST_MODEL")
	delete(s.Env, "ANTHROPIC_DEFAULT_SONNET_MODEL")
	delete(s.Env, "ANTHROPIC_DEFAULT_OPUS_MODEL")
	delete(s.Env, "ANTHROPIC_DEFAULT_HAIKU_MODEL")
	delete(s.Env, "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC")
}

// GetCurrentEnvConfig returns a summary of the current env configuration
func (s *Settings) GetCurrentEnvConfig() map[string]string {
	result := make(map[string]string)

	keys := []string{
		"ANTHROPIC_BASE_URL",
		"ANTHROPIC_AUTH_TOKEN",
		"API_TIMEOUT_MS",
		"ANTHROPIC_MODEL",
		"ANTHROPIC_SMALL_FAST_MODEL",
		"ANTHROPIC_DEFAULT_SONNET_MODEL",
		"ANTHROPIC_DEFAULT_OPUS_MODEL",
		"ANTHROPIC_DEFAULT_HAIKU_MODEL",
	}

	for _, key := range keys {
		if val, ok := s.Env[key]; ok {
			switch v := val.(type) {
			case string:
				result[key] = v
			case float64:
				result[key] = strconv.FormatFloat(v, 'f', -1, 64)
			case int:
				result[key] = strconv.Itoa(v)
			default:
				result[key] = ""
			}
		}
	}

	return result
}
