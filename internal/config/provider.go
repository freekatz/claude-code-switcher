package config

// Provider represents a Claude Code API provider configuration
type Provider struct {
	Name        string `json:"name"`         // Provider display name
	Alias       string `json:"alias"`        // Provider short alias
	BaseURL     string `json:"base_url"`     // API Base URL
	APIKey      string `json:"api_key"`      // API Key / Auth Token
	Model       string `json:"model"`        // Main model (ANTHROPIC_MODEL)
	SmallModel  string `json:"small_model"`  // Small/fast model (ANTHROPIC_SMALL_FAST_MODEL)
	SonnetModel string `json:"sonnet_model"` // Sonnet model (ANTHROPIC_DEFAULT_SONNET_MODEL)
	OpusModel   string `json:"opus_model"`   // Opus model (ANTHROPIC_DEFAULT_OPUS_MODEL)
	HaikuModel  string `json:"haiku_model"`  // Haiku model (ANTHROPIC_DEFAULT_HAIKU_MODEL)
	Timeout     int    `json:"timeout_ms"`   // API timeout in milliseconds
}

// FillDefaults fills empty model fields with the main model value
func (p *Provider) FillDefaults() {
	if p.Model == "" {
		return
	}
	if p.SmallModel == "" {
		p.SmallModel = p.Model
	}
	if p.SonnetModel == "" {
		p.SonnetModel = p.Model
	}
	if p.OpusModel == "" {
		p.OpusModel = p.Model
	}
	if p.HaikuModel == "" {
		p.HaikuModel = p.Model
	}
	if p.Timeout == 0 {
		p.Timeout = 300000 // Default 5 minutes
	}
}
