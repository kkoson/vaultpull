package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client.
type Client struct {
	v *api.Client
}

// Config holds the parameters needed to create a Vault client.
type Config struct {
	Address string
	Token   string
	RoleID  string
	SecretID string
}

// New creates and authenticates a new Vault client.
func New(cfg Config) (*Client, error) {
	vcfg := api.DefaultConfig()
	vcfg.Address = cfg.Address

	v, err := api.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("vault: create client: %w", err)
	}

	if cfg.Token != "" {
		v.SetToken(cfg.Token)
		return &Client{v: v}, nil
	}

	if cfg.RoleID != "" && cfg.SecretID != "" {
		token, err := approleLogin(v, cfg.RoleID, cfg.SecretID)
		if err != nil {
			return nil, err
		}
		v.SetToken(token)
		return &Client{v: v}, nil
	}

	return nil, fmt.Errorf("vault: no authentication method provided")
}

func approleLogin(v *api.Client, roleID, secretID string) (string, error) {
	data := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}
	secret, err := v.Logical().Write("auth/approle/login", data)
	if err != nil {
		return "", fmt.Errorf("vault: approle login: %w", err)
	}
	if secret == nil || secret.Auth == nil {
		return "", fmt.Errorf("vault: approle login returned no auth info")
	}
	return secret.Auth.ClientToken, nil
}
