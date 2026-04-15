package sync

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
)

func buildObserveProfile(name string) config.Profile {
	return config.Profile{Name: name, VaultPath: "secret/data/test", OutputFile: "/tmp/test.env"}
}

func TestWithObserver_PanicsOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil observer")
		}
	}()
	WithObserver(nil)
}

func TestWithObserver_RecordsSyncedOnSuccess(t *testing.T) {
	var buf bytes.Buffer
	obs := NewObserver(ObserveSummary, &buf)
	mw := WithObserver(obs)
	stage := mw(func(_ context.Context, _ config.Profile) error { return nil })

	if err := stage(context.Background(), buildObserveProfile("prod")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	events := obs.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Change != "synced" {
		t.Errorf("expected 'synced', got %q", events[0].Change)
	}
}

func TestWithObserver_RecordsFailedOnError(t *testing.T) {
	var buf bytes.Buffer
	obs := NewObserver(ObserveSummary, &buf)
	mw := WithObserver(obs)
	stage := mw(func(_ context.Context, _ config.Profile) error {
		return errors.New("vault unavailable")
	})

	err := stage(context.Background(), buildObserveProfile("staging"))
	if err == nil {
		t.Fatal("expected error to propagate")
	}
	events := obs.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Change != "failed" {
		t.Errorf("expected 'failed', got %q", events[0].Change)
	}
}

func TestWithObserver_StoresInContext(t *testing.T) {
	var buf bytes.Buffer
	obs := NewObserver(ObserveFull, &buf)
	mw := WithObserver(obs)

	var captured *Observer
	stage := mw(func(ctx context.Context, _ config.Profile) error {
		captured = ObserverFromContext(ctx)
		return nil
	})

	_ = stage(context.Background(), buildObserveProfile("dev"))
	if captured == nil {
		t.Fatal("expected observer in context")
	}
	if captured != obs {
		t.Fatal("expected same observer instance")
	}
}

func TestObserverFromContext_NilWhenMissing(t *testing.T) {
	obs := ObserverFromContext(context.Background())
	if obs != nil {
		t.Fatal("expected nil when observer not in context")
	}
}
