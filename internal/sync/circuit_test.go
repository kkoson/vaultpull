package sync

import (
	"testing"
	"time"
)

func TestDefaultCircuitBreaker_InitialState(t *testing.T) {
	cb := DefaultCircuitBreaker()
	if cb.State() != CircuitClosed {
		t.Fatalf("expected CircuitClosed, got %d", cb.State())
	}
}

func TestCircuitBreaker_AllowWhenClosed(t *testing.T) {
	cb := DefaultCircuitBreaker()
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCircuitBreaker_OpensAfterMaxFailures(t *testing.T) {
	cb := NewCircuitBreaker(3, 10*time.Second)
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	if cb.State() != CircuitOpen {
		t.Fatalf("expected CircuitOpen, got %d", cb.State())
	}
	if err := cb.Allow(); err != ErrCircuitOpen {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreaker_RecordSuccess_Resets(t *testing.T) {
	cb := NewCircuitBreaker(1, 10*time.Second)
	cb.RecordFailure()
	if cb.State() != CircuitOpen {
		t.Fatalf("expected CircuitOpen after failure")
	}
	cb.RecordSuccess()
	if cb.State() != CircuitClosed {
		t.Fatalf("expected CircuitClosed after success")
	}
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected nil after reset, got %v", err)
	}
}

func TestCircuitBreaker_HalfOpenAfterTimeout(t *testing.T) {
	cb := NewCircuitBreaker(1, 10*time.Millisecond)
	cb.RecordFailure()
	if cb.State() != CircuitOpen {
		t.Fatalf("expected CircuitOpen")
	}
	time.Sleep(20 * time.Millisecond)
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected nil in half-open state, got %v", err)
	}
	if cb.State() != CircuitHalfOpen {
		t.Fatalf("expected CircuitHalfOpen, got %d", cb.State())
	}
}

func TestNewCircuitBreaker_ZeroMaxFailuresDefaultsToOne(t *testing.T) {
	cb := NewCircuitBreaker(0, 5*time.Second)
	cb.RecordFailure()
	if cb.State() != CircuitOpen {
		t.Fatalf("expected CircuitOpen with maxFailures=1, got %d", cb.State())
	}
}

func TestCircuitBreaker_StaysOpenWithinTimeout(t *testing.T) {
	cb := NewCircuitBreaker(1, 1*time.Hour)
	cb.RecordFailure()
	if err := cb.Allow(); err != ErrCircuitOpen {
		t.Fatalf("expected ErrCircuitOpen within timeout, got %v", err)
	}
}
