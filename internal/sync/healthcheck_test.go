package sync

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewHealthChecker_NotNil(t *testing.T) {
	h := NewHealthChecker()
	if h == nil {
		t.Fatal("expected non-nil HealthChecker")
	}
}

func TestHealthChecker_NoProbes_ReturnsEmpty(t *testing.T) {
	h := NewHealthChecker()
	results := h.RunAll(context.Background())
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestHealthChecker_HealthyProbe(t *testing.T) {
	h := NewHealthChecker()
	h.Register("vault", func(ctx context.Context) error {
		return nil
	})

	results := h.RunAll(context.Background())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != HealthOK {
		t.Errorf("expected status ok, got %s", results[0].Status)
	}
	if results[0].Name != "vault" {
		t.Errorf("expected name vault, got %s", results[0].Name)
	}
}

func TestHealthChecker_UnhealthyProbe(t *testing.T) {
	h := NewHealthChecker()
	h.Register("db", func(ctx context.Context) error {
		return errors.New("connection refused")
	})

	results := h.RunAll(context.Background())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != HealthDown {
		t.Errorf("expected status down, got %s", results[0].Status)
	}
	if results[0].Message != "connection refused" {
		t.Errorf("unexpected message: %s", results[0].Message)
	}
}

func TestHealthChecker_MultipleProbes(t *testing.T) {
	h := NewHealthChecker()
	h.Register("ok", func(ctx context.Context) error { return nil })
	h.Register("fail", func(ctx context.Context) error { return errors.New("down") })

	results := h.RunAll(context.Background())
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestHealthChecker_CheckedAt_IsRecent(t *testing.T) {
	h := NewHealthChecker()
	h.Register("ts", func(ctx context.Context) error { return nil })

	before := time.Now()
	results := h.RunAll(context.Background())
	after := time.Now()

	if results[0].CheckedAt.Before(before) || results[0].CheckedAt.After(after) {
		t.Error("CheckedAt timestamp out of expected range")
	}
}

func TestOverall_EmptyResults(t *testing.T) {
	if Overall(nil) != HealthOK {
		t.Error("expected ok for empty results")
	}
}

func TestOverall_AllHealthy(t *testing.T) {
	results := []HealthResult{
		{Status: HealthOK},
		{Status: HealthOK},
	}
	if Overall(results) != HealthOK {
		t.Error("expected ok when all healthy")
	}
}

func TestOverall_OneDown(t *testing.T) {
	results := []HealthResult{
		{Status: HealthOK},
		{Status: HealthDown},
	}
	if Overall(results) != HealthDown {
		t.Error("expected down when one probe is down")
	}
}

func TestSummary_Format(t *testing.T) {
	results := []HealthResult{
		{Status: HealthOK},
	}
	s := Summary(results)
	expected := "health: ok (1 probes)"
	if s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}
