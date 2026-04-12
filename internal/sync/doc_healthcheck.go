// Package sync provides synchronization primitives and utilities for
// vaultpull's secret sync pipeline.
//
// # HealthChecker
//
// HealthChecker runs named probe functions concurrently to determine
// the operational status of external dependencies (e.g. Vault connectivity,
// filesystem availability).
//
// Usage:
//
//	hc := sync.NewHealthChecker()
//	hc.Register("vault", func(ctx context.Context) error {
//	    return client.Ping(ctx)
//	})
//
//	results := hc.RunAll(ctx)
//	fmt.Println(sync.Summary(results))
//
// Each probe returns nil for healthy or a descriptive error for unhealthy.
// The Overall helper aggregates all results into a single HealthStatus.
package sync
