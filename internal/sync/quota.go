package sync

import (
	"fmt"
	"sync"
	"time"
)

// QuotaConfig defines limits for sync operations within a rolling window.
type QuotaConfig struct {
	// MaxSyncsPerWindow is the maximum number of syncs allowed per window.
	MaxSyncsPerWindow int
	// Window is the duration of the rolling window.
	Window time.Duration
}

// DefaultQuotaConfig returns a QuotaConfig with sensible defaults.
func DefaultQuotaConfig() QuotaConfig {
	return QuotaConfig{
		MaxSyncsPerWindow: 100,
		Window:            time.Hour,
	}
}

// quotaEntry tracks usage for a single profile.
type quotaEntry struct {
	count     int
	windowEnd time.Time
}

// QuotaEnforcer tracks and enforces per-profile sync quotas.
type QuotaEnforcer struct {
	mu      sync.Mutex
	cfg     QuotaConfig
	entries map[string]*quotaEntry
	now     func() time.Time
}

// NewQuotaEnforcer creates a new QuotaEnforcer with the given config.
func NewQuotaEnforcer(cfg QuotaConfig) *QuotaEnforcer {
	if cfg.MaxSyncsPerWindow <= 0 {
		cfg.MaxSyncsPerWindow = DefaultQuotaConfig().MaxSyncsPerWindow
	}
	if cfg.Window <= 0 {
		cfg.Window = DefaultQuotaConfig().Window
	}
	return &QuotaEnforcer{
		cfg:     cfg,
		entries: make(map[string]*quotaEntry),
		now:     time.Now,
	}
}

// Allow checks whether a sync for the given profile is within quota.
// It increments the counter if allowed and returns an error if the quota is exceeded.
func (q *QuotaEnforcer) Allow(profile string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := q.now()
	entry, ok := q.entries[profile]
	if !ok || now.After(entry.windowEnd) {
		q.entries[profile] = &quotaEntry{
			count:     1,
			windowEnd: now.Add(q.cfg.Window),
		}
		return nil
	}

	if entry.count >= q.cfg.MaxSyncsPerWindow {
		return fmt.Errorf("quota exceeded for profile %q: %d/%d syncs in window ending %s",
			profile, entry.count, q.cfg.MaxSyncsPerWindow, entry.windowEnd.Format(time.RFC3339))
	}

	entry.count++
	return nil
}

// Stats returns the current count and window end for a profile.
// Returns zeros if no usage has been recorded.
func (q *QuotaEnforcer) Stats(profile string) (count int, windowEnd time.Time) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if entry, ok := q.entries[profile]; ok {
		return entry.count, entry.windowEnd
	}
	return 0, time.Time{}
}

// Reset clears quota tracking for the given profile.
func (q *QuotaEnforcer) Reset(profile string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.entries, profile)
}
