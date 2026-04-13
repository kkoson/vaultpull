package sync

import (
	"bytes"
	"errors"
	"testing"
)

func TestNewShadowWriter_DefaultsToStderr(t *testing.T) {
	sw := NewShadowWriter(ShadowOn, nil)
	if sw == nil {
		t.Fatal("expected non-nil ShadowWriter")
	}
	if sw.dest == nil {
		t.Fatal("expected dest to default to stderr")
	}
}

func TestShadowWriter_Off_DoesNotRecord(t *testing.T) {
	buf := &bytes.Buffer{}
	sw := NewShadowWriter(ShadowOff, buf)
	sw.Write("prod", func() error { return errors.New("boom") })
	if len(sw.Records()) != 0 {
		t.Errorf("expected 0 records, got %d", len(sw.Records()))
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got %q", buf.String())
	}
}

func TestShadowWriter_On_RecordsSuccess(t *testing.T) {
	buf := &bytes.Buffer{}
	sw := NewShadowWriter(ShadowOn, buf)
	sw.Write("staging", func() error { return nil })
	recs := sw.Records()
	if len(recs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(recs))
	}
	if recs[0].Profile != "staging" {
		t.Errorf("expected profile=staging, got %q", recs[0].Profile)
	}
	if recs[0].Err != nil {
		t.Errorf("expected no error, got %v", recs[0].Err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no stderr output on success, got %q", buf.String())
	}
}

func TestShadowWriter_On_RecordsError(t *testing.T) {
	buf := &bytes.Buffer{}
	sw := NewShadowWriter(ShadowOn, buf)
	sw.Write("dev", func() error { return errors.New("vault unreachable") })
	recs := sw.Records()
	if len(recs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(recs))
	}
	if recs[0].Err == nil {
		t.Error("expected error to be recorded")
	}
	if buf.Len() == 0 {
		t.Error("expected error message written to dest")
	}
}

func TestShadowWriter_ErrorCount(t *testing.T) {
	buf := &bytes.Buffer{}
	sw := NewShadowWriter(ShadowOn, buf)
	sw.Write("a", func() error { return nil })
	sw.Write("b", func() error { return errors.New("fail") })
	sw.Write("c", func() error { return errors.New("fail") })
	if sw.ErrorCount() != 2 {
		t.Errorf("expected 2 errors, got %d", sw.ErrorCount())
	}
}

func TestShadowWriter_NilFn_NoError(t *testing.T) {
	sw := NewShadowWriter(ShadowOn, &bytes.Buffer{})
	sw.Write("x", nil)
	recs := sw.Records()
	if len(recs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(recs))
	}
	if recs[0].Err != nil {
		t.Errorf("expected nil error for nil fn, got %v", recs[0].Err)
	}
}

func TestShadowWriter_Records_IsCopy(t *testing.T) {
	sw := NewShadowWriter(ShadowOn, &bytes.Buffer{})
	sw.Write("p", nil)
	r1 := sw.Records()
	r1[0].Profile = "mutated"
	r2 := sw.Records()
	if r2[0].Profile == "mutated" {
		t.Error("Records() should return a copy, not a reference")
	}
}
