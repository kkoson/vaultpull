package sync

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/drew/vaultpull/internal/config"
)

func buildReplayProfile(name string) config.Profile {
	return config.Profile{Name: name, VaultPath: "secret/data/" + name, OutputFile: "/tmp/" + name + ".env"}
}

func TestWithReplay_PanicsOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	WithReplay(nil, func(_ context.Context, _ config.Profile) error { return nil })
}

func TestWithReplay_SuccessPassesThrough(t *testing.T) {
	store := NewReplayStore(tmpReplayDir(t))
	called := false
	fn := WithReplay(store, func(_ context.Context, _ config.Profile) error {
		called = true
		return nil
	})
	if err := fn(context.Background(), buildReplayProfile("prod")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("inner fn not called")
	}
}

func TestWithReplay_NoSnapshot_ReturnsOriginalError(t *testing.T) {
	store := NewReplayStore(tmpReplayDir(t))
	origErr := errors.New("vault down")
	fn := WithReplay(store, func(_ context.Context, _ config.Profile) error { return origErr })

	err := fn(context.Background(), buildReplayProfile("prod"))
	if !errors.Is(err, origErr) {
		t.Fatalf("expected original error, got %v", err)
	}
}

func TestWithReplay_WithSnapshot_ReturnsReplayMessage(t *testing.T) {
	dir := tmpReplayDir(t)
	store := NewReplayStore(dir)
	store.Save("prod", map[string]string{"K": "v"})

	fn := WithReplay(store, func(_ context.Context, _ config.Profile) error {
		return errors.New("vault down")
	})

	err := fn(context.Background(), buildReplayProfile("prod"))
	if err == nil {
		t.Fatal("expected error wrapping replay message")
	}
	if !strings.Contains(err.Error(), "replayed last snapshot") {
		t.Errorf("unexpected message: %v", err)
	}
}
