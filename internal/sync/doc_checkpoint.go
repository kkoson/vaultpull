// Package sync provides synchronisation primitives for vaultpull.
//
// # Checkpoint
//
// Checkpoint tracks the last successful sync timestamp for every profile
// and persists that state to a JSON file on disk.  On subsequent runs the
// syncer can compare the recorded timestamp against the Vault secret's
// last-modified time and skip profiles whose secrets have not changed,
// reducing unnecessary writes and Vault API calls.
//
// Usage:
//
//	cp, err := sync.NewCheckpoint(".vaultpull.checkpoint.json")
//	if err != nil { ... }
//
//	if entry, ok := cp.Get(profile.Name); ok {
//	    // compare entry.SyncedAt with vault metadata
//	}
//
//	_ = cp.Set(sync.CheckpointEntry{
//	    Profile:    profile.Name,
//	    SyncedAt:   time.Now().UTC(),
//	    SecretPath: profile.VaultPath,
//	})
package sync
