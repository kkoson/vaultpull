package sync

import (
	"errors"
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
)

func TestWithLoadShedder_PanicsOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on nil shedder")
		}
	}()
	WithLoadShedder(nil)
}

func TestWithLoadShedder_AllowsWhenUnderLimit(t *testing.T) {
	shedder := NewLoadShedder(ShedderConfig{MaxLoad: 5, WindowSize: 0, DropPercent: 0.5})
	called := false
	mw := WithLoadShedder(shedder)(func(_ config.Profile) error {
		called = true
		return nil
	})
	if err := mw(config.Profile{Name: "test"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected inner function to be called")
	}
}

func TestWithLoadShedder_ShedsWhenFull(t *testing.T) {
	shedder := NewLoadShedder(ShedderConfig{MaxLoad: 1, WindowSize: 0, DropPercent: 0.5})
	// Fill the shedder without releasing
	shedder.Admit()
	mw := WithLoadShedder(shedder)(func(_ config.Profile) error {
		return nil
	})
	err := mw(config.Profile{Name: "test"})
	if !errors.Is(err, ErrLoadShed) {
		t.Errorf("expected ErrLoadShed, got %v", err)
	}
}

func TestWithLoadShedder_ReleasesAfterSuccess(t *testing.T) {
	shedder := NewLoadShedder(ShedderConfig{MaxLoad: 1, WindowSize: 0, DropPercent: 0.5})
	mw := WithLoadShedder(shedder)(func(_ config.Profile) error {
		return nil
	})
	if err := mw(config.Profile{Name: "a"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// After release, should admit again
	if err := mw(config.Profile{Name: "b"}); err != nil {
		t.Errorf("expected admission after release, got %v", err)
	}
}

func TestWithLoadShedder_ReleasesAfterError(t *testing.T) {
	shedder := NewLoadShedder(ShedderConfig{MaxLoad: 1, WindowSize: 0, DropPercent: 0.5})
	expected := errors.New("inner error")
	mw := WithLoadShedder(shedder)(func(_ config.Profile) error {
		return expected
	})
	if err := mw(config.Profile{Name: "a"}); !errors.Is(err, expected) {
		t.Fatalf("expected inner error, got %v", err)
	}
	// Should still release and admit next
	if err := mw(config.Profile{Name: "b"}); !errors.Is(err, expected) {
		t.Errorf("expected inner error on second call, got %v", err)
	}
}
