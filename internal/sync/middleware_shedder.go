package sync

import (
	"errors"

	"github.com/yourusername/vaultpull/internal/config"
)

// ErrLoadShed is returned when a request is rejected by the load shedder.
var ErrLoadShed = errors.New("request shed: service overloaded")

// WithLoadShedder wraps a SyncFunc with load-shedding logic.
// If the shedder does not admit the request, ErrLoadShed is returned
// without invoking the inner function.
func WithLoadShedder(shedder *LoadShedder) func(SyncFunc) SyncFunc {
	if shedder == nil {
		panic("WithLoadShedder: shedder must not be nil")
	}
	return func(next SyncFunc) SyncFunc {
		return func(profile config.Profile) error {
			if !shedder.Admit() {
				return ErrLoadShed
			}
			defer shedder.Release()
			return next(profile)
		}
	}
}
