package sync

import (
	"context"
	"strings"
	"testing"
)

func TestNewCorrelator_NotNil(t *testing.T) {
	c := NewCorrelator("")
	if c == nil {
		t.Fatal("expected non-nil Correlator")
	}
}

func TestCorrelationIDFromContext_EmptyWhenMissing(t *testing.T) {
	id := CorrelationIDFromContext(context.Background())
	if id != "" {
		t.Fatalf("expected empty string, got %q", id)
	}
}

func TestCorrelator_Inject_ReturnsNonEmpty(t *testing.T) {
	c := NewCorrelator("")
	ctx := c.Inject(context.Background(), "staging")
	id := CorrelationIDFromContext(ctx)
	if id == "" {
		t.Fatal("expected non-empty correlation ID")
	}
}

func TestCorrelator_Inject_ContainsProfileName(t *testing.T) {
	c := NewCorrelator("")
	ctx := c.Inject(context.Background(), "production")
	id := CorrelationIDFromContext(ctx)
	if !strings.Contains(id, "production") {
		t.Fatalf("expected ID to contain profile name, got %q", id)
	}
}

func TestCorrelator_Inject_ContainsPrefix(t *testing.T) {
	c := NewCorrelator("vp")
	ctx := c.Inject(context.Background(), "dev")
	id := CorrelationIDFromContext(ctx)
	if !strings.HasPrefix(id, "vp-") {
		t.Fatalf("expected ID to start with prefix, got %q", id)
	}
}

func TestCorrelator_Inject_UniqueIDs(t *testing.T) {
	c := NewCorrelator("")
	ctx1 := c.Inject(context.Background(), "alpha")
	ctx2 := c.Inject(context.Background(), "alpha")
	id1 := CorrelationIDFromContext(ctx1)
	id2 := CorrelationIDFromContext(ctx2)
	if id1 == id2 {
		t.Fatalf("expected unique IDs, both were %q", id1)
	}
}

func TestCorrelator_Inject_DoesNotMutateParent(t *testing.T) {
	c := NewCorrelator("")
	parent := context.Background()
	_ = c.Inject(parent, "child")
	if id := CorrelationIDFromContext(parent); id != "" {
		t.Fatalf("parent context should not carry ID, got %q", id)
	}
}
