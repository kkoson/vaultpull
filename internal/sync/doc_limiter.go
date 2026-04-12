// Package sync provides synchronisation primitives for vaultpull.
//
// # Rate Limiter
//
// NewLimiter creates a token-bucket rate limiter that controls how many
// Vault API requests are issued per second across concurrent profile syncs.
//
// A rate of 0 disables limiting entirely — every call to Wait returns
// immediately without consuming a token.
//
// Usage:
//
//	limiter := sync.NewLimiter(cfg.RequestsPerSecond)
//
//	// Before each Vault API call:
//	limiter.Wait(ctx)
//
// The limiter respects context cancellation; if the context is cancelled
// while waiting for a token, Wait returns immediately so the caller can
// propagate the cancellation.
package sync
