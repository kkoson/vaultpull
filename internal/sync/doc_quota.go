// Package sync provides synchronisation primitives for vaultpull.
//
// # Quota Enforcer
//
// QuotaEnforcer limits the number of sync operations performed for each
// profile within a rolling time window. This prevents runaway sync loops
// or misconfigured schedules from hammering Vault.
//
// Usage:
//
//	q := sync.NewQuotaEnforcer(sync.QuotaConfig{
//	    MaxSyncsPerWindow: 50,
//	    Window:            time.Hour,
//	})
//
//	if err := q.Allow(profile.Name); err != nil {
//	    // quota exceeded — skip or surface the error
//	    return err
//	}
//
// The enforcer is safe for concurrent use across multiple goroutines.
// Each profile maintains an independent counter and window, so a quota
// breach on one profile does not affect others.
package sync
