package sync

import (
	"context"
	"errors"

	"github.com/your-org/vaultpull/internal/config"
)

// ErrBlackoutActive is returned when a sync is attempted during a blackout window.
var ErrBlackoutActive = errors.New("sync suppressed: blackout window active")

// WithBlackout returns a StageFunc middleware that skips execution during
// active blackout windows. If the manager is nil the stage is always executed.
func WithBlackout(bm *BlackoutManager, next func(context.Context, config.Profile) error) func(context.Context, config.Profile) error {
	if bm == nil {
		panic("WithBlackout: BlackoutManager must not be nil")
	}
	return func(ctx context.Context, p config.Profile) error {
		if bm.IsBlackedOut() {
			return ErrBlackoutActive
		}
		return next(ctx, p)
	}
}
