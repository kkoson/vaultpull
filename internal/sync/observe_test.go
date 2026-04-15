package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewObserver_NilWriterUsesStderr(t *testing.T) {
	obs := NewObserver(ObserveSummary, nil)
	if obs == nil {
		t.Fatal("expected non-nil observer")
	}
	if obs.out == nil {
		t.Fatal("expected non-nil writer fallback")
	}
}

func TestObserver_Off_RecordsNothing(t *testing.T) {
	var buf bytes.Buffer
	obs := NewObserver(ObserveOff, &buf)
	obs.Record("prod", "API_KEY", "added")
	if len(obs.Events()) != 0 {
		t.Fatalf("expected 0 events, got %d", len(obs.Events()))
	}
	if buf.Len() != 0 {
		t.Fatal("expected no output for ObserveOff")
	}
}

func TestObserver_Summary_RecordsProfileEvent(t *testing.T) {
	var buf bytes.Buffer
	obs := NewObserver(ObserveSummary, &buf)
	obs.Record("prod", "", "synced")
	events := obs.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Profile != "prod" || events[0].Change != "synced" {
		t.Errorf("unexpected event: %+v", events[0])
	}
	if !strings.Contains(buf.String(), "prod") {
		t.Error("expected output to contain profile name")
	}
}

func TestObserver_Summary_SkipsKeyEvents(t *testing.T) {
	var buf bytes.Buffer
	obs := NewObserver(ObserveSummary, &buf)
	obs.Record("prod", "SECRET", "added")
	// event stored but not printed at summary level
	if len(obs.Events()) != 1 {
		t.Fatalf("expected event stored")
	}
	if strings.Contains(buf.String(), "SECRET") {
		t.Error("summary level should not print key-level events")
	}
}

func TestObserver_Full_PrintsKeyEvents(t *testing.T) {
	var buf bytes.Buffer
	obs := NewObserver(ObserveFull, &buf)
	obs.Record("staging", "DB_PASS", "updated")
	if !strings.Contains(buf.String(), "DB_PASS") {
		t.Error("full level should print key name")
	}
	if !strings.Contains(buf.String(), "updated") {
		t.Error("full level should print change type")
	}
}

func TestObserver_Reset_ClearsEvents(t *testing.T) {
	var buf bytes.Buffer
	obs := NewObserver(ObserveFull, &buf)
	obs.Record("prod", "KEY", "added")
	obs.Record("prod", "OTHER", "removed")
	if len(obs.Events()) != 2 {
		t.Fatal("expected 2 events before reset")
	}
	obs.Reset()
	if len(obs.Events()) != 0 {
		t.Fatal("expected 0 events after reset")
	}
}

func TestObserver_NilReceiver_Safe(t *testing.T) {
	var obs *Observer
	obs.Record("prod", "KEY", "added") // must not panic
	if obs.Events() != nil {
		t.Fatal("nil observer should return nil events")
	}
	obs.Reset() // must not panic
}
