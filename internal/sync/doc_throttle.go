// Package sync provides synchronisation primitives for vaultpull.
//
// # Throttle
//
// Throttle enforces a minimum interval between successive profile sync
// operations to avoid overwhelming Vault or the local filesystem.
//
// Usage:
//
//	th := sync.NewThrottle(sync.DefaultThrottleConfig())
//	if err := th.Wait(ctx); err != nil {
//	    return err
//	}
//	// perform sync ...
//
// BurstSize allows a short burst of operations before the interval is
// enforced, which is useful when syncing many profiles at startup.
package sync
