package sync

import (
	"sync"
	"time"
)

// CooldownConfig holds configuration for the cooldown manager.
type CooldownConfig struct {
	// Duration is the minimum time between successive syncs for the same profile.
	Duration time.Duration
}

// DefaultCooldownConfig returns a CooldownConfig with sensible defaults.
func DefaultCooldownConfig() CooldownConfig {
	return CooldownConfig{
		Duration: 30 * time.Second,
	}
}

// CooldownManager tracks the last sync time per profile and enforces a
// minimum interval between successive syncs.
type CooldownManager struct {
	mu       sync.Mutex
	cfg      CooldownConfig
	lastSync map[string]time.Time
}

// NewCooldownManager creates a new CooldownManager. If cfg.Duration is zero
// the default duration is used.
func NewCooldownManager(cfg CooldownConfig) *CooldownManager {
	if cfg.Duration <= 0 {
		cfg.Duration = DefaultCooldownConfig().Duration
	}
	return &CooldownManager{
		cfg:      cfg,
		lastSync: make(map[string]time.Time),
	}
}

// Allow returns true when the profile is not within the cooldown window.
// It records the current time as the last sync time when allowed.
func (c *CooldownManager) Allow(profile string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	if last, ok := c.lastSync[profile]; ok {
		if now.Sub(last) < c.cfg.Duration {
			return false
		}
	}
	c.lastSync[profile] = now
	return true
}

// Reset clears the recorded sync time for the given profile, allowing an
// immediate subsequent sync.
func (c *CooldownManager) Reset(profile string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.lastSync, profile)
}

// Remaining returns the time left in the cooldown window for a profile.
// It returns 0 if the profile is not in cooldown.
func (c *CooldownManager) Remaining(profile string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	last, ok := c.lastSync[profile]
	if !ok {
		return 0
	}
	remaining := c.cfg.Duration - time.Since(last)
	if remaining < 0 {
		return 0
	}
	return remaining
}
