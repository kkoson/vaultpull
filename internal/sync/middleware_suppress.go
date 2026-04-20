package sync

import (
	"context"
	"fmt"

	"github.com/yourusername/vaultpull/internal/config"
)

// WithSuppress returns a StageFunc middleware that skips execution for
// profiles whose names appear in the given Suppressor.
// If s is nil, the middleware is a no-op pass-through.
func WithSuppress(s *Suppressor) func(StageFunc) StageFunc {
	return func(next StageFunc) StageFunc {
		return func(ctx context.Context, p config.Profile) error {
			if s != nil && s.IsSuppressed(p.Name) {
				return fmt.Errorf("profile %q is suppressed: %w", p.Name, ErrSuppressed)
			}
			return next(ctx, p)
		}
	}
}

// ErrSuppressed is returned when a profile is skipped due to suppression.
var ErrSuppressed = fmt.Errorf("suppressed")
