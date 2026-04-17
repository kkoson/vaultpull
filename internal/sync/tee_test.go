package sync

import (
	"errors"
	"testing"
)

type mockWriter struct {
	called  bool
	payload map[string]string
	err     error
}

func (m *mockWriter) Write(secrets map[string]string) error {
	m.called = true
	m.payload = secrets
	return m.err
}

func TestNewTeeWriter_PanicsOnEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for empty writers")
		}
	}()
	NewTeeWriter()
}

func TestTeeWriter_Len(t *testing.T) {
	a, b := &mockWriter{}, &mockWriter{}
	tee := NewTeeWriter(a, b)
	if tee.Len() != 2 {
		t.Fatalf("expected 2, got %d", tee.Len())
	}
}

func TestTeeWriter_WritesAllTargets(t *testing.T) {
	a, b := &mockWriter{}, &mockWriter{}
	tee := NewTeeWriter(a, b)
	secrets := map[string]string{"KEY": "val"}
	if err := tee.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !a.called || !b.called {
		t.Fatal("expected both writers to be called")
	}
	if a.payload["KEY"] != "val" || b.payload["KEY"] != "val" {
		t.Fatal("unexpected payload")
	}
}

func TestTeeWriter_ReturnsFirstError(t *testing.T) {
	errA := errors.New("writer a failed")
	a := &mockWriter{err: errA}
	b := &mockWriter{}
	tee := NewTeeWriter(a, b)
	err := tee.Write(map[string]string{})
	if !errors.Is(err, errA) {
		t.Fatalf("expected errA, got %v", err)
	}
	// second writer must still have been called
	if !b.called {
		t.Fatal("expected second writer to be called despite first error")
	}
}

func TestTeeWriter_SingleWriter(t *testing.T) {
	w := &mockWriter{}
	tee := NewTeeWriter(w)
	if err := tee.Write(map[string]string{"X": "1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !w.called {
		t.Fatal("writer was not called")
	}
}
