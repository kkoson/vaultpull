package sync

import (
	"context"
	"time"

	"github.com/yourusername/vaultpull/internal/config"
)

// CheckpointMiddleware wraps a SyncFunc so that:
//  1. Before syncing, it checks whether the profile was already synced
//     within the given staleness window and skips it if so.
//  2. After a successful sync, it records a new checkpoint entry.
func CheckpointMiddleware(cp *Checkpoint, maxAge time.Duration, next func(context.Context, config.Profile) error) func(context.Context, config.Profile) error {
	return func(ctx context.Context, p config.Profile) error {
		if cp != nil && maxAge > 0 {
			if entry, ok := cp.Get(p.Name); ok {
				if time.Since(entry.SyncedAt) < maxAge {
					// Profile is fresh — skip.
					return nil
				}
			}
		}

		if err := next(ctx, p); err != nil {
			return err
		}

		if cp != nil {
			_ = cp.Set(CheckpointEntry{
				Profile:    p.Name,
				SyncedAt:   time.Now().UTC(),
				SecretPath: p.VaultPath,
			})
		}
		return nil
	}
}
