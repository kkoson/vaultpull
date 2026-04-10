package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupViper(t *testing.T, raw map[string]interface{}) {
	t.Helper()
	viper.Reset()
	for k, v := range raw {
		viper.Set(k, v)
	}
}

func TestLoad_EmptyConfig(t *testing.T) {
	setupViper(t, map[string]interface{}{})
	cfg, err := Load()
	require.NoError(t, err)
	assert.NotNil(t, cfg.Profiles)
}

func TestLoad_WithProfile(t *testing.T) {
	setupViper(t, map[string]interface{}{
		"profiles": map[string]interface{}{
			"staging": map[string]interface{}{
				"vault_addr":  "https://vault.example.com",
				"vault_token": "s.abc123",
				"env_file":    ".env.staging",
			},
		},
	})

	cfg, err := Load()
	require.NoError(t, err)

	p, err := cfg.GetProfile("staging")
	require.NoError(t, err)
	assert.Equal(t, "https://vault.example.com", p.VaultAddr)
	assert.Equal(t, "s.abc123", p.VaultToken)
	assert.Equal(t, ".env.staging", p.EnvFile)
}

func TestGetProfile_NotFound(t *testing.T) {
	setupViper(t, map[string]interface{}{})
	cfg, err := Load()
	require.NoError(t, err)

	_, err = cfg.GetProfile("nonexistent")
	assert.ErrorContains(t, err, "nonexistent")
}

func TestGetProfile_Default(t *testing.T) {
	setupViper(t, map[string]interface{}{
		"profiles": map[string]interface{}{
			"default": map[string]interface{}{
				"vault_addr": "http://localhost:8200",
				"env_file":   ".env",
			},
		},
	})

	cfg, err := Load()
	require.NoError(t, err)

	p, err := cfg.GetProfile("default")
	require.NoError(t, err)
	assert.Equal(t, ".env", p.EnvFile)
}
