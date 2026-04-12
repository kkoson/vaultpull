package sync

import (
	"context"
	"fmt"
)

// RunAllConcurrent runs all profiles concurrently using a WorkerPool.
// workers controls the parallelism level (>= 1).
// It collects all errors and returns a combined error if any profile failed.
func (r *Runner) RunAllConcurrent(ctx context.Context, workers int) error {
	profiles := r.cfg.Profiles
	if len(profiles) == 0 {
		return nil
	}

	pool := NewWorkerPool(workers)
	results := pool.Start(ctx)

	go func() {
		for _, p := range profiles {
			name := p.Name
			_ = pool.Submit(ctx, name, func(ctx context.Context) error {
				return r.RunProfile(ctx, name)
			})
		}
		pool.Close()
	}()

	var errs []string
	for res := range results {
		if res.Err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", res.ProfileName, res.Err))
		}
	}

	if len(errs) > 0 {
		combined := ""
		for i, e := range errs {
			if i > 0 {
				combined += "; "
			}
			combined += e
		}
		return fmt.Errorf("concurrent sync errors: %s", combined)
	}
	return nil
}
