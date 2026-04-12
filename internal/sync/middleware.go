package sync

import (
	"context"
	"fmt"
	"time"
)

// Middleware wraps a stage function, adding cross-cutting behaviour such as
// logging elapsed time or injecting a per-stage timeout.
type Middleware func(StageFunc) StageFunc

// StageFunc is the unit of work executed by a pipeline stage.
type StageFunc func(ctx context.Context) error

// Chain applies a slice of middlewares to fn, outermost first.
//
//	Chain(fn, m1, m2) → m1(m2(fn))
func Chain(fn StageFunc, middlewares ...Middleware) StageFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		fn = middlewares[i](fn)
	}
	return fn
}

// WithTiming returns a Middleware that records the elapsed time of a stage
// and writes it to the supplied recorder function.
func WithTiming(record func(name string, d time.Duration)) Middleware {
	return func(next StageFunc) StageFunc {
		return func(ctx context.Context) error {
			start := time.Now()
			err := next(ctx)
			if record != nil {
				record("stage", time.Since(start))
			}
			return err
		}
	}
}

// WithStageTimeout returns a Middleware that cancels the stage context after d.
// If d is zero the middleware is a no-op pass-through.
func WithStageTimeout(d time.Duration) Middleware {
	return func(next StageFunc) StageFunc {
		return func(ctx context.Context) error {
			if d <= 0 {
				return next(ctx)
			}
			tctx, cancel := context.WithTimeout(ctx, d)
			defer cancel()
			return next(tctx)
		}
	}
}

// WithRecover returns a Middleware that converts any panic in the stage into a
// non-nil error, preventing the whole process from crashing.
func WithRecover() Middleware {
	return func(next StageFunc) StageFunc {
		return func(ctx context.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("stage panic: %v", r)
				}
			}()
			return next(ctx)
		}
	}
}
