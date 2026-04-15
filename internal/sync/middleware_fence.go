package sync

import (
	"context"
	"fmt"

	"github.com/yourusername/vaultpull/internal/config"
)

// WithWriteFence wraps a SyncFunc so that each profile's output file path
// is checked against the fence before execution. If the fence blocks the
// write, the stage is skipped and no error is returned.
func WithWriteFence(fence *WriteFence, next func(context.Context, config.Profile) error) func(context.Context, config.Profile) error {
	if fence == nil {
		panic("WithWriteFence: fence must not be nil")
	}
	return func(ctx context.Context, p config.Profile) error {
		if err := fence.Allow(p.OutputFile); err != nil {
			// Fenced — log and skip rather than fail.
			fmt.Printf("[fence] skipping profile %q: %v\n", p.Name, err)
			return nil
		}
		return next(ctx, p)
	}
}
