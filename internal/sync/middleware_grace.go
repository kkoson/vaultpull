package sync

import (
	"context"
	"fmt"

	"github.com/your-org/vaultpull/internal/config"
)

// WithGrace wraps a SyncFunc so that failures within the grace period are
// suppressed — the error is swallowed and a nil is returned instead, giving
// the profile time to recover before being treated as truly failed.
//
// On success the grace record for the profile is reset so the next failure
// starts a fresh grace window.
func WithGrace(gm *GraceManager, next SyncFunc) SyncFunc {
	if gm == nil {
		panic("WithGrace: GraceManager must not be nil")
	}
	return func(ctx context.Context, p config.Profile) error {
		err := next(ctx, p)
		if err == nil {
			gm.Reset(p.Name)
			return nil
		}
		gm.RecordFailure(p.Name)
		if gm.InGrace(p.Name) {
			// Still within grace window — suppress the error.
			return nil
		}
		return fmt.Errorf("grace period expired for profile %q: %w", p.Name, err)
	}
}
