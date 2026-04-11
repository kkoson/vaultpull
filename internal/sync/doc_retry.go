// Package sync provides orchestration for syncing secrets from HashiCorp Vault
// into local .env files.
//
// # Retry
//
// WithRetry executes a fallible operation according to a RetryPolicy, supporting
// configurable attempt counts, initial delay, and exponential backoff via a
// multiplier. Context cancellation is respected both before the first attempt
// and during inter-attempt delays, allowing callers to abort long retry loops
// cleanly.
//
// Example:
//
//	policy := sync.DefaultRetryPolicy()
//	err := sync.WithRetry(ctx, policy, func() error {
//		return syncer.Run(ctx)
//	})
package sync
