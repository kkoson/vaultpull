package sync

import (
	"context"
	"fmt"

	"github.com/yourusername/vaultpull/internal/config"
)

// ResolverMiddleware wraps a SyncFn and resolves the profile name using the
// provided Resolver before delegating to next. If the name resolves to
// multiple profiles only the first match is forwarded; callers that need
// fan-out should iterate the result of Resolver.Resolve themselves.
func ResolverMiddleware(r *Resolver, next func(ctx context.Context, p config.Profile) error) func(ctx context.Context, name string) error {
	return func(ctx context.Context, name string) error {
		if r == nil {
			return fmt.Errorf("resolver middleware: nil resolver")
		}

		profiles, err := r.Resolve(name)
		if err != nil {
			return fmt.Errorf("resolver middleware: %w", err)
		}

		var lastErr error
		for _, p := range profiles {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if err := next(ctx, p); err != nil {
				lastErr = err
			}
		}
		return lastErr
	}
}
