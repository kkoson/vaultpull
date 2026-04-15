package sync

import (
	"testing"
)

func TestDefaultDigestConfig_Values(t *testing.T) {
	cfg := DefaultDigestConfig()
	if cfg.Algorithm != "sha256" {
		t.Fatalf("expected sha256, got %s", cfg.Algorithm)
	}
}

func TestNewDigester_ZeroConfigUsesDefaults(t *testing.T) {
	d := NewDigester(DigestConfig{})
	if d.cfg.Algorithm != "sha256" {
		t.Fatalf("expected default algorithm, got %s", d.cfg.Algorithm)
	}
}

func TestDigester_Compute_EmptyMap(t *testing.T) {
	d := NewDigester(DefaultDigestConfig())
	hex := d.Compute(map[string]string{})
	if hex == "" {
		t.Fatal("expected non-empty digest for empty map")
	}
}

func TestDigester_Compute_Stable(t *testing.T) {
	d := NewDigester(DefaultDigestConfig())
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	h1 := d.Compute(secrets)
	h2 := d.Compute(secrets)
	if h1 != h2 {
		t.Fatalf("digest not stable: %s != %s", h1, h2)
	}
}

func TestDigester_Compute_OrderIndependent(t *testing.T) {
	d := NewDigester(DefaultDigestConfig())
	a := map[string]string{"KEY_A": "1", "KEY_B": "2"}
	b := map[string]string{"KEY_B": "2", "KEY_A": "1"}
	if d.Compute(a) != d.Compute(b) {
		t.Fatal("digest should be order-independent")
	}
}

func TestDigester_Compute_DifferentValues(t *testing.T) {
	d := NewDigester(DefaultDigestConfig())
	a := map[string]string{"KEY": "value1"}
	b := map[string]string{"KEY": "value2"}
	if d.Compute(a) == d.Compute(b) {
		t.Fatal("different values should produce different digests")
	}
}

func TestDigester_Equal_SameMaps(t *testing.T) {
	d := NewDigester(DefaultDigestConfig())
	secrets := map[string]string{"TOKEN": "abc", "SECRET": "xyz"}
	if !d.Equal(secrets, secrets) {
		t.Fatal("identical maps should be equal")
	}
}

func TestDigester_Equal_DifferentMaps(t *testing.T) {
	d := NewDigester(DefaultDigestConfig())
	a := map[string]string{"TOKEN": "abc"}
	b := map[string]string{"TOKEN": "def"}
	if d.Equal(a, b) {
		t.Fatal("maps with different values should not be equal")
	}
}

func TestDigester_Equal_EmptyMaps(t *testing.T) {
	d := NewDigester(DefaultDigestConfig())
	if !d.Equal(map[string]string{}, map[string]string{}) {
		t.Fatal("two empty maps should be equal")
	}
}
