package rotate

import (
	"testing"
	"time"
)

func TestDefaultSecretIDTTL(t *testing.T) {
	if defaultSecretIDTTL != 24*time.Hour {
		t.Errorf("expected defaultSecretIDTTL to be 24h, got %s", defaultSecretIDTTL)
	}
}

func TestSecretIDResponse_Fields(t *testing.T) {
	resp := SecretIDResponse{
		SecretID: "abc-123",
		TTL:      12 * time.Hour,
	}
	if resp.SecretID != "abc-123" {
		t.Errorf("unexpected SecretID: %s", resp.SecretID)
	}
	if resp.TTL != 12*time.Hour {
		t.Errorf("unexpected TTL: %s", resp.TTL)
	}
}

func TestTTLParsing_ZeroFallsBackToDefault(t *testing.T) {
	// Simulate the TTL parsing branch: ttlSec == 0 should keep defaultSecretIDTTL.
	ttl := defaultSecretIDTTL
	ttlSec := float64(0)
	if ttlSec > 0 {
		ttl = time.Duration(ttlSec) * time.Second
	}
	if ttl != defaultSecretIDTTL {
		t.Errorf("expected fallback to default TTL, got %s", ttl)
	}
}

func TestTTLParsing_PositiveOverridesDefault(t *testing.T) {
	ttl := defaultSecretIDTTL
	ttlSec := float64(3600)
	if ttlSec > 0 {
		ttl = time.Duration(ttlSec) * time.Second
	}
	if ttl != time.Hour {
		t.Errorf("expected 1h TTL, got %s", ttl)
	}
}
