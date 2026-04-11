// Package sync provides synchronisation primitives for vaultpull.
//
// # Backoff
//
// BackoffConfig controls how long the syncer waits between consecutive retry
// attempts when a Vault request fails transiently.
//
// Three strategies are available:
//
//   - BackoffFixed      — constant delay on every attempt
//   - BackoffLinear     — delay grows by one BaseDelay unit per attempt
//   - BackoffExponential — delay is multiplied by Multiplier on each attempt
//
// All strategies honour MaxDelay so that delays never grow unboundedly.
// Use DefaultBackoffConfig for a sensible starting point and pass the result
// to WithRetry via RetryPolicy.Backoff.
package sync
