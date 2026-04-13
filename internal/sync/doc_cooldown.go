// Package sync provides synchronisation primitives and orchestration
// utilities for vaultpull.
//
// # Cooldown
//
// CooldownManager enforces a minimum interval between successive syncs for
// the same profile. This prevents hammering Vault when a profile is
// triggered repeatedly in a short period (e.g. by a file-watcher or a
// rapid schedule tick).
//
// Usage:
//
//	cm := sync.NewCooldownManager(sync.CooldownConfig{
//		Duration: 15 * time.Second,
//	})
//
//	if !cm.Allow(profile.Name) {
//		log.Printf("profile %s is in cooldown (%v remaining)",
//			profile.Name, cm.Remaining(profile.Name))
//		return nil
//	}
//
// Call Reset to clear the cooldown for a profile early, for example after
// a forced manual sync.
package sync
