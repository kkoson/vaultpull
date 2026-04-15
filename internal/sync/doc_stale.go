// Package sync provides synchronisation primitives for vaultpull.
//
// # Staleness Tracking
//
// StalenessTracker records the last successful sync timestamp for each
// profile and exposes a simple IsStale query.
//
// Usage:
//
//	tracker := sync.NewStalenessTracker(sync.DefaultStalenessConfig())
//
//	// After a successful sync:
//	tracker.Record(profile.Name)
//
//	// Before scheduling the next sync:
//	if tracker.IsStale(profile.Name) {
//		// re-sync immediately
//	}
//
// The default MaxAge is 30 minutes. A zero MaxAge in the supplied config
// falls back to the default automatically.
package sync
