package sync

import (
	"context"
	"fmt"
)

// WithThrottle returns a StageFunc middleware that applies the given Throttle
// before executing the wrapped stage. If the context is cancelled while
// waiting for the throttle, the error is returned immediately.
//
// Example:
//
//	stage := WithThrottle(th)(myStageFunc)
func WithThrottle(th *Throttle) func(StageFunc) StageFunc {
	return func(next StageFunc) StageFunc {
		return func(ctx context.Context, profile string) error {
			if th == nil {
				return next(ctx, profile)
			}
			if err := th.Wait(ctx); err != nil {
				return fmt.Errorf("throttle wait for profile %q: %w", profile, err)
			}
			return next(ctx, profile)
		}
	}
}
