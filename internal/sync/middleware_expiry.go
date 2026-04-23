package sync

import (
	"context"
	"fmt"

	"github.com/yourusername/vaultpull/internal/config"
)

// WithExpiry wraps a SyncFunc so that it skips profiles whose last-sync
// entry in the ExpiryTracker has not yet expired (i.e. they are still fresh).
// On a successful sync the profile is recorded in the tracker.
//
// This is useful to prevent redundant syncs when a profile was recently
// refreshed by another path (e.g. cache warmer).
func WithExpiry(tracker *ExpiryTracker, next SyncFunc) SyncFunc {
	if tracker == nil {
		panic("WithExpiry: tracker must not be nil")
	}
	return func(ctx context.Context, p config.Profile) error {
		if !tracker.IsExpired(p.Name) {
			// Entry is still fresh — skip.
			return fmt.Errorf("expiry: profile %q is still fresh, skipping", p.Name)
		}
		if err := next(ctx, p); err != nil {
			return err
		}
		tracker.Record(p.Name)
		return nil
	}
}
