package sync

import (
	"errors"
	"testing"
)

func TestHooks_RunPreSync_NilHooks(t *testing.T) {
	var h *Hooks
	if err := h.runPreSync("dev"); err != nil {
		t.Fatalf("expected nil error for nil Hooks, got %v", err)
	}
}

func TestHooks_RunPreSync_NoCallback(t *testing.T) {
	h := &Hooks{}
	if err := h.runPreSync("dev"); err != nil {
		t.Fatalf("expected nil error when PreSync is nil, got %v", err)
	}
}

func TestHooks_RunPreSync_CallbackInvoked(t *testing.T) {
	called := false
	h := &Hooks{
		PreSync: func(profile string, event HookEvent, err error) error {
			called = true
			if profile != "staging" {
				t.Errorf("expected profile %q, got %q", "staging", profile)
			}
			if event != HookPreSync {
				t.Errorf("expected event %q, got %q", HookPreSync, event)
			}
			return nil
		},
	}
	if err := h.runPreSync("staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("PreSync hook was not called")
	}
}

func TestHooks_RunPreSync_CallbackError(t *testing.T) {
	hookErr := errors.New("hook boom")
	h := &Hooks{
		PreSync: func(_ string, _ HookEvent, _ error) error { return hookErr },
	}
	err := h.runPreSync("prod")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, hookErr) {
		t.Errorf("expected wrapped hookErr, got %v", err)
	}
}

func TestHooks_RunPostSync_NilHooks(t *testing.T) {
	var h *Hooks
	if err := h.runPostSync("dev", nil); err != nil {
		t.Fatalf("expected nil error for nil Hooks, got %v", err)
	}
}

func TestHooks_RunPostSync_ReceivesSyncError(t *testing.T) {
	syncErr := errors.New("vault unreachable")
	var received error
	h := &Hooks{
		PostSync: func(_ string, _ HookEvent, err error) error {
			received = err
			return nil
		},
	}
	if err := h.runPostSync("dev", syncErr); err != nil {
		t.Fatalf("unexpected hook error: %v", err)
	}
	if !errors.Is(received, syncErr) {
		t.Errorf("expected syncErr to be forwarded, got %v", received)
	}
}

func TestHooks_RunPostSync_CallbackError(t *testing.T) {
	hookErr := errors.New("post hook failed")
	h := &Hooks{
		PostSync: func(_ string, _ HookEvent, _ error) error { return hookErr },
	}
	err := h.runPostSync("prod", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, hookErr) {
		t.Errorf("expected wrapped hookErr, got %v", err)
	}
}
