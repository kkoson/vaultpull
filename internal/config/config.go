package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Profile holds the configuration for a single named profile.
type Profile struct {
	VaultAddr  string            `mapstructure:"vault_addr"`
	VaultToken string            `mapstructure:"vault_token"`
	VaultRole  string            `mapstructure:"vault_role"`
	AuthMethod string            `mapstructure:"auth_method"`
	Secrets    []SecretMapping   `mapstructure:"secrets"`
	EnvFile    string            `mapstructure:"env_file"`
	ExtraVars  map[string]string `mapstructure:"extra_vars"`
}

// SecretMapping maps a Vault path to an optional key override.
type SecretMapping struct {
	Path   string            `mapstructure:"path"`
	Keys   map[string]string `mapstructure:"keys"` // vault_key -> env_key
}

// Config is the top-level configuration structure.
type Config struct {
	Profiles map[string]Profile `mapstructure:"profiles"`
}

// Load reads the configuration using viper and returns a Config.
func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}
	return &cfg, nil
}

// GetProfile returns the named profile or an error if it does not exist.
func (c *Config) GetProfile(name string) (*Profile, error) {
	p, ok := c.Profiles[name]
	if !ok {
		return nil, fmt.Errorf("profile %q not found in config", name)
	}
	return &p, nil
}
