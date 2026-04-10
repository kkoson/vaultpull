package sync

import (
	"context"
	"fmt"
	"log"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

// Syncer orchestrates fetching secrets from Vault and writing them
// to the configured .env file for a given profile.
type Syncer struct {
	cfg    *config.Config
	client *vault.Client
}

// New creates a Syncer from the provided configuration.
func New(cfg *config.Config, client *vault.Client) *Syncer {
	return &Syncer{cfg: cfg, client: client}
}

// Run executes the sync for the named profile.
// Pass an empty string to use the default profile.
func (s *Syncer) Run(ctx context.Context, profileName string) error {
	profile, err := s.cfg.GetProfile(profileName)
	if err != nil {
		return fmt.Errorf("sync: resolve profile: %w", err)
	}

	log.Printf("syncing profile %q → %s", profile.Name, profile.EnvFile)

	merged := make(map[string]string)
	for _, path := range profile.Paths {
		log.Printf("  fetching %s", path)
		secrets, err := s.client.GetSecrets(ctx, path)
		if err != nil {
			return fmt.Errorf("sync: fetch %q: %w", path, err)
		}
		for k, v := range secrets {
			merged[k] = v
		}
	}

	w := env.NewWriter(profile.EnvFile)
	if err := w.Write(merged); err != nil {
		return fmt.Errorf("sync: write env file: %w", err)
	}

	log.Printf("wrote %d secret(s) to %s", len(merged), profile.EnvFile)
	return nil
}
