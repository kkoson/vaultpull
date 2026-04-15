package sync

import (
	"context"

	"github.com/yourusername/vaultpull/internal/config"
)

type observerKey struct{}

// WithObserver is a pipeline middleware that records a summary observation
// for each profile execution and stores the observer in the context.
func WithObserver(obs *Observer) func(StageFunc) StageFunc {
	if obs == nil {
		panic("WithObserver: observer must not be nil")
	}
	return func(next StageFunc) StageFunc {
		return func(ctx context.Context, p config.Profile) error {
			ctx = context.WithValue(ctx, observerKey{}, obs)
			err := next(ctx, p)
			if err != nil {
				obs.Record(p.Name, "", "failed")
			} else {
				obs.Record(p.Name, "", "synced")
			}
			return err
		}
	}
}

// ObserverFromContext retrieves the Observer stored by WithObserver.
// Returns nil if none is present.
func ObserverFromContext(ctx context.Context) *Observer {
	v, _ := ctx.Value(observerKey{}).(*Observer)
	return v
}
