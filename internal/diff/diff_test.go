package diff_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/diff"
)

func TestCompare_AllAdded(t *testing.T) {
	existing := map[string]string{}
	incoming := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := diff.Compare(existing, incoming)
	added, updated, removed, unchanged := result.Summary()

	if added != 2 || updated != 0 || removed != 0 || unchanged != 0 {
		t.Errorf("expected 2 added, got added=%d updated=%d removed=%d unchanged=%d", added, updated, removed, unchanged)
	}
	if !result.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestCompare_AllUnchanged(t *testing.T) {
	existing := map[string]string{"FOO": "bar"}
	incoming := map[string]string{"FOO": "bar"}

	result := diff.Compare(existing, incoming)
	added, updated, removed, unchanged := result.Summary()

	if added != 0 || updated != 0 || removed != 0 || unchanged != 1 {
		t.Errorf("expected 1 unchanged, got added=%d updated=%d removed=%d unchanged=%d", added, updated, removed, unchanged)
	}
	if result.HasChanges() {
		t.Error("expected HasChanges to be false")
	}
}

func TestCompare_Updated(t *testing.T) {
	existing := map[string]string{"FOO": "old"}
	incoming := map[string]string{"FOO": "new"}

	result := diff.Compare(existing, incoming)
	added, updated, removed, _ := result.Summary()

	if added != 0 || updated != 1 || removed != 0 {
		t.Errorf("expected 1 updated, got added=%d updated=%d removed=%d", added, updated, removed)
	}

	for _, c := range result.Changes {
		if c.Key == "FOO" {
			if c.OldValue != "old" || c.NewValue != "new" {
				t.Errorf("unexpected old/new values: %q / %q", c.OldValue, c.NewValue)
			}
		}
	}
}

func TestCompare_Removed(t *testing.T) {
	existing := map[string]string{"FOO": "bar", "GONE": "bye"}
	incoming := map[string]string{"FOO": "bar"}

	result := diff.Compare(existing, incoming)
	added, updated, removed, unchanged := result.Summary()

	if added != 0 || updated != 0 || removed != 1 || unchanged != 1 {
		t.Errorf("expected 1 removed 1 unchanged, got added=%d updated=%d removed=%d unchanged=%d", added, updated, removed, unchanged)
	}
}

func TestCompare_Mixed(t *testing.T) {
	existing := map[string]string{"A": "1", "B": "old", "C": "keep"}
	incoming := map[string]string{"B": "new", "C": "keep", "D": "fresh"}

	result := diff.Compare(existing, incoming)
	added, updated, removed, unchanged := result.Summary()

	if added != 1 || updated != 1 || removed != 1 || unchanged != 1 {
		t.Errorf("mixed diff mismatch: added=%d updated=%d removed=%d unchanged=%d", added, updated, removed, unchanged)
	}
	if !result.HasChanges() {
		t.Error("expected HasChanges true")
	}
}
