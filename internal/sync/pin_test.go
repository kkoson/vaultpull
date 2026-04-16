package sync

import (
	"testing"
)

func TestNewPinStore_NotNil(t *testing.T) {
	if NewPinStore() == nil {
		t.Fatal("expected non-nil PinStore")
	}
}

func TestPinStore_GetMiss(t *testing.T) {
	s := NewPinStore()
	ok, ver := s.Get("prod")
	if ok || ver != "" {
		t.Fatalf("expected miss, got ok=%v ver=%q", ok, ver)
	}
}

func TestPinStore_PinAndGet(t *testing.T) {
	s := NewPinStore()
	s.Pin("prod", "v5")
	ok, ver := s.Get("prod")
	if !ok {
		t.Fatal("expected pin to be found")
	}
	if ver != "v5" {
		t.Fatalf("expected v5, got %q", ver)
	}
}

func TestPinStore_OverwritePin(t *testing.T) {
	s := NewPinStore()
	s.Pin("prod", "v1")
	s.Pin("prod", "v2")
	_, ver := s.Get("prod")
	if ver != "v2" {
		t.Fatalf("expected v2, got %q", ver)
	}
}

func TestPinStore_Unpin(t *testing.T) {
	s := NewPinStore()
	s.Pin("staging", "v3")
	s.Unpin("staging")
	ok, _ := s.Get("staging")
	if ok {
		t.Fatal("expected pin to be removed")
	}
}

func TestPinStore_Unpin_Noop(t *testing.T) {
	s := NewPinStore()
	s.Unpin("nonexistent") // should not panic
}

func TestPinStore_Pinned_Snapshot(t *testing.T) {
	s := NewPinStore()
	s.Pin("a", "v1")
	s.Pin("b", "v2")
	snap := s.Pinned()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
	if snap["a"] != "v1" || snap["b"] != "v2" {
		t.Fatalf("unexpected snapshot: %v", snap)
	}
}

func TestPinStore_Pinned_IsCopy(t *testing.T) {
	s := NewPinStore()
	s.Pin("x", "v9")
	snap := s.Pinned()
	snap["x"] = "mutated"
	_, ver := s.Get("x")
	if ver != "v9" {
		t.Fatal("Pinned snapshot should be a copy")
	}
}
