package sync

import (
	"testing"
)

func TestNewFingerprintStore_NotNil(t *testing.T) {
	fs := NewFingerprintStore()
	if fs == nil {
		t.Fatal("expected non-nil FingerprintStore")
	}
}

func TestFingerprintStore_Compute_Empty(t *testing.T) {
	fs := NewFingerprintStore()
	fp := fs.Compute(map[string]string{})
	if fp == "" {
		t.Fatal("expected non-empty fingerprint for empty map")
	}
}

func TestFingerprintStore_Compute_OrderIndependent(t *testing.T) {
	fs := NewFingerprintStore()
	a := fs.Compute(map[string]string{"FOO": "1", "BAR": "2"})
	b := fs.Compute(map[string]string{"BAR": "2", "FOO": "1"})
	if a != b {
		t.Fatalf("expected same fingerprint regardless of key order, got %q vs %q", a, b)
	}
}

func TestFingerprintStore_Changed_TrueOnFirstCall(t *testing.T) {
	fs := NewFingerprintStore()
	secrets := map[string]string{"KEY": "value"}
	if !fs.Changed("dev", secrets) {
		t.Fatal("expected Changed=true when no fingerprint recorded yet")
	}
}

func TestFingerprintStore_Changed_FalseAfterRecord(t *testing.T) {
	fs := NewFingerprintStore()
	secrets := map[string]string{"KEY": "value"}
	fs.Record("dev", secrets)
	if fs.Changed("dev", secrets) {
		t.Fatal("expected Changed=false after recording same secrets")
	}
}

func TestFingerprintStore_Changed_TrueAfterSecretUpdated(t *testing.T) {
	fs := NewFingerprintStore()
	fs.Record("dev", map[string]string{"KEY": "old"})
	if !fs.Changed("dev", map[string]string{"KEY": "new"}) {
		t.Fatal("expected Changed=true after secret value updated")
	}
}

func TestFingerprintStore_Clear_ReturnsToUnrecorded(t *testing.T) {
	fs := NewFingerprintStore()
	secrets := map[string]string{"KEY": "value"}
	fs.Record("dev", secrets)
	fs.Clear("dev")
	if !fs.Changed("dev", secrets) {
		t.Fatal("expected Changed=true after Clear")
	}
}

func TestFingerprintStore_MultipleProfiles_Independent(t *testing.T) {
	fs := NewFingerprintStore()
	fs.Record("dev", map[string]string{"A": "1"})
	fs.Record("prod", map[string]string{"A": "2"})

	if fs.Changed("dev", map[string]string{"A": "1"}) {
		t.Fatal("dev should not have changed")
	}
	if !fs.Changed("prod", map[string]string{"A": "1"}) {
		t.Fatal("prod should have changed")
	}
}
