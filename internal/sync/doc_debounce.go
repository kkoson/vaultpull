// Package sync provides synchronisation primitives and orchestration
// utilities for vaultpull profile syncing.
//
// # Debouncer
//
// Debouncer coalesces rapid successive sync triggers into a single
// execution after a configurable quiet period. This prevents hammering
// Vault when multiple file-system events or config reloads fire in quick
// succession.
//
// Usage:
//
//	d := sync.NewDebouncer(300 * time.Millisecond)
//
//	// Called many times in rapid succession — only the last one fires.
//	for i := 0; i < 10; i++ {
//		d.Trigger(ctx, func() { runner.RunAll(ctx) })
//	}
//
// Call Flush to cancel a pending trigger without executing it.
package sync
