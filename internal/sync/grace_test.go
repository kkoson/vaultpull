package sync

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/config"
)

func TestDefaultGraceConfig_Values(t *testing.T) {
	cfg := DefaultGraceConfig()
	if cfg.Period != 5*time.Second {
		t.Fatalf("expected 5s, got %v", cfg.Period)
	}
}

func TestNewGraceManager_ZeroPeriodUsesDefault(t *testing.T) {
	gm := NewGraceManager(GraceConfig{})
	if gm.cfg.Period != 5*time.Second {
		t.Fatalf("expected default period, got %v", gm.cfg.Period)
	}
}

func TestGraceManager_NoFailure_NotInGrace(t *testing.T) {
	gm := NewGraceManager(DefaultGraceConfig())
	if gm.InGrace("prod") {
		t.Fatal("expected not in grace before any failure")
	}
}

func TestGraceManager_RecordAndInGrace(t *testing.T) {
	gm := NewGraceManager(GraceConfig{Period: time.Minute})
	gm.RecordFailure("prod")
	if !gm.InGrace("prod") {
		t.Fatal("expected in grace immediately after failure")
	}
}

func TestGraceManager_ExpiredGrace(t *testing.T) {
	gm := NewGraceManager(GraceConfig{Period: time.Millisecond})
	gm.RecordFailure("prod")
	time.Sleep(5 * time.Millisecond)
	if gm.InGrace("prod") {
		t.Fatal("expected grace period to have expired")
	}
}

func TestGraceManager_Reset_ClearsRecord(t *testing.T) {
	gm := NewGraceManager(GraceConfig{Period: time.Minute})
	gm.RecordFailure("prod")
	gm.Reset("prod")
	if gm.InGrace("prod") {
		t.Fatal("expected not in grace after reset")
	}
}

func TestWithGrace_PanicsOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on nil GraceManager")
		}
	}()
	WithGrace(nil, func(_ context.Context, _ config.Profile) error { return nil })
}

func TestWithGrace_SuccessResetsGrace(t *testing.T) {
	gm := NewGraceManager(GraceConfig{Period: time.Minute})
	p := config.Profile{Name: "dev"}
	gm.RecordFailure(p.Name)

	fn := WithGrace(gm, func(_ context.Context, _ config.Profile) error { return nil })
	if err := fn(context.Background(), p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gm.InGrace(p.Name) {
		t.Fatal("expected grace to be reset after success")
	}
}

func TestWithGrace_FailureWithinGrace_Suppressed(t *testing.T) {
	gm := NewGraceManager(GraceConfig{Period: time.Minute})
	p := config.Profile{Name: "dev"}
	boom := errors.New("vault unreachable")

	fn := WithGrace(gm, func(_ context.Context, _ config.Profile) error { return boom })
	if err := fn(context.Background(), p); err != nil {
		t.Fatalf("expected error suppressed within grace, got: %v", err)
	}
}

func TestWithGrace_FailureAfterGrace_ReturnsError(t *testing.T) {
	gm := NewGraceManager(GraceConfig{Period: time.Millisecond})
	p := config.Profile{Name: "dev"}
	boom := errors.New("vault unreachable")

	gm.RecordFailure(p.Name)
	time.Sleep(5 * time.Millisecond)

	fn := WithGrace(gm, func(_ context.Context, _ config.Profile) error { return boom })
	if err := fn(context.Background(), p); err == nil {
		t.Fatal("expected error after grace period expired")
	}
}
