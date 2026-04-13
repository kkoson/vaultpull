package sync

import (
	"math/rand"
	"time"
)

// JitterConfig holds configuration for jitter applied to delays.
type JitterConfig struct {
	// Factor is the fraction of the base delay to use as the jitter range.
	// For example, 0.25 means ±25% of the base delay. Must be in [0, 1].
	Factor float64

	// MaxJitter caps the absolute jitter duration regardless of Factor.
	// Zero means no cap.
	MaxJitter time.Duration
}

// DefaultJitterConfig returns a JitterConfig with sensible defaults.
func DefaultJitterConfig() JitterConfig {
	return JitterConfig{
		Factor:    0.2,
		MaxJitter: 5 * time.Second,
	}
}

// Jitter applies randomised jitter to base according to cfg.
// The returned duration is base ± (base * Factor), capped by MaxJitter when set.
// A nil or zero-value cfg falls back to DefaultJitterConfig.
func Jitter(base time.Duration, cfg *JitterConfig) time.Duration {
	if base <= 0 {
		return 0
	}

	c := DefaultJitterConfig()
	if cfg != nil {
		if cfg.Factor >= 0 && cfg.Factor <= 1 {
			c.Factor = cfg.Factor
		}
		if cfg.MaxJitter >= 0 {
			c.MaxJitter = cfg.MaxJitter
		}
	}

	if c.Factor == 0 {
		return base
	}

	// window is the full jitter range (2 * half-range)
	window := float64(base) * c.Factor
	if c.MaxJitter > 0 && time.Duration(window) > c.MaxJitter {
		window = float64(c.MaxJitter)
	}

	// random offset in [-window/2, +window/2]
	offset := (rand.Float64() - 0.5) * window
	result := time.Duration(float64(base) + offset)
	if result < 0 {
		return 0
	}
	return result
}
