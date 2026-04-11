package sync

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDefaultTimeoutConfig_Values(t *testing.T) {
	cfg := DefaultTimeoutConfig()
	if cfg.ProfileTimeout != 30*time.Second {
		t.Errorf("expected ProfileTimeout 30s, got %s", cfg.ProfileTimeout)
	}
	if cfg.GlobalTimeout != 5*time.Minute {
		t.Errorf("expected GlobalTimeout 5m, got %s", cfg.GlobalTimeout)
	}
}

func TestWithProfileTimeout_ZeroDisablesDeadline(t *testing.T) {
	cfg := TimeoutConfig{ProfileTimeout: 0}
	ctx, cancel := WithProfileTimeout(context.Background(), cfg)
	defer cancel()
	if _, ok := ctx.Deadline(); ok {
		t.Error("expected no deadline when ProfileTimeout is zero")
	}
}

func TestWithProfileTimeout_SetsDeadline(t *testing.T) {
	cfg := TimeoutConfig{ProfileTimeout: 10 * time.Second}
	ctx, cancel := WithProfileTimeout(context.Background(), cfg)
	defer cancel()
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected a deadline to be set")
	}
	if time.Until(deadline) > 10*time.Second {
		t.Error("deadline is further away than expected")
	}
}

func TestWithGlobalTimeout_ZeroDisablesDeadline(t *testing.T) {
	cfg := TimeoutConfig{GlobalTimeout: 0}
	ctx, cancel := WithGlobalTimeout(context.Background(), cfg)
	defer cancel()
	if _, ok := ctx.Deadline(); ok {
		t.Error("expected no deadline when GlobalTimeout is zero")
	}
}

func TestWithGlobalTimeout_SetsDeadline(t *testing.T) {
	cfg := TimeoutConfig{GlobalTimeout: 2 * time.Minute}
	ctx, cancel := WithGlobalTimeout(context.Background(), cfg)
	defer cancel()
	_, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected a deadline to be set")
	}
}

func TestTimeoutError_ProfileMessage(t *testing.T) {
	err := &TimeoutError{Profile: "prod", Limit: 30 * time.Second}
	want := `sync timed out for profile "prod" after 30s`
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestTimeoutError_GlobalMessage(t *testing.T) {
	err := &TimeoutError{Limit: 5 * time.Minute}
	want := "global sync timed out after 5m0s"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestIsTimeout_Nil(t *testing.T) {
	if IsTimeout(nil) {
		t.Error("expected false for nil error")
	}
}

func TestIsTimeout_DeadlineExceeded(t *testing.T) {
	if !IsTimeout(context.DeadlineExceeded) {
		t.Error("expected true for context.DeadlineExceeded")
	}
}

func TestIsTimeout_TimeoutError(t *testing.T) {
	err := &TimeoutError{Profile: "dev", Limit: 10 * time.Second}
	if !IsTimeout(err) {
		t.Error("expected true for *TimeoutError")
	}
}

func TestIsTimeout_OtherError(t *testing.T) {
	if IsTimeout(errors.New("some other error")) {
		t.Error("expected false for unrelated error")
	}
}
