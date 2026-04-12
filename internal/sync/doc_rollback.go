// Package sync provides synchronisation primitives for vaultpull.
//
// # Rollback
//
// RollbackStore gives each sync operation a safety net: before writing new
// secrets to a .env file, the caller should call Backup to snapshot the
// current file contents.  If the sync subsequently fails, Restore brings the
// file back to its previous state so that the running application is never
// left with a partially-written or corrupted env file.
//
// Usage:
//
//	store, err := sync.NewRollbackStore(".vaultpull/backups")
//	if err != nil { ... }
//
//	// Before sync:
//	store.Backup(profile.Name, profile.OutputFile)
//
//	// On failure:
//	store.Restore(profile.Name, profile.OutputFile)
//
//	// On success:
//	store.Clear(profile.Name)
//
// Backup files are stored as <profile>.bak inside the configured directory
// and are only readable by the current user (mode 0600).
package sync
