package sync

import (
	"testing"
)

func TestNewPriorityQueue_NotNil(t *testing.T) {
	q := NewPriorityQueue()
	if q == nil {
		t.Fatal("expected non-nil PriorityQueue")
	}
}

func TestPriorityQueue_Len_Empty(t *testing.T) {
	q := NewPriorityQueue()
	if q.Len() != 0 {
		t.Fatalf("expected 0, got %d", q.Len())
	}
}

func TestPriorityQueue_Add_IncrementsLen(t *testing.T) {
	q := NewPriorityQueue()
	q.Add("alpha", PriorityNormal)
	q.Add("beta", PriorityLow)
	if q.Len() != 2 {
		t.Fatalf("expected 2, got %d", q.Len())
	}
}

func TestPriorityQueue_Sorted_HighBeforeLow(t *testing.T) {
	q := NewPriorityQueue()
	q.Add("low", PriorityLow)
	q.Add("high", PriorityHigh)
	q.Add("normal", PriorityNormal)

	got := q.Sorted()
	if got[0] != "high" || got[1] != "normal" || got[2] != "low" {
		t.Fatalf("unexpected order: %v", got)
	}
}

func TestPriorityQueue_Sorted_StableOnEqual(t *testing.T) {
	q := NewPriorityQueue()
	q.Add("first", PriorityNormal)
	q.Add("second", PriorityNormal)

	got := q.Sorted()
	if got[0] != "first" || got[1] != "second" {
		t.Fatalf("expected stable order, got %v", got)
	}
}

func TestPriorityQueue_Sorted_DoesNotMutateQueue(t *testing.T) {
	q := NewPriorityQueue()
	q.Add("a", PriorityLow)
	q.Add("b", PriorityHigh)

	_ = q.Sorted()
	if q.Len() != 2 {
		t.Fatalf("sorted mutated the queue, len=%d", q.Len())
	}
}

func TestPriorityQueue_Clear_ResetsLen(t *testing.T) {
	q := NewPriorityQueue()
	q.Add("x", PriorityHigh)
	q.Clear()
	if q.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", q.Len())
	}
}

func TestPriorityLevel_Constants(t *testing.T) {
	if PriorityLow >= PriorityNormal {
		t.Error("PriorityLow should be less than PriorityNormal")
	}
	if PriorityNormal >= PriorityHigh {
		t.Error("PriorityNormal should be less than PriorityHigh")
	}
}
