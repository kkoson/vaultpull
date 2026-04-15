package sync

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/config"
)

func buildFenceProfile(outputFile string) config.Profile {
	return config.Profile{Name: "test", OutputFile: outputFile, VaultPath: "secret/test"}
}

func TestWithFence_PanicsOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil fence")
		}
	}()
	WithWriteFence(nil, func(_ context.Context, _ config.Profile) error { return nil })
}

func TestWithFence_AllowsFirstWrite(t *testing.T) {
	fence := NewWriteFence(DefaultFenceConfig())
	called := false
	stage := WithWriteFence(fence, func(_ context.Context, _ config.Profile) error {
		called = true
		return nil
	})
	if err := stage(context.Background(), buildFenceProfile(".env")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected inner stage to be called")
	}
}

func TestWithFence_SkipsSecondWrite(t *testing.T) {
	fence := NewWriteFence(DefaultFenceConfig())
	calls := 0
	stage := WithWriteFence(fence, func(_ context.Context, _ config.Profile) error {
		calls++
		return nil
	})
	p := buildFenceProfile(".env")
	_ = stage(context.Background(), p)
	_ = stage(context.Background(), p)
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestWithFence_PropagatesInnerError(t *testing.T) {
	fence := NewWriteFence(DefaultFenceConfig())
	want := errors.New("write failed")
	stage := WithWriteFence(fence, func(_ context.Context, _ config.Profile) error {
		return want
	})
	got := stage(context.Background(), buildFenceProfile(".env.prod"))
	if !errors.Is(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestWithFence_AllowsAfterWindowExpires(t *testing.T) {
	fence := NewWriteFence(FenceConfig{Window: 10 * time.Millisecond})
	now := time.Now()
	fence.nowFn = func() time.Time { return now }
	calls := 0
	stage := WithWriteFence(fence, func(_ context.Context, _ config.Profile) error {
		calls++
		return nil
	})
	p := buildFenceProfile(".env.staging")
	_ = stage(context.Background(), p)
	fence.nowFn = func() time.Time { return now.Add(20 * time.Millisecond) }
	_ = stage(context.Background(), p)
	if calls != 2 {
		t.Fatalf("expected 2 calls after window expired, got %d", calls)
	}
}
