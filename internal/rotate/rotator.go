package rotate

import (
	"context"
	"fmt"
	"time"

	"github.com/your-org/vaultpull/internal/audit"
	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/token"
	"github.com/your-org/vaultpull/internal/vault"
)

// Rotator handles AppRole secret-id rotation for a profile.
type Rotator struct {
	cfg    *config.Profile
	cache  *token.Cache
	logger *audit.Logger
}

// New creates a new Rotator for the given profile.
func New(cfg *config.Profile, cache *token.Cache, logger *audit.Logger) *Rotator {
	return &Rotator{
		cfg:    cfg,
		cache:  cache,
		logger: logger,
	}
}

// Rotate generates a new AppRole secret-id, persists it to the token cache,
// and returns the new token TTL deadline.
func (r *Rotator) Rotate(ctx context.Context, client *vault.Client) (time.Time, error) {
	if r.cfg.Auth.Method != "approle" {
		return time.Time{}, fmt.Errorf("rotate: profile %q does not use approle auth", r.cfg.Name)
	}

	roleID := r.cfg.Auth.RoleID
	if roleID == "" {
		return time.Time{}, fmt.Errorf("rotate: role_id is required for approle rotation")
	}

	newSecretID, ttl, err := generateSecretID(ctx, client, roleID)
	if err != nil {
		r.logger.Write(audit.Entry{
			Profile: r.cfg.Name,
			Action:  "rotate",
			Error:   err,
		})
		return time.Time{}, fmt.Errorf("rotate: failed to generate secret-id: %w", err)
	}

	expiry := time.Now().Add(ttl)
	if err := r.cache.Save(r.cfg.Name, newSecretID, expiry); err != nil {
		return time.Time{}, fmt.Errorf("rotate: failed to cache new secret-id: %w", err)
	}

	r.logger.Write(audit.Entry{
		Profile: r.cfg.Name,
		Action:  "rotate",
		Detail:  fmt.Sprintf("secret-id rotated, expires %s", expiry.Format(time.RFC3339)),
	})

	return expiry, nil
}
