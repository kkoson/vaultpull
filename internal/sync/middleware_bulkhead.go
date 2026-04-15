package sync

import (
	"context"
	"fmt"

	"github.com/yourusername/vaultpull/internal/config"
)

// WithBulkhead returns a StageFunc middleware that wraps the next stage with
// bulkhead isolation. Execution is rejected immediately when the bulkhead is
// full (no waiting slots remain).
//
// Usage:
//
//	stage := WithBulkhead(bulkhead)(myStage)
func WithBulkhead(b *Bulkhead) func(StageFunc) StageFunc {
	if b == nil {
		panic("WithBulkhead: bulkhead must not be nil")
	}
	return func(next StageFunc) StageFunc {
		return func(ctx context.Context, p config.Profile) error {
			err := b.Execute(ctx, func(ctx context.Context) error {
				return next(ctx, p)
			})
			if err != nil {
				return fmt.Errorf("bulkhead [%s]: %w", p.Name, err)
			}
			return nil
		}
	}
}
