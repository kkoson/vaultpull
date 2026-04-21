package sync

import (
	"context"
	"fmt"
	"sync"
)

// PreSyncCheck is a function that validates conditions before a sync operation.
type PreSyncCheck func(ctx context.Context, profileName string) error

// PreSyncGuard runs a set of named checks before allowing a sync to proceed.
// All checks must pass; the first failure aborts the guard.
type PreSyncGuard struct {
	mu     sync.RWMutex
	checks []namedCheck
}

type namedCheck struct {
	name string
	fn   PreSyncCheck
}

// NewPreSyncGuard returns an empty PreSyncGuard.
func NewPreSyncGuard() *PreSyncGuard {
	return &PreSyncGuard{}
}

// Register adds a named check to the guard. Checks are evaluated in
// registration order.
func (g *PreSyncGuard) Register(name string, fn PreSyncCheck) {
	if fn == nil {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.checks = append(g.checks, namedCheck{name: name, fn: fn})
}

// Run evaluates all registered checks for the given profile. It returns the
// first error encountered, wrapped with the check name for context.
func (g *PreSyncGuard) Run(ctx context.Context, profileName string) error {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, c := range g.checks {
		if err := c.fn(ctx, profileName); err != nil {
			return fmt.Errorf("pre-sync check %q failed for profile %q: %w", c.name, profileName, err)
		}
	}
	return nil
}

// Len returns the number of registered checks.
func (g *PreSyncGuard) Len() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.checks)
}

// WithPreSyncGuard returns a StageFunc middleware that runs the guard before
// invoking the next stage function.
func WithPreSyncGuard(guard *PreSyncGuard) func(StageFunc) StageFunc {
	if guard == nil {
		panic("presync: guard must not be nil")
	}
	return func(next StageFunc) StageFunc {
		return func(ctx context.Context, p ProfileContext) error {
			if err := guard.Run(ctx, p.Name); err != nil {
				return err
			}
			return next(ctx, p)
		}
	}
}
