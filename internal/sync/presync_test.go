package sync

import (
	"context"
	"errors"
	"testing"
)

func TestNewPreSyncGuard_NotNil(t *testing.T) {
	g := NewPreSyncGuard()
	if g == nil {
		t.Fatal("expected non-nil guard")
	}
}

func TestPreSyncGuard_Len_Empty(t *testing.T) {
	g := NewPreSyncGuard()
	if g.Len() != 0 {
		t.Fatalf("expected 0 checks, got %d", g.Len())
	}
}

func TestPreSyncGuard_Register_NilIgnored(t *testing.T) {
	g := NewPreSyncGuard()
	g.Register("noop", nil)
	if g.Len() != 0 {
		t.Fatalf("expected nil check to be ignored, got len %d", g.Len())
	}
}

func TestPreSyncGuard_Register_IncrementsLen(t *testing.T) {
	g := NewPreSyncGuard()
	g.Register("check1", func(_ context.Context, _ string) error { return nil })
	g.Register("check2", func(_ context.Context, _ string) error { return nil })
	if g.Len() != 2 {
		t.Fatalf("expected 2 checks, got %d", g.Len())
	}
}

func TestPreSyncGuard_Run_AllPass(t *testing.T) {
	g := NewPreSyncGuard()
	g.Register("ok", func(_ context.Context, _ string) error { return nil })
	if err := g.Run(context.Background(), "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPreSyncGuard_Run_FirstFailureReturned(t *testing.T) {
	g := NewPreSyncGuard()
	sentinel := errors.New("vault unreachable")
	g.Register("connectivity", func(_ context.Context, _ string) error { return sentinel })
	g.Register("never", func(_ context.Context, _ string) error { return errors.New("should not run") })

	err := g.Run(context.Background(), "staging")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got: %v", err)
	}
}

func TestPreSyncGuard_Run_ErrorContainsCheckName(t *testing.T) {
	g := NewPreSyncGuard()
	g.Register("health-check", func(_ context.Context, _ string) error {
		return errors.New("down")
	})
	err := g.Run(context.Background(), "dev")
	if err == nil {
		t.Fatal("expected error")
	}
	if msg := err.Error(); len(msg) == 0 {
		t.Fatal("expected non-empty error message")
	}
}

func TestWithPreSyncGuard_PanicsOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil guard")
		}
	}()
	WithPreSyncGuard(nil)
}

func TestWithPreSyncGuard_AllowsWhenChecksPassed(t *testing.T) {
	g := NewPreSyncGuard()
	g.Register("ok", func(_ context.Context, _ string) error { return nil })

	mw := WithPreSyncGuard(g)
	called := false
	next := func(_ context.Context, _ ProfileContext) error {
		called = true
		return nil
	}
	err := mw(next)(context.Background(), ProfileContext{Name: "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected next stage to be called")
	}
}

func TestWithPreSyncGuard_BlocksWhenCheckFails(t *testing.T) {
	g := NewPreSyncGuard()
	g.Register("fail", func(_ context.Context, _ string) error {
		return errors.New("blocked")
	})

	mw := WithPreSyncGuard(g)
	called := false
	next := func(_ context.Context, _ ProfileContext) error {
		called = true
		return nil
	}
	err := mw(next)(context.Background(), ProfileContext{Name: "prod"})
	if err == nil {
		t.Fatal("expected error")
	}
	if called {
		t.Fatal("next stage should not have been called")
	}
}
