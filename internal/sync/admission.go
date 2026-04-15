package sync

import (
	"fmt"
	"sync"
	"time"
)

// AdmissionPolicy controls whether a profile sync should be admitted
// based on recent failure history and a configurable failure threshold.
type AdmissionPolicy struct {
	mu            sync.Mutex
	maxFailures   int
	window        time.Duration
	failures      map[string][]time.Time
}

// DefaultAdmissionPolicy returns an AdmissionPolicy with sensible defaults:
// up to 3 failures within a 2-minute window before a profile is denied.
func DefaultAdmissionPolicy() *AdmissionPolicy {
	return NewAdmissionPolicy(3, 2*time.Minute)
}

// NewAdmissionPolicy creates an AdmissionPolicy with the given failure
// threshold and observation window. Zero values fall back to defaults.
func NewAdmissionPolicy(maxFailures int, window time.Duration) *AdmissionPolicy {
	if maxFailures <= 0 {
		maxFailures = 3
	}
	if window <= 0 {
		window = 2 * time.Minute
	}
	return &AdmissionPolicy{
		maxFailures: maxFailures,
		window:      window,
		failures:    make(map[string][]time.Time),
	}
}

// Admit returns nil if the profile is allowed to proceed, or an error
// if the failure count within the observation window exceeds the threshold.
func (a *AdmissionPolicy) Admit(profile string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-a.window)

	recent := a.failures[profile][:0]
	for _, t := range a.failures[profile] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}
	a.failures[profile] = recent

	if len(recent) >= a.maxFailures {
		return fmt.Errorf("admission denied for profile %q: %d failures in last %s",
			profile, len(recent), a.window)
	}
	return nil
}

// RecordFailure records a failure timestamp for the given profile.
func (a *AdmissionPolicy) RecordFailure(profile string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.failures[profile] = append(a.failures[profile], time.Now())
}

// Reset clears all failure history for the given profile.
func (a *AdmissionPolicy) Reset(profile string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.failures, profile)
}
