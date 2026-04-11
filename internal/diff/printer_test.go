package diff

import (
	"bytes"
	"strings"
	"testing"
)

func newTestPrinter() (*Printer, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	p := NewPrinter(buf)
	p.color = false // disable color for deterministic output
	return p, buf
}

func TestPrint_NoChanges(t *testing.T) {
	p, buf := newTestPrinter()
	p.Print(nil)
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected 'No changes', got %q", buf.String())
	}
}

func TestPrint_Added(t *testing.T) {
	p, buf := newTestPrinter()
	p.Print([]Change{{Type: Added, Key: "FOO", NewValue: "bar"}})
	out := buf.String()
	if !strings.Contains(out, "+ FOO=bar") {
		t.Errorf("expected added line, got %q", out)
	}
}

func TestPrint_Removed(t *testing.T) {
	p, buf := newTestPrinter()
	p.Print([]Change{{Type: Removed, Key: "OLD", OldValue: "val"}})
	out := buf.String()
	if !strings.Contains(out, "- OLD=val") {
		t.Errorf("expected removed line, got %q", out)
	}
}

func TestPrint_Updated(t *testing.T) {
	p, buf := newTestPrinter()
	p.Print([]Change{{Type: Updated, Key: "KEY", OldValue: "old", NewValue: "new"}})
	out := buf.String()
	if !strings.Contains(out, "~ KEY: old → new") {
		t.Errorf("expected updated line, got %q", out)
	}
}

func TestPrint_Unchanged(t *testing.T) {
	p, buf := newTestPrinter()
	p.Print([]Change{{Type: Unchanged, Key: "SAME", OldValue: "x", NewValue: "x"}})
	out := buf.String()
	if !strings.Contains(out, "  SAME=x") {
		t.Errorf("expected unchanged line, got %q", out)
	}
}

func TestPrint_SortedByKey(t *testing.T) {
	p, buf := newTestPrinter()
	p.Print([]Change{
		{Type: Added, Key: "ZZZ", NewValue: "1"},
		{Type: Added, Key: "AAA", NewValue: "2"},
	})
	out := buf.String()
	idxAAA := strings.Index(out, "AAA")
	idxZZZ := strings.Index(out, "ZZZ")
	if idxAAA > idxZZZ {
		t.Errorf("expected AAA before ZZZ in output")
	}
}

func TestSummary_Counts(t *testing.T) {
	p, buf := newTestPrinter()
	p.Summary([]Change{
		{Type: Added},
		{Type: Added},
		{Type: Removed},
		{Type: Updated},
		{Type: Unchanged},
	})
	out := buf.String()
	if !strings.Contains(out, "+2 added") {
		t.Errorf("expected +2 added in summary, got %q", out)
	}
	if !strings.Contains(out, "-1 removed") {
		t.Errorf("expected -1 removed in summary, got %q", out)
	}
	if !strings.Contains(out, "~1 updated") {
		t.Errorf("expected ~1 updated in summary, got %q", out)
	}
	if !strings.Contains(out, "1 unchanged") {
		t.Errorf("expected 1 unchanged in summary, got %q", out)
	}
}
