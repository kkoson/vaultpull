package sync

import (
	"context"
	"errors"
	"time"
)

// RetryPolicy defines how sync operations are retried on failure.
type RetryPolicy struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// Delay is the initial wait duration between attempts.
	Delay time.Duration
	// Multiplier scales the delay after each attempt (exponential backoff).
	Multiplier float64
}

// DefaultRetryPolicy returns a sensible default retry policy.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// ErrMaxAttemptsReached is returned when all retry attempts are exhausted.
var ErrMaxAttemptsReached = errors.New("max retry attempts reached")

// WithRetry executes fn according to the given RetryPolicy, returning the
// last error if all attempts fail. The context is checked before each attempt.
func WithRetry(ctx context.Context, policy RetryPolicy, fn func() error) error {
	if policy.MaxAttempts <= 0 {
		policy.MaxAttempts = 1
	}
	if policy.Multiplier <= 0 {
		policy.Multiplier = 1.0
	}

	delay := policy.Delay
	var lastErr error

	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		if lastErr = fn(); lastErr == nil {
			return nil
		}

		if attempt == policy.MaxAttempts {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		delay = time.Duration(float64(delay) * policy.Multiplier)
	}

	return errors.Join(ErrMaxAttemptsReached, lastErr)
}
