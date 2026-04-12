package sync

import (
	"context"
	"fmt"
	"sync"
)

// ConcurrencyConfig holds configuration for concurrent profile syncing.
type ConcurrencyConfig struct {
	// MaxWorkers is the maximum number of profiles synced concurrently.
	MaxWorkers int
	// FailFast stops all workers on the first encountered error.
	FailFast bool
}

// DefaultConcurrencyConfig returns a ConcurrencyConfig with sensible defaults.
func DefaultConcurrencyConfig() ConcurrencyConfig {
	return ConcurrencyConfig{
		MaxWorkers: 4,
		FailFast:   false,
	}
}

// ConcurrencyResult holds the outcome of a single profile run.
type ConcurrencyResult struct {
	Profile string
	Err     error
}

// RunConcurrent executes fn for each profile name concurrently, respecting
// MaxWorkers and FailFast semantics. It returns a slice of ConcurrencyResult
// in the order results arrive (not necessarily input order).
func RunConcurrent(ctx context.Context, profiles []string, cfg ConcurrencyConfig, fn func(ctx context.Context, profile string) error) []ConcurrencyResult {
	if cfg.MaxWorkers < 1 {
		cfg.MaxWorkers = 1
	}

	results := make([]ConcurrencyResult, 0, len(profiles))
	resultCh := make(chan ConcurrencyResult, len(profiles))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sem := make(chan struct{}, cfg.MaxWorkers)
	var wg sync.WaitGroup

	for _, p := range profiles {
		profile := p
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				resultCh <- ConcurrencyResult{Profile: profile, Err: fmt.Errorf("context cancelled before start: %w", ctx.Err())}
				return
			}

			err := fn(ctx, profile)
			resultCh <- ConcurrencyResult{Profile: profile, Err: err}
			if err != nil && cfg.FailFast {
				cancel()
			}
		}()
	}

	wg.Wait()
	close(resultCh)

	for r := range resultCh {
		results = append(results, r)
	}
	return results
}
