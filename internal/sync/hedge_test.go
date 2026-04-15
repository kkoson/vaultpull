package sync

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestDefaultHedgeConfig_Values(t *testing.T) {
	cfg := DefaultHedgeConfig()
	if cfg.Delay != 200*time.Millisecond {
		t.Errorf("expected 200ms delay, got %v", cfg.Delay)
	}
	if cfg.MaxHedges != 1 {
		t.Errorf("expected MaxHedges=1, got %d", cfg.MaxHedges)
	}
}

func TestHedge_SuccessOnFirstAttempt(t *testing.T) {
	cfg := HedgeConfig{Delay: 50 * time.Millisecond, MaxHedges: 1}
	v, err := Hedge(context.Background(), cfg, func(_ context.Context) (interface{}, error) {
		return "ok", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.(string) != "ok" {
		t.Errorf("expected 'ok', got %v", v)
	}
}

func TestHedge_HedgedRequestWins(t *testing.T) {
	var calls int32
	cfg := HedgeConfig{Delay: 20 * time.Millisecond, MaxHedges: 1}
	v, err := Hedge(context.Background(), cfg, func(_ context.Context) (interface{}, error) {
		n := atomic.AddInt32(&calls, 1)
		if n == 1 {
			time.Sleep(200 * time.Millisecond)
			return nil, errors.New("slow")
		}
		return "hedge", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.(string) != "hedge" {
		t.Errorf("expected 'hedge', got %v", v)
	}
}

func TestHedge_AllFail_ReturnsError(t *testing.T) {
	cfg := HedgeConfig{Delay: 10 * time.Millisecond, MaxHedges: 1}
	expected := errors.New("boom")
	_, err := Hedge(context.Background(), cfg, func(_ context.Context) (interface{}, error) {
		return nil, expected
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestHedge_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	cfg := HedgeConfig{Delay: 10 * time.Millisecond, MaxHedges: 1}
	_, err := Hedge(ctx, cfg, func(c context.Context) (interface{}, error) {
		<-c.Done()
		return nil, c.Err()
	})
	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestHedge_ZeroMaxHedgesDefaultsToOne(t *testing.T) {
	cfg := HedgeConfig{Delay: 10 * time.Millisecond, MaxHedges: 0}
	v, err := Hedge(context.Background(), cfg, func(_ context.Context) (interface{}, error) {
		return 42, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.(int) != 42 {
		t.Errorf("expected 42, got %v", v)
	}
}
