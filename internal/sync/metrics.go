package sync

import (
	"sync"
	"time"
)

// Metrics tracks runtime statistics for sync operations across profiles.
type Metrics struct {
	mu sync.Mutex

	TotalRuns     int
	SuccessCount  int
	FailureCount  int
	SkippedCount  int
	TotalDuration time.Duration
	LastRunAt     time.Time
	ProfileStats  map[string]*ProfileMetric
}

// ProfileMetric holds per-profile statistics.
type ProfileMetric struct {
	Name        string
	Runs        int
	Failures    int
	LastError   string
	LastRunAt   time.Time
	AvgDuration time.Duration
	totalDur    time.Duration
}

// NewMetrics returns an initialised Metrics instance.
func NewMetrics() *Metrics {
	return &Metrics{
		ProfileStats: make(map[string]*ProfileMetric),
	}
}

// RecordSuccess records a successful sync for the given profile.
func (m *Metrics) RecordSuccess(profile string, dur time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TotalRuns++
	m.SuccessCount++
	m.TotalDuration += dur
	m.LastRunAt = time.Now()

	pm := m.getOrCreate(profile)
	pm.Runs++
	pm.totalDur += dur
	pm.AvgDuration = pm.totalDur / time.Duration(pm.Runs)
	pm.LastRunAt = m.LastRunAt
}

// RecordFailure records a failed sync for the given profile.
func (m *Metrics) RecordFailure(profile string, dur time.Duration, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TotalRuns++
	m.FailureCount++
	m.TotalDuration += dur
	m.LastRunAt = time.Now()

	pm := m.getOrCreate(profile)
	pm.Runs++
	pm.Failures++
	pm.totalDur += dur
	pm.AvgDuration = pm.totalDur / time.Duration(pm.Runs)
	pm.LastRunAt = m.LastRunAt
	if err != nil {
		pm.LastError = err.Error()
	}
}

// RecordSkipped increments the skipped counter.
func (m *Metrics) RecordSkipped(profile string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SkippedCount++
	m.getOrCreate(profile).Runs++
}

// Snapshot returns a shallow copy of the current metrics (safe for reading).
func (m *Metrics) Snapshot() Metrics {
	m.mu.Lock()
	defer m.mu.Unlock()
	copy := *m
	return copy
}

func (m *Metrics) getOrCreate(profile string) *ProfileMetric {
	if pm, ok := m.ProfileStats[profile]; ok {
		return pm
	}
	pm := &ProfileMetric{Name: profile}
	m.ProfileStats[profile] = pm
	return pm
}
