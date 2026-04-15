package sync

import (
	"context"
	"errors"

	"github.com/yourusername/vaultpull/internal/config"
)

// ErrPolicyDenied is returned when the RetentionPolicy blocks a sync.
var ErrPolicyDenied = errors.New("sync denied by retention policy")

// WithRetentionPolicy returns a StageFunc middleware that enforces the given
// PolicyEnforcer before allowing a sync to proceed. On success, it records
// the sync time so future calls reflect the updated age.
//
// If the enforcer is nil, WithRetentionPolicy panics.
func WithRetentionPolicy(enforcer *PolicyEnforcer) func(StageFunc) StageFunc {
	if enforcer == nil {
		panic("WithRetentionPolicy: enforcer must not be nil")
	}
	return func(next StageFunc) StageFunc {
		return func(ctx context.Context, profile config.Profile) error {
			if !enforcer.Allow(profile.Name) {
				return ErrPolicyDenied
			}
			err := next(ctx, profile)
			if err == nil {
				enforcer.Record(profile.Name)
			}
			return err
		}
	}
}
