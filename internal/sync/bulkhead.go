package sync

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ErrBulkheadFull is returned when the bulkhead has no available capacity.
var ErrBulkheadFull = errors.New("bulkhead: capacity exceeded")

// BulkheadConfig holds configuration for a Bulkhead.
type BulkheadConfig struct {
	// MaxConcurrent is the maximum number of concurrent executions allowed.
	MaxConcurrent int
	// MaxWaiting is the maximum number of callers that may queue waiting for a slot.
	MaxWaiting int
}

// DefaultBulkheadConfig returns a BulkheadConfig with sensible defaults.
func DefaultBulkheadConfig() BulkheadConfig {
	return BulkheadConfig{
		MaxConcurrent: 10,
		MaxWaiting:    20,
	}
}

// Bulkhead limits concurrent and queued executions to isolate failures
// and prevent resource exhaustion across profiles.
type Bulkhead struct {
	mu      sync.Mutex
	cfg     BulkheadConfig
	active  int
	waiting int
}

// NewBulkhead creates a Bulkhead. Zero values in cfg fall back to defaults.
func NewBulkhead(cfg BulkheadConfig) *Bulkhead {
	def := DefaultBulkheadConfig()
	if cfg.MaxConcurrent <= 0 {
		cfg.MaxConcurrent = def.MaxConcurrent
	}
	if cfg.MaxWaiting < 0 {
		cfg.MaxWaiting = def.MaxWaiting
	}
	return &Bulkhead{cfg: cfg}
}

// Execute runs fn within the bulkhead constraints. It returns ErrBulkheadFull
// when neither an active slot nor a waiting slot is available.
func (b *Bulkhead) Execute(ctx context.Context, fn func(context.Context) error) error {
	if err := b.acquire(ctx); err != nil {
		return err
	}
	defer b.release()
	return fn(ctx)
}

func (b *Bulkhead) acquire(ctx context.Context) error {
	b.mu.Lock()
	if b.active < b.cfg.MaxConcurrent {
		b.active++
		b.mu.Unlock()
		return nil
	}
	if b.waiting >= b.cfg.MaxWaiting {
		b.mu.Unlock()
		return fmt.Errorf("%w: active=%d waiting=%d", ErrBulkheadFull, b.active, b.waiting)
	}
	b.waiting++
	b.mu.Unlock()

	// Poll until a slot opens or context is cancelled.
	for {
		select {
		case <-ctx.Done():
			b.mu.Lock()
			b.waiting--
			b.mu.Unlock()
			return ctx.Err()
		default:
			b.mu.Lock()
			if b.active < b.cfg.MaxConcurrent {
				b.active++
				b.waiting--
				b.mu.Unlock()
				return nil
			}
			b.mu.Unlock()
		}
	}
}

func (b *Bulkhead) release() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.active > 0 {
		b.active--
	}
}

// Stats returns a snapshot of current active and waiting counts.
func (b *Bulkhead) Stats() (active, waiting int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.active, b.waiting
}
