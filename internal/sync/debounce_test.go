package sync

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewDebouncer_DefaultDelay(t *testing.T) {
	d := NewDebouncer(0)
	if d.delay != 500*time.Millisecond {
		t.Fatalf("expected default delay 500ms, got %v", d.delay)
	}
}

func TestNewDebouncer_CustomDelay(t *testing.T) {
	d := NewDebouncer(200 * time.Millisecond)
	if d.delay != 200*time.Millisecond {
		t.Fatalf("expected 200ms, got %v", d.delay)
	}
}

func TestDebouncer_FunctionCalledAfterDelay(t *testing.T) {
	var called atomic.Int32
	d := NewDebouncer(50 * time.Millisecond)

	d.Trigger(context.Background(), func() { called.Add(1) })

	time.Sleep(120 * time.Millisecond)
	if called.Load() != 1 {
		t.Fatalf("expected fn called once, got %d", called.Load())
	}
}

func TestDebouncer_RapidTriggers_OnlyLastFires(t *testing.T) {
	var called atomic.Int32
	d := NewDebouncer(60 * time.Millisecond)

	for i := 0; i < 5; i++ {
		d.Trigger(context.Background(), func() { called.Add(1) })
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(150 * time.Millisecond)
	if called.Load() != 1 {
		t.Fatalf("expected exactly 1 call, got %d", called.Load())
	}
}

func TestDebouncer_Flush_CancelsPending(t *testing.T) {
	var called atomic.Int32
	d := NewDebouncer(80 * time.Millisecond)

	d.Trigger(context.Background(), func() { called.Add(1) })
	d.Flush()

	time.Sleep(150 * time.Millisecond)
	if called.Load() != 0 {
		t.Fatalf("expected fn not called after flush, got %d", called.Load())
	}
}

func TestDebouncer_ContextCancelled_DoesNotFire(t *testing.T) {
	var called atomic.Int32
	d := NewDebouncer(60 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	d.Trigger(ctx, func() { called.Add(1) })
	cancel()

	time.Sleep(150 * time.Millisecond)
	if called.Load() != 0 {
		t.Fatalf("expected fn not called after context cancel, got %d", called.Load())
	}
}
