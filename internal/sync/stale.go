package sync

import (
	"sync"
	"time"
)

// StalenessConfig holds configuration for stale profile detection.
type StalenessConfig struct {
	// MaxAge is the maximum time a profile result is considered fresh.
	MaxAge time.Duration
	// WarnOnly logs a warning instead of returning an error when stale.
	WarnOnly bool
}

// DefaultStalenessConfig returns a StalenessConfig with sensible defaults.
func DefaultStalenessConfig() StalenessConfig {
	return StalenessConfig{
		MaxAge:   30 * time.Minute,
		WarnOnly: false,
	}
}

// StalenessTracker tracks the last successful sync time per profile.
type StalenessTracker struct {
	mu      sync.RWMutex
	cfg     StalenessConfig
	lastRun map[string]time.Time
}

// NewStalenessTracker creates a StalenessTracker with the given config.
// Zero-value MaxAge falls back to the default.
func NewStalenessTracker(cfg StalenessConfig) *StalenessTracker {
	if cfg.MaxAge <= 0 {
		cfg.MaxAge = DefaultStalenessConfig().MaxAge
	}
	return &StalenessTracker{
		cfg:     cfg,
		lastRun: make(map[string]time.Time),
	}
}

// Record marks the current time as the last successful sync for profile.
func (s *StalenessTracker) Record(profile string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastRun[profile] = time.Now()
}

// IsStale reports whether the profile's last sync exceeds MaxAge.
// Profiles that have never been synced are always considered stale.
func (s *StalenessTracker) IsStale(profile string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.lastRun[profile]
	if !ok {
		return true
	}
	return time.Since(t) > s.cfg.MaxAge
}

// LastRun returns the last recorded sync time for profile and whether it exists.
func (s *StalenessTracker) LastRun(profile string) (time.Time, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.lastRun[profile]
	return t, ok
}

// Reset removes the recorded sync time for profile.
func (s *StalenessTracker) Reset(profile string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.lastRun, profile)
}
