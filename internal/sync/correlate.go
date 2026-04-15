package sync

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type correlationKey struct{}

// Correlator generates and propagates correlation IDs through a context.
type Correlator struct {
	prefix string
}

// NewCorrelator returns a Correlator whose IDs are optionally prefixed.
func NewCorrelator(prefix string) *Correlator {
	return &Correlator{prefix: prefix}
}

// Inject creates a child context that carries a new correlation ID derived
// from the profile name and a random suffix.
func (c *Correlator) Inject(ctx context.Context, profile string) context.Context {
	id := c.newID(profile)
	return context.WithValue(ctx, correlationKey{}, id)
}

// newID builds a collision-resistant ID of the form [prefix-]profile-<hex>.
func (c *Correlator) newID(profile string) string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	suffix := hex.EncodeToString(b)
	if c.prefix != "" {
		return fmt.Sprintf("%s-%s-%s", c.prefix, profile, suffix)
	}
	return fmt.Sprintf("%s-%s", profile, suffix)
}

// CorrelationIDFromContext retrieves the correlation ID stored in ctx.
// It returns an empty string when no ID has been injected.
func CorrelationIDFromContext(ctx context.Context) string {
	if v := ctx.Value(correlationKey{}); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}
