// Package sync provides synchronization primitives and orchestration
// logic for syncing secrets from HashiCorp Vault into local .env files.
//
// # Concurrency Limiter
//
// The Limiter type provides a token-bucket style concurrency limiter
// that caps the number of profiles synced in parallel. This prevents
// overwhelming the Vault server with simultaneous requests when many
// profiles are configured.
//
// Example usage:
//
//	limiter := sync.NewLimiter(5) // allow 5 concurrent syncs
//	err := limiter.Run(ctx, func() error {
//	    return syncer.Run(ctx)
//	})
//
// The Limiter is safe for concurrent use by multiple goroutines.
package sync
