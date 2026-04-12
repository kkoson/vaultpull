package sync

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNewWatchdog_DefaultDeadline(t *testing.T) {
	w := NewWatchdog(0, nil)
	if w.deadline != 30*time.Second {
		t.Fatalf("expected 30s deadline, got %v", w.deadline)
	}
}

func TestNewWatchdog_NilRestartNoopSafe(t *testing.T) {
	w := NewWatchdog(100*time.Millisecond, nil)
	w.Beat("worker")
	// calling restart should not panic
	w.restart("worker")
}

func TestWatchdog_Beat_UpdatesTimestamp(t *testing.T) {
	w := NewWatchdog(time.Second, nil)
	before := time.Now()
	w.Beat("w1")
	w.mu.Lock()
	last := w.heartbeat["w1"]
	w.mu.Unlock()
	if last.Before(before) {
		t.Fatal("heartbeat timestamp not updated")
	}
}

func TestWatchdog_Unregister_RemovesWorker(t *testing.T) {
	w := NewWatchdog(time.Second, nil)
	w.Beat("w1")
	w.Unregister("w1")
	w.mu.Lock()
	_, ok := w.heartbeat["w1"]
	w.mu.Unlock()
	if ok {
		t.Fatal("worker should have been removed")
	}
}

func TestWatchdog_StaleWorker_TriggersRestart(t *testing.T) {
	var mu sync.Mutex
	restarted := []string{}

	w := NewWatchdog(50*time.Millisecond, func(name string) {
		mu.Lock()
		restarted = append(restarted, name)
		mu.Unlock()
	})

	w.Beat("stale")
	// back-date the heartbeat so it appears stale
	w.mu.Lock()
	w.heartbeat["stale"] = time.Now().Add(-200 * time.Millisecond)
	w.mu.Unlock()

	w.check(time.Now())

	mu.Lock()
	defer mu.Unlock()
	if len(restarted) != 1 || restarted[0] != "stale" {
		t.Fatalf("expected restart for 'stale', got %v", restarted)
	}
}

func TestWatchdog_FreshWorker_NoRestart(t *testing.T) {
	restarted := false
	w := NewWatchdog(50*time.Millisecond, func(string) { restarted = true })
	w.Beat("fresh")
	w.check(time.Now())
	if restarted {
		t.Fatal("fresh worker should not trigger restart")
	}
}

func TestWatchdog_Start_CancelStopsLoop(t *testing.T) {
	w := NewWatchdog(time.Second, nil)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		w.Start(ctx, 10*time.Millisecond)
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("watchdog did not stop after context cancel")
	}
}
