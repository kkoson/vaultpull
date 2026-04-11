package sync

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestProgress_ProfileStarted_Verbose(t *testing.T) {
	var buf bytes.Buffer
	p := NewProgress(&buf, 3, true)
	p.ProfileStarted("production")
	if !strings.Contains(buf.String(), "production") {
		t.Errorf("expected profile name in output, got: %s", buf.String())
	}
}

func TestProgress_ProfileStarted_Silent(t *testing.T) {
	var buf bytes.Buffer
	p := NewProgress(&buf, 3, false)
	p.ProfileStarted("production")
	if buf.Len() != 0 {
		t.Errorf("expected no output in non-verbose mode, got: %s", buf.String())
	}
}

func TestProgress_ProfileDone_IncrementsCount(t *testing.T) {
	var buf bytes.Buffer
	p := NewProgress(&buf, 2, true)
	p.ProfileDone("staging", 3, 1, 0)
	if p.FailedCount() != 0 {
		t.Errorf("expected 0 failures, got %d", p.FailedCount())
	}
	if !strings.Contains(buf.String(), "+3") {
		t.Errorf("expected added count in output, got: %s", buf.String())
	}
}

func TestProgress_ProfileFailed_IncrementsFailure(t *testing.T) {
	var buf bytes.Buffer
	p := NewProgress(&buf, 2, false)
	p.ProfileFailed("dev", errors.New("vault unreachable"))
	if p.FailedCount() != 1 {
		t.Errorf("expected 1 failure, got %d", p.FailedCount())
	}
	if !strings.Contains(buf.String(), "vault unreachable") {
		t.Errorf("expected error message in output, got: %s", buf.String())
	}
}

func TestProgress_Summary_AllSuccess(t *testing.T) {
	var buf bytes.Buffer
	p := NewProgress(&buf, 2, false)
	p.ProfileDone("a", 1, 0, 0)
	p.ProfileDone("b", 2, 0, 0)
	p.Summary()
	if !strings.Contains(buf.String(), "2/2 profiles succeeded") {
		t.Errorf("unexpected summary: %s", buf.String())
	}
}

func TestProgress_Summary_PartialFailure(t *testing.T) {
	var buf bytes.Buffer
	p := NewProgress(&buf, 2, false)
	p.ProfileDone("a", 1, 0, 0)
	p.ProfileFailed("b", errors.New("timeout"))
	p.Summary()
	out := buf.String()
	if !strings.Contains(out, "1/2 profiles succeeded") {
		t.Errorf("unexpected summary: %s", out)
	}
	if !strings.Contains(out, "1 failed") {
		t.Errorf("expected failure count in summary: %s", out)
	}
}

func TestNewProgress_NilWriterUsesStdout(t *testing.T) {
	p := NewProgress(nil, 1, false)
	if p.out == nil {
		t.Error("expected non-nil writer when nil is passed")
	}
}
