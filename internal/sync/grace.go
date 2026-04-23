package sync

import (
	"sync"
	"time"
)

// DefaultGraceConfig returns a GraceConfig with sensible defaults.
func DefaultGraceConfig() GraceConfig {
	return GraceConfig{
		Period: 5 * time.Second,
	}
}

// GraceConfig controls the grace period behaviour.
type GraceConfig struct {
	// Period is the duration to wait before marking a profile as failed after
	// its first error. Zero disables the grace period.
	Period time.Duration
}

// GraceManager tracks per-profile first-failure timestamps and decides whether
// the profile is still within its grace period.
type GraceManager struct {
	mu     sync.Mutex
	cfg    GraceConfig
	first  map[string]time.Time
	nowFn  func() time.Time
}

// NewGraceManager creates a GraceManager using cfg. If cfg.Period is zero the
// default config is used.
func NewGraceManager(cfg GraceConfig) *GraceManager {
	if cfg.Period == 0 {
		cfg = DefaultGraceConfig()
	}
	return &GraceManager{
		cfg:   cfg,
		first: make(map[string]time.Time),
		nowFn: time.Now,
	}
}

// RecordFailure notes the first failure time for profile. Subsequent calls for
// the same profile before a Reset are no-ops.
func (g *GraceManager) RecordFailure(profile string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.first[profile]; !ok {
		g.first[profile] = g.nowFn()
	}
}

// InGrace reports whether profile is still within its grace period.
// Returns false if no failure has been recorded.
func (g *GraceManager) InGrace(profile string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	t, ok := g.first[profile]
	if !ok {
		return false
	}
	return g.nowFn().Before(t.Add(g.cfg.Period))
}

// Reset clears the recorded failure for profile, allowing the grace period to
// restart on the next failure.
func (g *GraceManager) Reset(profile string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.first, profile)
}
