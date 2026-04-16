// Package sync provides synchronisation primitives for vaultpull.
//
// # Secret Pinning
//
// PinStore allows callers to "pin" a secret version for a profile so that
// subsequent syncs skip fetching a newer version from Vault until the pin
// is explicitly cleared.
//
// Usage:
//
//	store := sync.NewPinStore()
//	store.Pin("prod", "v3")
//
//	if pinned, ver := store.Get("prod"); pinned {
//	    // use cached version ver
//	}
//
//	store.Unpin("prod")
package sync
