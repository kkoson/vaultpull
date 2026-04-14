package sync

import (
	"math"
	"sync"
	"time"
)

// DecayConfig controls exponential decay behaviour for weighted metrics.
type DecayConfig struct {
	// HalfLife is the time after which a value decays to 50% of its original weight.
	HalfLife time.Duration
}

// DefaultDecayConfig returns a DecayConfig with sensible defaults.
func DefaultDecayConfig() DecayConfig {
	return DecayConfig{
		HalfLife: 5 * time.Minute,
	}
}

// DecayCounter tracks an exponentially decaying count over time.
type DecayCounter struct {
	mu     sync.Mutex
	cfg    DecayConfig
	value  float64
	lastAt time.Time
}

// NewDecayCounter creates a new DecayCounter using the provided config.
// If cfg.HalfLife is zero the default is used.
func NewDecayCounter(cfg DecayConfig) *DecayCounter {
	if cfg.HalfLife <= 0 {
		cfg = DefaultDecayConfig()
	}
	return &DecayCounter{cfg: cfg}
}

// Add increments the counter by delta after decaying the existing value.
func (d *DecayCounter) Add(delta float64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := time.Now()
	if !d.lastAt.IsZero() {
		elapsed := now.Sub(d.lastAt)
		decay := math.Pow(0.5, elapsed.Seconds()/d.cfg.HalfLife.Seconds())
		d.value *= decay
	}
	d.value += delta
	d.lastAt = now
}

// Value returns the current decayed value.
func (d *DecayCounter) Value() float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.lastAt.IsZero() {
		return 0
	}
	elapsed := time.Since(d.lastAt)
	decay := math.Pow(0.5, elapsed.Seconds()/d.cfg.HalfLife.Seconds())
	return d.value * decay
}

// Reset zeroes the counter.
func (d *DecayCounter) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.value = 0
	d.lastAt = time.Time{}
}
