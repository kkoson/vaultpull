package sync

import (
	"context"
	"errors"
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
)

func TestNewSuppressor_NilAllowsAll(t *testing.T) {
	s := NewSuppressor(nil)
	if s.IsSuppressed("any") {
		t.Fatal("expected no suppression for nil input")
	}
}

func TestNewSuppressor_EmptyAllowsAll(t *testing.T) {
	s := NewSuppressor([]string{})
	if s.IsSuppressed("prod") {
		t.Fatal("expected no suppression for empty list")
	}
}

func TestSuppressor_IsSuppressed_Match(t *testing.T) {
	s := NewSuppressor([]string{"staging", "legacy"})
	if !s.IsSuppressed("staging") {
		t.Fatal("expected staging to be suppressed")
	}
	if !s.IsSuppressed("LEGACY") {
		t.Fatal("expected case-insensitive match for legacy")
	}
}

func TestSuppressor_IsSuppressed_NoMatch(t *testing.T) {
	s := NewSuppressor([]string{"staging"})
	if s.IsSuppressed("prod") {
		t.Fatal("prod should not be suppressed")
	}
}

func TestSuppressor_AddAndRemove(t *testing.T) {
	s := NewSuppressor(nil)
	s.Add("dev")
	if !s.IsSuppressed("dev") {
		t.Fatal("expected dev to be suppressed after Add")
	}
	s.Remove("dev")
	if s.IsSuppressed("dev") {
		t.Fatal("expected dev to be unsuppressed after Remove")
	}
}

func TestSuppressor_Len(t *testing.T) {
	s := NewSuppressor([]string{"a", "b", "c"})
	if s.Len() != 3 {
		t.Fatalf("expected Len 3, got %d", s.Len())
	}
}

func TestWithSuppress_AllowsUnsuppressedProfile(t *testing.T) {
	s := NewSuppressor([]string{"staging"})
	called := false
	next := func(_ context.Context, _ config.Profile) error {
		called = true
		return nil
	}
	mw := WithSuppress(s)(next)
	p := config.Profile{Name: "prod"}
	if err := mw(context.Background(), p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected next to be called")
	}
}

func TestWithSuppress_SkipsSuppressedProfile(t *testing.T) {
	s := NewSuppressor([]string{"staging"})
	next := func(_ context.Context, _ config.Profile) error {
		return nil
	}
	mw := WithSuppress(s)(next)
	p := config.Profile{Name: "staging"}
	err := mw(context.Background(), p)
	if err == nil {
		t.Fatal("expected error for suppressed profile")
	}
	if !errors.Is(err, ErrSuppressed) {
		t.Fatalf("expected ErrSuppressed, got %v", err)
	}
}

func TestWithSuppress_NilSuppressor_PassesThrough(t *testing.T) {
	called := false
	next := func(_ context.Context, _ config.Profile) error {
		called = true
		return nil
	}
	mw := WithSuppress(nil)(next)
	p := config.Profile{Name: "anything"}
	if err := mw(context.Background(), p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected next to be called with nil suppressor")
	}
}
