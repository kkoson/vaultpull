// Package sync provides lease management for profile sync operations.
//
// # Lease Manager
//
// LeaseManager tracks time-bounded leases per profile to avoid redundant
// sync operations when secrets have been recently fetched.
//
// A lease is acquired before each sync and released when the sync completes.
// If a valid lease exists and has not crossed the renewal threshold, the sync
// is skipped entirely.
//
// Usage:
//
//	cfg := sync.DefaultLeaseConfig()
//	manager := sync.NewLeaseManager(cfg)
//
//	protected := sync.WithLease(manager, "production", func(ctx context.Context) error {
//		return syncer.Run(ctx)
//	})
//
//	if err := protected(ctx); err != nil {
//		log.Fatal(err)
//	}
package sync
