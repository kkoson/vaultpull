package sync

import (
	"testing"
	"time"
)

func TestDefaultCachePolicyConfig_Values(t *testing.T) {
	cfg := DefaultCachePolicyConfig()
	if cfg.Mode != CachePolicyReadThrough {
		t.Fatalf("expected ReadThrough, got %d", cfg.Mode)
	}
	if cfg.TTL != 5*time.Minute {
		t.Fatalf("expected 5m TTL, got %v", cfg.TTL)
	}
}

func TestNewCachePolicy_ZeroTTLUsesDefault(t *testing.T) {
	p := NewCachePolicy(CachePolicyConfig{Mode: CachePolicyReadThrough})
	if p.cfg.TTL != 5*time.Minute {
		t.Fatalf("expected default TTL, got %v", p.cfg.TTL)
	}
}

func TestCachePolicy_IsFresh_Miss(t *testing.T) {
	p := NewCachePolicy(DefaultCachePolicyConfig())
	if p.IsFresh("dev") {
		t.Fatal("expected not fresh on empty store")
	}
}

func TestCachePolicy_IsFresh_Hit(t *testing.T) {
	p := NewCachePolicy(DefaultCachePolicyConfig())
	p.Set("dev", map[string]string{"KEY": "val"})
	if !p.IsFresh("dev") {
		t.Fatal("expected fresh after set")
	}
}

func TestCachePolicy_IsFresh_Bypass(t *testing.T) {
	p := NewCachePolicy(CachePolicyConfig{Mode: CachePolicyBypass, TTL: time.Hour})
	p.Set("dev", map[string]string{"KEY": "val"})
	if p.IsFresh("dev") {
		t.Fatal("bypass mode must never report fresh")
	}
}

func TestCachePolicy_Get_ReturnsNilOnMiss(t *testing.T) {
	p := NewCachePolicy(DefaultCachePolicyConfig())
	if got := p.Get("missing"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestCachePolicy_Get_ReturnsCopy(t *testing.T) {
	p := NewCachePolicy(DefaultCachePolicyConfig())
	orig := map[string]string{"A": "1"}
	p.Set("dev", orig)
	got := p.Get("dev")
	got["A"] = "mutated"
	if p.Get("dev")["A"] != "1" {
		t.Fatal("Get must return an independent copy")
	}
}

func TestCachePolicy_Invalidate(t *testing.T) {
	p := NewCachePolicy(DefaultCachePolicyConfig())
	p.Set("dev", map[string]string{"K": "v"})
	p.Invalidate("dev")
	if p.IsFresh("dev") {
		t.Fatal("expected not fresh after invalidation")
	}
}

func TestCachePolicy_Expired_NotFresh(t *testing.T) {
	cfg := CachePolicyConfig{Mode: CachePolicyReadThrough, TTL: time.Millisecond}
	p := NewCachePolicy(cfg)
	p.Set("dev", map[string]string{"K": "v"})
	time.Sleep(5 * time.Millisecond)
	if p.IsFresh("dev") {
		t.Fatal("expected stale after TTL expiry")
	}
}
