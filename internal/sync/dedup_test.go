package sync

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewDeduplicator_NotNil(t *testing.T) {
	d := NewDeduplicator()
	if d == nil {
		t.Fatal("expected non-nil Deduplicator")
	}
}

func TestDedup_SingleCall_ExecutesFn(t *testing.T) {
	d := NewDeduplicator()
	called := 0
	err := d.Do("profile-a", func() error {
		called++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called != 1 {
		t.Fatalf("expected fn called once, got %d", called)
	}
}

func TestDedup_SingleCall_PropagatesError(t *testing.T) {
	d := NewDeduplicator()
	sentinel := errors.New("vault error")
	err := d.Do("profile-b", func() error { return sentinel })
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestDedup_ConcurrentCalls_OnlyOneExecution(t *testing.T) {
	d := NewDeduplicator()
	var executions atomic.Int32
	var wg sync.WaitGroup

	const goroutines = 10
	start := make(chan struct{})

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			d.Do("shared-key", func() error { //nolint:errcheck
				executions.Add(1)
				time.Sleep(20 * time.Millisecond)
				return nil
			})
		}()
	}

	close(start)
	wg.Wait()

	if n := executions.Load(); n != 1 {
		t.Fatalf("expected exactly 1 execution, got %d", n)
	}
}

func TestDedup_InFlight_ZeroAfterCompletion(t *testing.T) {
	d := NewDeduplicator()
	done := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		d.Do("key", func() error { //nolint:errcheck
			<-done
			return nil
		})
	}()

	time.Sleep(10 * time.Millisecond)
	if d.InFlight() != 1 {
		t.Fatal("expected 1 in-flight call")
	}
	close(done)
	wg.Wait()

	if d.InFlight() != 0 {
		t.Fatal("expected 0 in-flight calls after completion")
	}
}

func TestDedup_SequentialCalls_BothExecute(t *testing.T) {
	d := NewDeduplicator()
	var count atomic.Int32

	d.Do("seq", func() error { count.Add(1); return nil }) //nolint:errcheck
	d.Do("seq", func() error { count.Add(1); return nil }) //nolint:errcheck

	if count.Load() != 2 {
		t.Fatalf("expected 2 executions for sequential calls, got %d", count.Load())
	}
}
