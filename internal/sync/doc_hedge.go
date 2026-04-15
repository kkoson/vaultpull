// Package sync provides synchronisation utilities for vaultpull.
//
// # Hedge
//
// Hedge implements the hedged-requests pattern to reduce tail latency.
// When a call to the supplied function does not return within the configured
// Delay, an additional (hedged) request is issued concurrently. The first
// successful result is returned and the losing goroutine's context is
// cancelled.
//
// Example:
//
//	cfg := sync.DefaultHedgeConfig()
//	cfg.Delay = 100 * time.Millisecond
//
//	val, err := sync.Hedge(ctx, cfg, func(ctx context.Context) (interface{}, error) {
//		return fetchSecret(ctx, path)
//	})
//
// MaxHedges controls how many additional requests may be in-flight at once.
// The default is 1, meaning at most two concurrent attempts.
package sync
