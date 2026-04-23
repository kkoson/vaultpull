package sync

import (
	"sync"
	"time"
)

// DefaultExpiryConfig returns a sensible default ExpiryConfig.
func DefaultExpiryConfig() ExpiryConfig {
	return ExpiryConfig{
		TTL:           30 * time.Minute,
		CleanupPeriod: 5 * time.Minute,
	}
}

// ExpiryConfig controls how long entries live in the ExpiryTracker.
type ExpiryConfig struct {
	TTL           time.Duration
	CleanupPeriod time.Duration
}

type expiryEntry struct {
	recordedAt time.Time
}

// ExpiryTracker tracks when each profile was last recorded and reports
// whether the entry has passed its TTL.
type ExpiryTracker struct {
	mu      sync.Mutex
	cfg     ExpiryConfig
	entries map[string]expiryEntry
}

// NewExpiryTracker creates an ExpiryTracker with the given config.
// Zero-value fields fall back to DefaultExpiryConfig.
func NewExpiryTracker(cfg ExpiryConfig) *ExpiryTracker {
	def := DefaultExpiryConfig()
	if cfg.TTL <= 0 {
		cfg.TTL = def.TTL
	}
	if cfg.CleanupPeriod <= 0 {
		cfg.CleanupPeriod = def.CleanupPeriod
	}
	return &ExpiryTracker{
		cfg:     cfg,
		entries: make(map[string]expiryEntry),
	}
}

// Record marks the given profile as active right now.
func (e *ExpiryTracker) Record(profile string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.entries[profile] = expiryEntry{recordedAt: time.Now()}
}

// IsExpired reports whether the profile's entry has exceeded the TTL.
// Profiles that have never been recorded are considered expired.
func (e *ExpiryTracker) IsExpired(profile string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	ent, ok := e.entries[profile]
	if !ok {
		return true
	}
	return time.Since(ent.recordedAt) > e.cfg.TTL
}

// Evict removes all entries whose TTL has elapsed.
func (e *ExpiryTracker) Evict() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	removed := 0
	for k, ent := range e.entries {
		if time.Since(ent.recordedAt) > e.cfg.TTL {
			delete(e.entries, k)
			removed++
		}
	}
	return removed
}

// Len returns the number of tracked entries.
func (e *ExpiryTracker) Len() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.entries)
}
