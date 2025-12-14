package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// Config represents the CCS configuration
type Config struct {
	CurrentProvider string     `json:"current_provider"` // Current active provider alias
	Providers       []Provider `json:"providers"`        // List of configured providers
}

var (
	ErrProviderNotFound = errors.New("provider not found")
	ErrProviderExists   = errors.New("provider with this alias already exists")
	ErrNoProviders      = errors.New("no providers configured")
	ErrInvalidAlias     = errors.New("invalid provider alias")
)

// GetConfigDir returns the CCS configuration directory path
func GetConfigDir() (string, error) {
	if runtime.GOOS == "windows" {
		// Windows: %APPDATA%\ccs
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "ccs"), nil
	}

	// macOS/Linux: ~/.config/ccs
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "ccs"), nil
}

// GetConfigPath returns the CCS configuration file path
func GetConfigPath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// Load loads the configuration from disk
func Load() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				Providers: []Provider{},
			}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save saves the configuration to disk
func (c *Config) Save() error {
	dir, err := GetConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// GetProvider returns a provider by alias
func (c *Config) GetProvider(alias string) (*Provider, error) {
	for i := range c.Providers {
		if c.Providers[i].Alias == alias {
			return &c.Providers[i], nil
		}
	}
	return nil, ErrProviderNotFound
}

// GetCurrentProvider returns the current active provider
func (c *Config) GetCurrentProvider() (*Provider, error) {
	if c.CurrentProvider == "" {
		return nil, ErrNoProviders
	}
	return c.GetProvider(c.CurrentProvider)
}

// AddProvider adds a new provider
func (c *Config) AddProvider(p Provider) error {
	for _, existing := range c.Providers {
		if existing.Alias == p.Alias {
			return ErrProviderExists
		}
	}

	p.FillDefaults()
	c.Providers = append(c.Providers, p)
	return nil
}

// UpdateProvider updates an existing provider
func (c *Config) UpdateProvider(alias string, p Provider) error {
	for i := range c.Providers {
		if c.Providers[i].Alias == alias {
			if alias != p.Alias {
				for j := range c.Providers {
					if i != j && c.Providers[j].Alias == p.Alias {
						return ErrProviderExists
					}
				}
				if c.CurrentProvider == alias {
					c.CurrentProvider = p.Alias
				}
			}
			p.FillDefaults()
			c.Providers[i] = p
			return nil
		}
	}
	return ErrProviderNotFound
}

// RemoveProvider removes a provider by alias
func (c *Config) RemoveProvider(alias string) error {
	for i := range c.Providers {
		if c.Providers[i].Alias == alias {
			c.Providers = append(c.Providers[:i], c.Providers[i+1:]...)
			if c.CurrentProvider == alias {
				c.CurrentProvider = ""
			}
			return nil
		}
	}
	return ErrProviderNotFound
}

// SetCurrentProvider sets the current active provider
func (c *Config) SetCurrentProvider(alias string) error {
	_, err := c.GetProvider(alias)
	if err != nil {
		return err
	}
	c.CurrentProvider = alias
	return nil
}
