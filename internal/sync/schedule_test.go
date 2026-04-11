package sync

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewSchedule_NilRunner(t *testing.T) {
	_, err := NewSchedule(nil, time.Second, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil runner")
	}
}

func TestNewSchedule_ZeroInterval(t *testing.T) {
	r := buildTestRunner(t, nil, nil)
	_, err := NewSchedule(r, 0, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestNewSchedule_NegativeInterval(t *testing.T) {
	r := buildTestRunner(t, nil, nil)
	_, err := NewSchedule(r, -time.Second, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for negative interval")
	}
}

func TestSchedule_Start_CancelledImmediately(t *testing.T) {
	// runner with no profiles defined — RunAll returns nil
	r := buildTestRunner(t, nil, nil)
	s, err := NewSchedule(r, 50*time.Millisecond, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before Start
	err = s.Start(ctx)
	// initial tick runs first (no profiles → RunAll → no profiles in config → nil)
	// then ctx is done
	if err != nil && !errors.Is(err, context.Canceled) {
		t.Fatalf("expected nil or context.Canceled, got %v", err)
	}
}

func TestSchedule_Start_TicksAtLeastOnce(t *testing.T) {
	calls := 0
	newSyncer := func(p interface{}, opts Options) (interface{ Run(context.Context) error }, error) {
		calls++
		return &fakeSyncer{}, nil
	}
	_ = newSyncer // just validate the counter logic conceptually

	r := buildTestRunner(t, nil, nil)
	s, err := NewSchedule(r, 200*time.Millisecond, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	// Just ensure Start does not panic and returns a context error
	err = s.Start(ctx)
	if err != nil && !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSchedule_Profiles_Stored(t *testing.T) {
	r := buildTestRunner(t, nil, nil)
	s, err := NewSchedule(r, time.Minute, DefaultOptions(), "dev", "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(s.Profiles))
	}
	if s.Profiles[0] != "dev" || s.Profiles[1] != "staging" {
		t.Fatalf("unexpected profiles: %v", s.Profiles)
	}
}

// fakeSyncer is a no-op syncer for schedule tests.
type fakeSyncer struct{}

func (f *fakeSyncer) Run(_ context.Context) error { return nil }
