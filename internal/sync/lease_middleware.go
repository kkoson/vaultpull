package sync

import (
	"context"
	"fmt"
	"time"
)

// WithLease wraps a SyncFunc so that it acquires a lease before execution
// and releases it afterward. If a valid lease already exists and does not
// need renewal, the sync is skipped.
func WithLease(manager *LeaseManager, profile string, fn func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		now := time.Now()

		if !manager.NeedsRenewal(profile, now) {
			// Lease is still healthy; skip execution.
			return nil
		}

		manager.Acquire(profile, now)
		defer manager.Release(profile)

		if err := fn(ctx); err != nil {
			return fmt.Errorf("lease(%s): %w", profile, err)
		}
		return nil
	}
}
