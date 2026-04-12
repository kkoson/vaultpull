// Package sync provides synchronisation primitives for vaultpull.
//
// # Rate Limit Store
//
// RateLimitStore enforces a per-profile cap on the number of sync operations
// that may be performed within a configurable time window.
//
// Usage:
//
//	store := sync.NewRateLimitStore(10, time.Minute)
//
//	if !store.Allow(profile.Name) {
//		// skip or return an error — limit exceeded
//		return ErrRateLimited
//	}
//	// proceed with sync …
//
// A maxRate of zero or less disables limiting entirely, which is useful
// during testing or when running in an unrestricted environment.
package sync
