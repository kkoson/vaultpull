package sync

import (
	"context"
	"fmt"
)

// WithDecayRateLimit wraps a StageFunc and rejects execution when the
// decayed error rate for the given profile exceeds maxRate (0–1).
//
// The counter is incremented by 1 on each failure and by 0 on success,
// so the decayed value approximates a recent error frequency. When the
// value divided by the call count exceeds maxRate the stage is skipped
// and an error is returned.
func WithDecayRateLimit(counter *DecayCounter, callCounter *DecayCounter, maxRate float64, next StageFunc) StageFunc {
	if counter == nil || callCounter == nil {
		panic("WithDecayRateLimit: counter and callCounter must not be nil")
	}
	return func(ctx context.Context, profile string) error {
		calls := callCounter.Value()
		if calls > 0 {
			rate := counter.Value() / calls
			if rate > maxRate {
				return fmt.Errorf("decay: error rate %.2f exceeds limit %.2f for profile %q", rate, maxRate, profile)
			}
		}
		callCounter.Add(1)
		err := next(ctx, profile)
		if err != nil {
			counter.Add(1)
		}
		return err
	}
}
