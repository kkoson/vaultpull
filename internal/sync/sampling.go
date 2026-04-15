package sync

import (
	"math/rand"
	"sync"
	"time"
)

// SamplingMode controls how profiles are selected for sampling.
type SamplingMode int

const (
	// SamplingModeRandom selects profiles randomly based on a rate.
	SamplingModeRandom SamplingMode = iota
	// SamplingModeAlways always samples every profile.
	SamplingModeAlways
	// SamplingModeNever never samples any profile.
	SamplingModeNever
)

// DefaultSamplingConfig returns a SamplingConfig with sensible defaults.
func DefaultSamplingConfig() SamplingConfig {
	return SamplingConfig{
		Mode: SamplingModeRandom,
		Rate: 1.0,
	}
}

// SamplingConfig holds configuration for the Sampler.
type SamplingConfig struct {
	Mode SamplingMode
	Rate float64 // 0.0–1.0; only used in SamplingModeRandom
}

// Sampler decides whether a given profile sync should be sampled
// (e.g. for tracing or shadow writes).
type Sampler struct {
	mu  sync.Mutex
	cfg SamplingConfig
	rng *rand.Rand
	hits   int64
	misses int64
}

// NewSampler creates a Sampler from the provided config.
// A zero Rate defaults to 1.0 when mode is Always; for Random mode a rate of
// 0 effectively disables sampling — use SamplingModeNever instead.
func NewSampler(cfg SamplingConfig) *Sampler {
	if cfg.Rate < 0 {
		cfg.Rate = 0
	}
	if cfg.Rate > 1.0 {
		cfg.Rate = 1.0
	}
	return &Sampler{
		cfg: cfg,
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Sample returns true if the named profile should be sampled.
func (s *Sampler) Sample(profile string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	var sampled bool
	switch s.cfg.Mode {
	case SamplingModeAlways:
		sampled = true
	case SamplingModeNever:
		sampled = false
	default: // SamplingModeRandom
		sampled = s.rng.Float64() < s.cfg.Rate
	}

	if sampled {
		s.hits++
	} else {
		s.misses++
	}
	return sampled
}

// Stats returns the total number of sampled and skipped decisions.
func (s *Sampler) Stats() (hits, misses int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.hits, s.misses
}
