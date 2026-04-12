package sync

import (
	"errors"
	"testing"
	"time"
)

func TestNewMetrics_InitialState(t *testing.T) {
	m := NewMetrics()
	if m.TotalRuns != 0 {
		t.Fatalf("expected 0 total runs, got %d", m.TotalRuns)
	}
	if m.ProfileStats == nil {
		t.Fatal("expected ProfileStats map to be initialised")
	}
}

func TestMetrics_RecordSuccess(t *testing.T) {
	m := NewMetrics()
	m.RecordSuccess("dev", 100*time.Millisecond)

	if m.TotalRuns != 1 {
		t.Fatalf("expected 1 total run, got %d", m.TotalRuns)
	}
	if m.SuccessCount != 1 {
		t.Fatalf("expected 1 success, got %d", m.SuccessCount)
	}
	if m.FailureCount != 0 {
		t.Fatalf("expected 0 failures, got %d", m.FailureCount)
	}
	pm := m.ProfileStats["dev"]
	if pm == nil {
		t.Fatal("expected profile metric for 'dev'")
	}
	if pm.Runs != 1 {
		t.Fatalf("expected 1 profile run, got %d", pm.Runs)
	}
	if pm.AvgDuration != 100*time.Millisecond {
		t.Fatalf("expected avg 100ms, got %v", pm.AvgDuration)
	}
}

func TestMetrics_RecordFailure(t *testing.T) {
	m := NewMetrics()
	err := errors.New("vault unreachable")
	m.RecordFailure("prod", 50*time.Millisecond, err)

	if m.FailureCount != 1 {
		t.Fatalf("expected 1 failure, got %d", m.FailureCount)
	}
	pm := m.ProfileStats["prod"]
	if pm.Failures != 1 {
		t.Fatalf("expected 1 profile failure, got %d", pm.Failures)
	}
	if pm.LastError != "vault unreachable" {
		t.Fatalf("unexpected last error: %s", pm.LastError)
	}
}

func TestMetrics_RecordSkipped(t *testing.T) {
	m := NewMetrics()
	m.RecordSkipped("staging")

	if m.SkippedCount != 1 {
		t.Fatalf("expected 1 skipped, got %d", m.SkippedCount)
	}
	if _, ok := m.ProfileStats["staging"]; !ok {
		t.Fatal("expected profile entry for skipped profile")
	}
}

func TestMetrics_AvgDuration_MultipleRuns(t *testing.T) {
	m := NewMetrics()
	m.RecordSuccess("dev", 100*time.Millisecond)
	m.RecordSuccess("dev", 200*time.Millisecond)

	pm := m.ProfileStats["dev"]
	if pm.AvgDuration != 150*time.Millisecond {
		t.Fatalf("expected avg 150ms, got %v", pm.AvgDuration)
	}
}

func TestMetrics_Snapshot_IsIndependent(t *testing.T) {
	m := NewMetrics()
	m.RecordSuccess("dev", 10*time.Millisecond)

	snap := m.Snapshot()
	m.RecordSuccess("dev", 10*time.Millisecond)

	if snap.TotalRuns != 1 {
		t.Fatalf("snapshot should not reflect later mutations, got %d", snap.TotalRuns)
	}
}
