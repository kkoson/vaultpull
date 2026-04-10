package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds the full application configuration loaded from the
// .vaultpull.yaml file.
type Config struct {
	Vault    VaultConfig `mapstructure:"vault"`
	Profiles []Profile   `mapstructure:"profiles"`
}

// VaultConfig holds connection and authentication settings for Vault.
type VaultConfig struct {
	Address  string      `mapstructure:"address"`
	Token    string      `mapstructure:"token"`
	AppRole  AppRoleCfg  `mapstructure:"approle"`
}

// AppRoleCfg holds AppRole authentication credentials.
type AppRoleCfg struct {
	RoleID   string `mapstructure:"role_id"`
	SecretID string `mapstructure:"secret_id"`
}

// Load reads configuration from Viper into a Config struct and performs
// basic validation.
func Load(v *viper.Viper) (*Config, error) {
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if cfg.Vault.Address == "" {
		cfg.Vault.Address = "http://127.0.0.1:8200"
	}

	for i := range cfg.Profiles {
		if err := cfg.Profiles[i].Validate(); err != nil {
			return nil, fmt.Errorf("invalid profile config: %w", err)
		}
	}

	return &cfg, nil
}
