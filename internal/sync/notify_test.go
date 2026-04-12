package sync

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNotify_NoneLevel_NoOutput(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf, NotifyNone)
	n.Notify(NotifyEvent{Profile: "dev", Success: false, Err: errors.New("boom")})
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got %q", buf.String())
	}
}

func TestNotify_FailureLevel_SkipsSuccess(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf, NotifyFailure)
	n.Notify(NotifyEvent{Profile: "dev", Success: true, Changes: 2})
	if buf.Len() != 0 {
		t.Fatalf("expected no output for success, got %q", buf.String())
	}
}

func TestNotify_FailureLevel_PrintsFailure(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf, NotifyFailure)
	n.Notify(NotifyEvent{
		Profile:   "staging",
		Success:   false,
		Err:       errors.New("vault unreachable"),
		Duration:  50 * time.Millisecond,
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	out := buf.String()
	if !strings.Contains(out, "staging") {
		t.Errorf("expected profile name in output, got %q", out)
	}
	if !strings.Contains(out, "vault unreachable") {
		t.Errorf("expected error message in output, got %q", out)
	}
	if !strings.Contains(out, "✘") {
		t.Errorf("expected failure icon in output, got %q", out)
	}
}

func TestNotify_AllLevel_PrintsSuccess(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf, NotifyAll)
	n.Notify(NotifyEvent{
		Profile:   "production",
		Success:   true,
		Changes:   5,
		Duration:  120 * time.Millisecond,
		Timestamp: time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC),
	})
	out := buf.String()
	if !strings.Contains(out, "production") {
		t.Errorf("expected profile name, got %q", out)
	}
	if !strings.Contains(out, "5 change(s)") {
		t.Errorf("expected change count, got %q", out)
	}
	if !strings.Contains(out, "✔") {
		t.Errorf("expected success icon, got %q", out)
	}
}

func TestNotify_NilWriter_UsesStderr(t *testing.T) {
	// Should not panic when out is nil.
	n := NewNotifier(nil, NotifyNone)
	if n.out == nil {
		t.Fatal("expected fallback writer, got nil")
	}
}

func TestNotify_ZeroTimestamp_Filled(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf, NotifyAll)
	n.Notify(NotifyEvent{Profile: "dev", Success: true})
	if buf.Len() == 0 {
		t.Fatal("expected output even with zero timestamp")
	}
}
