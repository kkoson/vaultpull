package sync

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNewTracer_NotNil(t *testing.T) {
	tr := NewTracer(TraceLevelBasic, nil)
	if tr == nil {
		t.Fatal("expected non-nil tracer")
	}
}

func TestTracer_Off_RecordsNothing(t *testing.T) {
	var buf bytes.Buffer
	tr := NewTracer(TraceLevelOff, &buf)
	tr.Record("prod", "", "start", 0)
	if len(tr.Entries()) != 0 {
		t.Errorf("expected 0 entries, got %d", len(tr.Entries()))
	}
	if buf.Len() != 0 {
		t.Error("expected no output for TraceLevelOff")
	}
}

func TestTracer_Basic_SkipsStageEvents(t *testing.T) {
	var buf bytes.Buffer
	tr := NewTracer(TraceLevelBasic, &buf)
	tr.Record("prod", "validate", "done", time.Millisecond)
	if len(tr.Entries()) != 0 {
		t.Errorf("basic level should skip stage events; got %d entries", len(tr.Entries()))
	}
}

func TestTracer_Basic_RecordsProfileEvent(t *testing.T) {
	var buf bytes.Buffer
	tr := NewTracer(TraceLevelBasic, &buf)
	tr.Record("prod", "", "start", 0)
	if len(tr.Entries()) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(tr.Entries()))
	}
	if tr.Entries()[0].Profile != "prod" {
		t.Errorf("unexpected profile: %s", tr.Entries()[0].Profile)
	}
}

func TestTracer_Full_RecordsStageEvents(t *testing.T) {
	var buf bytes.Buffer
	tr := NewTracer(TraceLevelFull, &buf)
	tr.Record("prod", "validate", "done", 5*time.Millisecond)
	if len(tr.Entries()) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(tr.Entries()))
	}
	e := tr.Entries()[0]
	if e.Stage != "validate" {
		t.Errorf("expected stage 'validate', got %q", e.Stage)
	}
	if e.Duration != 5*time.Millisecond {
		t.Errorf("unexpected duration: %v", e.Duration)
	}
}

func TestTracer_Emit_ContainsProfileAndEvent(t *testing.T) {
	var buf bytes.Buffer
	tr := NewTracer(TraceLevelBasic, &buf)
	tr.Record("staging", "", "end", 10*time.Millisecond)
	out := buf.String()
	if !strings.Contains(out, "staging") {
		t.Errorf("output missing profile name: %q", out)
	}
	if !strings.Contains(out, "end") {
		t.Errorf("output missing event: %q", out)
	}
}

func TestTracer_NilSafe(t *testing.T) {
	var tr *Tracer
	// should not panic
	tr.Record("prod", "", "start", 0)
}

func TestTracer_Entries_ReturnsCopy(t *testing.T) {
	var buf bytes.Buffer
	tr := NewTracer(TraceLevelFull, &buf)
	tr.Record("a", "", "start", 0)
	e1 := tr.Entries()
	e1[0].Profile = "mutated"
	e2 := tr.Entries()
	if e2[0].Profile == "mutated" {
		t.Error("Entries should return an independent copy")
	}
}
