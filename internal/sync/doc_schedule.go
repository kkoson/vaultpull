// Package sync provides the core synchronisation logic for vaultpull.
//
// Schedule wraps a Runner and executes syncs on a fixed interval using
// a context-aware ticker. An initial sync is performed immediately when
// Start is called, followed by subsequent ticks at the configured interval.
//
// Example usage:
//
//	runner := sync.NewRunner(cfg, vaultClient, auditLogger)
//	sched, err := sync.NewSchedule(runner, 5*time.Minute, sync.DefaultOptions(), "production")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if err := sched.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
//	    log.Fatal(err)
//	}
package sync
