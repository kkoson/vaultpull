package sync

import (
	"sync"
	"time"
)

// RateLimitEntry tracks per-profile rate limit state.
type RateLimitEntry struct {
	Profile   string
	Allowed   int
	Denied    int
	LastReset time.Time
	Window    time.Duration
}

// RateLimitStore tracks per-profile request counts within a sliding window.
type RateLimitStore struct {
	mu      sync.Mutex
	entries map[string]*RateLimitEntry
	maxRate int
	window  time.Duration
}

// NewRateLimitStore creates a store that allows up to maxRate syncs per window.
// A zero or negative maxRate disables limiting (always allows).
func NewRateLimitStore(maxRate int, window time.Duration) *RateLimitStore {
	if window <= 0 {
		window = time.Minute
	}
	return &RateLimitStore{
		entries: make(map[string]*RateLimitEntry),
		maxRate: maxRate,
		window:  window,
	}
}

// Allow reports whether the profile is within its rate limit and records the attempt.
func (s *RateLimitStore) Allow(profile string) bool {
	if s.maxRate <= 0 {
		return true
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	e, ok := s.entries[profile]
	if !ok || now.Sub(e.LastReset) >= s.window {
		s.entries[profile] = &RateLimitEntry{
			Profile:   profile,
			Allowed:   1,
			LastReset: now,
			Window:    s.window,
		}
		return true
	}

	if e.Allowed >= s.maxRate {
		e.Denied++
		return false
	}
	e.Allowed++
	return true
}

// Stats returns a copy of the entry for the given profile, or nil if absent.
func (s *RateLimitStore) Stats(profile string) *RateLimitEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[profile]
	if !ok {
		return nil
	}
	copy := *e
	return &copy
}

// Reset clears the rate limit state for a profile.
func (s *RateLimitStore) Reset(profile string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, profile)
}
