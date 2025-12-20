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

// managedEnvKeys are the env keys that ccs manages
var managedEnvKeys = []string{
	"ANTHROPIC_BASE_URL",
	"ANTHROPIC_AUTH_TOKEN",
	"API_TIMEOUT_MS",
	"ANTHROPIC_MODEL",
	"ANTHROPIC_SMALL_FAST_MODEL",
	"ANTHROPIC_DEFAULT_SONNET_MODEL",
	"ANTHROPIC_DEFAULT_OPUS_MODEL",
	"ANTHROPIC_DEFAULT_HAIKU_MODEL",
	"CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC",
}

// Settings wraps the raw settings.json
// ccs only manages specific keys in "env", everything else is preserved as-is
type Settings struct {
	raw map[string]interface{}
}

// getEnv returns the env map, creating it if needed
func (s *Settings) getEnv() map[string]interface{} {
	if s.raw == nil {
		s.raw = make(map[string]interface{})
	}
	env, ok := s.raw["env"].(map[string]interface{})
	if !ok {
		env = make(map[string]interface{})
		s.raw["env"] = env
	}
	return env
}

// stripJSONComments removes comments from JSON content (JSONC format)
func stripJSONComments(data []byte) []byte {
	content := string(data)
	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		originalLine := line

		// Remove # style comments
		if idx := strings.Index(line, "#"); idx >= 0 {
			beforeHash := line[:idx]
			quoteCount := strings.Count(beforeHash, `"`) - strings.Count(beforeHash, `\"`)
			if quoteCount%2 == 0 {
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

		trimmed := strings.TrimSpace(line)
		if trimmed == "" && strings.TrimSpace(originalLine) != "" {
			continue
		}

		result = append(result, line)
	}

	joined := strings.Join(result, "\n")
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
	baseDir := filepath.Join(home, ".claude")
	return filepath.Join(baseDir, "settings.json"), nil
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
			return &Settings{raw: make(map[string]interface{})}, nil
		}
		return nil, err
	}

	cleanData := stripJSONComments(data)

	var raw map[string]interface{}
	if err := json.Unmarshal(cleanData, &raw); err != nil {
		return nil, err
	}

	return &Settings{raw: raw}, nil
}

// backup creates a backup of the settings file
func backup(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	backupPath := path + ".bak"
	return os.WriteFile(backupPath, data, 0644)
}

// Save saves the settings.json with backup
func (s *Settings) Save() error {
	path, err := GetSettingsPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err := backup(path); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s.raw, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ApplyProvider applies a provider configuration to the settings
func (s *Settings) ApplyProvider(p *config.Provider) {
	env := s.getEnv()

	env["ANTHROPIC_BASE_URL"] = p.BaseURL
	env["ANTHROPIC_AUTH_TOKEN"] = p.APIKey

	delete(env, "API_TIMEOUT_MS")
	if p.Timeout > 0 {
		env["API_TIMEOUT_MS"] = strconv.Itoa(p.Timeout)
	}

	modelFields := map[string]string{
		"ANTHROPIC_MODEL":                p.Model,
		"ANTHROPIC_SMALL_FAST_MODEL":     p.SmallModel,
		"ANTHROPIC_DEFAULT_SONNET_MODEL": p.SonnetModel,
		"ANTHROPIC_DEFAULT_OPUS_MODEL":   p.OpusModel,
		"ANTHROPIC_DEFAULT_HAIKU_MODEL":  p.HaikuModel,
	}

	for key, value := range modelFields {
		delete(env, key)
		if value != "" {
			env[key] = value
		}
	}

	env["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"] = 1
}

// ClearProviderSettings removes provider-related settings from env
func (s *Settings) ClearProviderSettings() {
	env := s.getEnv()
	for _, key := range managedEnvKeys {
		delete(env, key)
	}
}

// GetCurrentEnvConfig returns a summary of the current env configuration
func (s *Settings) GetCurrentEnvConfig() map[string]string {
	result := make(map[string]string)
	env := s.getEnv()

	for _, key := range managedEnvKeys[:8] { // exclude DISABLE_NONESSENTIAL_TRAFFIC
		if val, ok := env[key]; ok {
			switch v := val.(type) {
			case string:
				result[key] = v
			case float64:
				result[key] = strconv.FormatFloat(v, 'f', -1, 64)
			case int:
				result[key] = strconv.Itoa(v)
			}
		}
	}

	return result
}
