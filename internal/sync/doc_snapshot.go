// Package sync provides synchronization primitives and utilities for
// pulling secrets from HashiCorp Vault into local .env files.
//
// # Snapshot Store
//
// SnapshotStore persists the last known secret values for each profile
// to disk. This allows vaultpull to detect drift between the local .env
// file and the remote Vault path without performing a full sync.
//
// Usage:
//
//	store, err := sync.NewSnapshotStore("/tmp/vaultpull-snapshots.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Save a snapshot after a successful sync.
//	store.Save("production", secrets)
//
//	// Retrieve the previous snapshot for diffing.
//	prev, ok := store.Get("production")
package sync
