package sync

import (
	"context"
	"fmt"
)

// WithCircuitBreaker wraps a StageFunc with circuit breaker protection.
// If the circuit is open, the stage is skipped and an error is returned
// immediately without invoking the underlying function.
//
// This middleware is useful for preventing cascading failures when a
// downstream dependency (e.g. Vault) is repeatedly unavailable.
func WithCircuitBreaker(cb *CircuitBreaker, next StageFunc) StageFunc {
	return func(ctx context.Context, profile string) error {
		if !cb.Allow() {
			return fmt.Errorf("circuit breaker open: skipping profile %q", profile)
		}

		err := next(ctx, profile)
		if err != nil {
			cb.RecordFailure()
			return err
		}

		cb.RecordSuccess()
		return nil
	}
}
