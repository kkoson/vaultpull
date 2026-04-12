// Package sync provides utilities for syncing secrets from Vault into
// local .env files.
//
// # WorkerPool
//
// WorkerPool executes profile sync jobs concurrently using a fixed number
// of goroutines. Jobs are submitted via Submit and results are streamed
// through the channel returned by Start.
//
// Basic usage:
//
//	pool := sync.NewWorkerPool(4)
//	results := pool.Start(ctx)
//
//	for _, profile := range profiles {
//		_ = pool.Submit(ctx, profile, func(ctx context.Context) error {
//			return runner.RunProfile(ctx, profile)
//		})
//	}
//	pool.Close()
//
//	for r := range results {
//		if r.Err != nil {
//			log.Printf("profile %s failed: %v", r.ProfileName, r.Err)
//		}
//	}
package sync
