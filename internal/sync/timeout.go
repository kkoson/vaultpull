package sync

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// TimeoutConfig holds configuration for per-profile sync timeouts.
type TimeoutConfig struct {
	// ProfileTimeout is the maximum duration allowed for a single profile sync.
	// A zero value disables the timeout.
	ProfileTimeout time.Duration

	// GlobalTimeout is the maximum duration allowed for a full sync run across
	// all profiles. A zero value disables the timeout.
	GlobalTimeout time.Duration
}

// DefaultTimeoutConfig returns a TimeoutConfig with sensible defaults.
func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		ProfileTimeout: 30 * time.Second,
		GlobalTimeout:  5 * time.Minute,
	}
}

// WithProfileTimeout wraps the given context with a per-profile deadline when
// cfg.ProfileTimeout is greater than zero. The returned cancel function must
// always be called by the caller.
func WithProfileTimeout(ctx context.Context, cfg TimeoutConfig) (context.Context, context.CancelFunc) {
	if cfg.ProfileTimeout <= 0 {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, cfg.ProfileTimeout)
}

// WithGlobalTimeout wraps the given context with a global deadline when
// cfg.GlobalTimeout is greater than zero.
func WithGlobalTimeout(ctx context.Context, cfg TimeoutConfig) (context.Context, context.CancelFunc) {
	if cfg.GlobalTimeout <= 0 {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, cfg.GlobalTimeout)
}

// TimeoutError is returned when a sync operation exceeds its deadline.
type TimeoutError struct {
	Profile string
	Limit   time.Duration
}

func (e *TimeoutError) Error() string {
	if e.Profile != "" {
		return fmt.Sprintf("sync timed out for profile %q after %s", e.Profile, e.Limit)
	}
	return fmt.Sprintf("global sync timed out after %s", e.Limit)
}

// IsTimeout reports whether err represents a context deadline or timeout.
// It checks for context.DeadlineExceeded, context.Canceled (when the parent
// deadline propagates), and *TimeoutError values anywhere in the error chain.
func IsTimeout(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var te *TimeoutError
	return errors.As(err, &te)
}
