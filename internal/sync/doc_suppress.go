// Package sync provides synchronization primitives for vaultpull.
//
// # Suppression
//
// The Suppressor type allows certain profiles to be temporarily suppressed
// from syncing based on a configurable list of profile names. Suppressed
// profiles are skipped silently or with an optional log message.
//
// Usage:
//
//	s := NewSuppressor([]string{"staging", "legacy"})
//	if s.IsSuppressed(profile) {
//	    // skip this profile
//	}
//
// The WithSuppress middleware integrates suppression into the sync pipeline.
package sync
