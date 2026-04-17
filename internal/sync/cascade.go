package sync

import (
	"context"
	"fmt"
	"sync"
)

// CascadeRunner runs a primary profile and, on success, triggers a
// set of dependent profiles in parallel.
type CascadeRunner struct {
	runner  Runner
	deps    map[string][]string // profile -> dependents
	mu      sync.RWMutex
}

// Runner is the minimal interface CascadeRunner depends on.
type Runner interface {
	RunProfile(ctx context.Context, name string) error
}

// NewCascadeRunner creates a CascadeRunner wrapping the given runner.
func NewCascadeRunner(r Runner) *CascadeRunner {
	if r == nil {
		panic("cascade: runner must not be nil")
	}
	return &CascadeRunner{
		runner: r,
		deps:   make(map[string][]string),
	}
}

// AddDependency registers dependent as a profile that should run
// after primary completes successfully.
func (c *CascadeRunner) AddDependency(primary, dependent string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.deps[primary] = append(c.deps[primary], dependent)
}

// Run executes primary and then all registered dependents concurrently.
func (c *CascadeRunner) Run(ctx context.Context, primary string) error {
	if err := c.runner.RunProfile(ctx, primary); err != nil {
		return fmt.Errorf("cascade: primary %q failed: %w", primary, err)
	}

	c.mu.RLock()
	deps := append([]string(nil), c.deps[primary]...)
	c.mu.RUnlock()

	if len(deps) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(deps))
	for _, dep := range deps {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			if err := c.runner.RunProfile(ctx, name); err != nil {
				errCh <- fmt.Errorf("cascade: dependent %q failed: %w", name, err)
			}
		}(dep)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}
	return nil
}
