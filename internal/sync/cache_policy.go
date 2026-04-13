package sync

import "time"

// CachePolicyMode controls how cached secrets are used during sync.
type CachePolicyMode int

const (
	// CachePolicyBypass always fetches from Vault, ignoring any cached value.
	CachePolicyBypass CachePolicyMode = iota
	// CachePolicyReadThrough returns a cached value when fresh, fetches otherwise.
	CachePolicyReadThrough
	// CachePolicyWriteThrough always fetches but writes the result back to cache.
	CachePolicyWriteThrough
)

// CachePolicyConfig holds tunables for the cache policy middleware.
type CachePolicyConfig struct {
	Mode CachePolicyMode
	TTL  time.Duration
}

// DefaultCachePolicyConfig returns a sensible default: read-through with a
// five-minute TTL.
func DefaultCachePolicyConfig() CachePolicyConfig {
	return CachePolicyConfig{
		Mode: CachePolicyReadThrough,
		TTL:  5 * time.Minute,
	}
}

// CachePolicy enforces the configured caching strategy for profile syncs.
type CachePolicy struct {
	cfg   CachePolicyConfig
	store map[string]cachePolicyEntry
}

type cachePolicyEntry struct {
	fetchedAt time.Time
	secrets   map[string]string
}

// NewCachePolicy constructs a CachePolicy. If cfg.TTL is zero the default TTL
// is used.
func NewCachePolicy(cfg CachePolicyConfig) *CachePolicy {
	if cfg.TTL <= 0 {
		cfg.TTL = DefaultCachePolicyConfig().TTL
	}
	return &CachePolicy{cfg: cfg, store: make(map[string]cachePolicyEntry)}
}

// IsFresh reports whether the cached entry for profile is still within TTL.
func (p *CachePolicy) IsFresh(profile string) bool {
	if p.cfg.Mode == CachePolicyBypass {
		return false
	}
	e, ok := p.store[profile]
	if !ok {
		return false
	}
	return time.Since(e.fetchedAt) < p.cfg.TTL
}

// Get returns the cached secrets for profile, or nil when not present.
func (p *CachePolicy) Get(profile string) map[string]string {
	e, ok := p.store[profile]
	if !ok {
		return nil
	}
	out := make(map[string]string, len(e.secrets))
	for k, v := range e.secrets {
		out[k] = v
	}
	return out
}

// Set stores secrets for profile with the current timestamp.
func (p *CachePolicy) Set(profile string, secrets map[string]string) {
	snap := make(map[string]string, len(secrets))
	for k, v := range secrets {
		snap[k] = v
	}
	p.store[profile] = cachePolicyEntry{fetchedAt: time.Now(), secrets: snap}
}

// Invalidate removes the cached entry for profile.
func (p *CachePolicy) Invalidate(profile string) {
	delete(p.store, profile)
}
