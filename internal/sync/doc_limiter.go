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
//	if err := limiter.Wait(ctx); err != nil {
//		return err
//	}
//
// The limiter respects context cancellation; if the context is cancelled
// while waiting for a token, Wait returns the context's error so the caller
// can propagate the cancellation.
//
// # Choosing a Rate
//
// Set RequestsPerSecond to a value slightly below the Vault server's
// configured rate limit to avoid 429 responses under concurrent load.
// A value of 50 is a reasonable default for most deployments.
package sync
