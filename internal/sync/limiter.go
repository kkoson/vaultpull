package sync

import (
	"context"

	"golang.org/x/time/rate"
)

// Limiter controls the rate of outbound Vault API requests.
type Limiter struct {
	r *rate.Limiter
}

// NewLimiter returns a Limiter that allows up to rps requests per second.
// If rps is zero or negative, no limiting is applied.
func NewLimiter(rps int) *Limiter {
	if rps <= 0 {
		return &Limiter{r: nil}
	}
	return &Limiter{
		r: rate.NewLimiter(rate.Limit(rps), rps),
	}
}

// Wait blocks until a token is available or the context is cancelled.
// If no rate limit is configured, Wait returns immediately.
func (l *Limiter) Wait(ctx context.Context) {
	if l == nil || l.r == nil {
		return
	}
	// Honour context cancellation; ignore the error — callers check ctx
	// themselves after Wait returns.
	_ = l.r.Wait(ctx) //nolint:errcheck
}

// Available reports whether the limiter would grant a token right now
// without blocking. Returns true when no rate limit is configured.
func (l *Limiter) Available() bool {
	if l == nil || l.r == nil {
		return true
	}
	return l.r.Allow()
}
