package sync

import (
	"math"
	"time"
)

// BackoffStrategy defines how delays are calculated between retry attempts.
type BackoffStrategy int

const (
	// BackoffFixed uses a constant delay between attempts.
	BackoffFixed BackoffStrategy = iota
	// BackoffExponential doubles the delay on each attempt.
	BackoffExponential
	// BackoffLinear increases the delay linearly on each attempt.
	BackoffLinear
)

// BackoffConfig holds configuration for a backoff policy.
type BackoffConfig struct {
	// Strategy determines how the delay grows between attempts.
	Strategy BackoffStrategy
	// BaseDelay is the initial delay duration.
	BaseDelay time.Duration
	// MaxDelay caps the computed delay so it never exceeds this value.
	MaxDelay time.Duration
	// Multiplier is used by exponential backoff (defaults to 2.0 if zero).
	Multiplier float64
}

// DefaultBackoffConfig returns a sensible exponential backoff configuration.
func DefaultBackoffConfig() BackoffConfig {
	return BackoffConfig{
		Strategy:   BackoffExponential,
		BaseDelay:  500 * time.Millisecond,
		MaxDelay:   30 * time.Second,
		Multiplier: 2.0,
	}
}

// Delay computes the wait duration for the given attempt number (0-indexed).
func (c BackoffConfig) Delay(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}

	multiplier := c.Multiplier
	if multiplier <= 0 {
		multiplier = 2.0
	}

	var d time.Duration
	switch c.Strategy {
	case BackoffExponential:
		factor := math.Pow(multiplier, float64(attempt))
		d = time.Duration(float64(c.BaseDelay) * factor)
	case BackoffLinear:
		d = c.BaseDelay * time.Duration(attempt+1)
	default: // BackoffFixed
		d = c.BaseDelay
	}

	if c.MaxDelay > 0 && d > c.MaxDelay {
		return c.MaxDelay
	}
	return d
}
