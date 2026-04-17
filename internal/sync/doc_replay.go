// Package sync provides synchronisation primitives for vaultpull.
//
// # Replay
//
// ReplayStore records the last successful secret payload for every profile so
// that a subsequent sync can replay (re-apply) it without contacting Vault.
// This is useful when Vault is temporarily unreachable but the operator still
// wants to ensure the local .env file is up-to-date.
//
// Usage:
//
//	store := sync.NewReplayStore("/tmp/vaultpull-replay")
//	store.Save("prod", secrets)
//	payload, ok := store.Load("prod")
package sync
