// Package sync provides synchronization primitives and utilities for
// syncing secrets from HashiCorp Vault into local .env files.
//
// # Metrics
//
// The Metrics type tracks per-profile sync statistics across a run session.
// It records success, failure, and skipped counts along with timing data.
//
// Usage:
//
//	m := sync.NewMetrics()
//	m.RecordSuccess("prod", 120*time.Millisecond)
//	m.RecordFailure("staging", 50*time.Millisecond)
//
//	fmt.Println(m.Summary())
//
// Metrics is safe for concurrent use via sync.Mutex.
package sync
