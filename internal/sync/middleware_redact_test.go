package sync

import (
	"context"
	"errors"
	"testing"
)

func TestWithRedact_StoresRedactorInContext(t *testing.T) {
	r := NewRedactor(DefaultRedactConfig())
	var capturedCtx context.Context

	stage := func(ctx context.Context, profile string) error {
		capturedCtx = ctx
		return nil
	}

	wrapped := WithRedact(r)(stage)
	if err := wrapped(context.Background(), "dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := RedactorFromContext(capturedCtx)
	if got == nil {
		t.Fatal("expected Redactor in context, got nil")
	}
}

func TestWithRedact_PropagatesError(t *testing.T) {
	r := NewRedactor(DefaultRedactConfig())
	want := errors.New("stage error")

	stage := func(ctx context.Context, profile string) error {
		return want
	}

	wrapped := WithRedact(r)(stage)
	got := wrapped(context.Background(), "dev")
	if !errors.Is(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestRedactorFromContext_NilWhenMissing(t *testing.T) {
	got := RedactorFromContext(context.Background())
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestWithRedact_PassesProfileUnchanged(t *testing.T) {
	r := NewRedactor(DefaultRedactConfig())
	var capturedProfile string

	stage := func(ctx context.Context, profile string) error {
		capturedProfile = profile
		return nil
	}

	wrapped := WithRedact(r)(stage)
	if err := wrapped(context.Background(), "staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedProfile != "staging" {
		t.Fatalf("expected staging, got %q", capturedProfile)
	}
}

func TestWithRedact_ChainedWithOtherMiddleware(t *testing.T) {
	r := NewRedactor(DefaultRedactConfig())
	calls := 0

	base := func(ctx context.Context, profile string) error {
		calls++
		if RedactorFromContext(ctx) == nil {
			return errors.New("no redactor in context")
		}
		return nil
	}

	wrapped := WithRedact(r)(base)
	for i := 0; i < 3; i++ {
		if err := wrapped(context.Background(), "prod"); err != nil {
			t.Fatalf("call %d failed: %v", i, err)
		}
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}
