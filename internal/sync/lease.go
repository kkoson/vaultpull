package sync

import (
	"sync"
	"time"
)

// LeaseConfig holds configuration for the lease manager.
type LeaseConfig struct {
	// TTL is the duration a lease is considered valid.
	TTL time.Duration
	// RenewThreshold is the fraction of TTL remaining that triggers renewal.
	RenewThreshold float64
}

// DefaultLeaseConfig returns a LeaseConfig with sensible defaults.
func DefaultLeaseConfig() LeaseConfig {
	return LeaseConfig{
		TTL:            5 * time.Minute,
		RenewThreshold: 0.25,
	}
}

// Lease represents a time-bounded lease for a profile sync operation.
type Lease struct {
	Profile   string
	AcquiredAt time.Time
	ExpiresAt  time.Time
}

// IsExpired reports whether the lease has passed its expiry time.
func (l Lease) IsExpired(now time.Time) bool {
	return now.After(l.ExpiresAt)
}

// NeedsRenewal reports whether the lease should be renewed based on the threshold.
func (l Lease) NeedsRenewal(now time.Time, threshold float64) bool {
	total := l.ExpiresAt.Sub(l.AcquiredAt)
	remaining := l.ExpiresAt.Sub(now)
	if total <= 0 {
		return true
	}
	return float64(remaining)/float64(total) <= threshold
}

// LeaseManager tracks active leases per profile.
type LeaseManager struct {
	mu     sync.Mutex
	cfg    LeaseConfig
	leases map[string]Lease
}

// NewLeaseManager creates a LeaseManager with the given config.
func NewLeaseManager(cfg LeaseConfig) *LeaseManager {
	if cfg.TTL <= 0 {
		cfg.TTL = DefaultLeaseConfig().TTL
	}
	if cfg.RenewThreshold <= 0 || cfg.RenewThreshold > 1 {
		cfg.RenewThreshold = DefaultLeaseConfig().RenewThreshold
	}
	return &LeaseManager{
		cfg:    cfg,
		leases: make(map[string]Lease),
	}
}

// Acquire creates or replaces a lease for the given profile.
func (m *LeaseManager) Acquire(profile string, now time.Time) Lease {
	m.mu.Lock()
	defer m.mu.Unlock()
	l := Lease{
		Profile:    profile,
		AcquiredAt: now,
		ExpiresAt:  now.Add(m.cfg.TTL),
	}
	m.leases[profile] = l
	return l
}

// Get returns the current lease for a profile and whether it exists.
func (m *LeaseManager) Get(profile string) (Lease, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.leases[profile]
	return l, ok
}

// Release removes the lease for a profile.
func (m *LeaseManager) Release(profile string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.leases, profile)
}

// NeedsRenewal reports whether the profile's lease should be renewed.
func (m *LeaseManager) NeedsRenewal(profile string, now time.Time) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.leases[profile]
	if !ok {
		return true
	}
	return l.NeedsRenewal(now, m.cfg.RenewThreshold)
}
