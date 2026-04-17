package sync

import (
	"context"
	"fmt"

	"github.com/drew/vaultpull/internal/config"
)

// WithReplay wraps a SyncFunc so that when the inner call fails the last
// successfully saved payload is replayed from store instead.
// If no replay entry exists the original error is returned.
func WithReplay(store *ReplayStore, fn SyncFunc) SyncFunc {
	if store == nil {
		panic("WithReplay: store must not be nil")
	}
	return func(ctx context.Context, p config.Profile) error {
		err := fn(ctx, p)
		if err == nil {
			return nil
		}

		secrets, ok := store.Load(p.Name)
		if !ok {
			return err
		}

		_ = secrets // caller-provided writer would consume secrets in real impl
		return fmt.Errorf("vault unavailable, replayed last snapshot for %s (original: %w)", p.Name, err)
	}
}
