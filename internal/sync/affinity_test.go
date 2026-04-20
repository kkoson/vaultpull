package sync

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
)

func TestDefaultAffinityConfig_Values(t *testing.T) {
	cfg := DefaultAffinityConfig()
	if cfg.TagKey != "affinity" {
		t.Fatalf("expected tag key 'affinity', got %q", cfg.TagKey)
	}
	if cfg.SlotCount != 16 {
		t.Fatalf("expected slot count 16, got %d", cfg.SlotCount)
	}
}

func TestNewAffinityRouter_ZeroConfigUsesDefaults(t *testing.T) {
	r := NewAffinityRouter(AffinityConfig{})
	if r.cfg.TagKey != "affinity" {
		t.Fatalf("expected default tag key, got %q", r.cfg.TagKey)
	}
	if r.cfg.SlotCount != 16 {
		t.Fatalf("expected default slot count, got %d", r.cfg.SlotCount)
	}
}

func TestAffinityRouter_SameProfileSameSlot(t *testing.T) {
	r := NewAffinityRouter(DefaultAffinityConfig())
	p := config.Profile{Name: "prod"}
	s1 := r.Assign(p)
	s2 := r.Assign(p)
	if s1 != s2 {
		t.Fatalf("expected stable slot, got %d then %d", s1, s2)
	}
}

func TestAffinityRouter_SharedTagSharesSlot(t *testing.T) {
	r := NewAffinityRouter(DefaultAffinityConfig())
	p1 := config.Profile{Name: "prod-us", Tags: []string{"affinity:us-east"}}
	p2 := config.Profile{Name: "prod-eu", Tags: []string{"affinity:us-east"}}
	if r.Assign(p1) != r.Assign(p2) {
		t.Fatal("profiles with same affinity tag should receive the same slot")
	}
}

func TestAffinityRouter_DifferentTagDifferentSlot(t *testing.T) {
	r := NewAffinityRouter(AffinityConfig{TagKey: "affinity", SlotCount: 1024})
	p1 := config.Profile{Name: "a", Tags: []string{"affinity:alpha"}}
	p2 := config.Profile{Name: "b", Tags: []string{"affinity:beta"}}
	// With 1024 slots a collision is astronomically unlikely for two short strings.
	if r.Assign(p1) == r.Assign(p2) {
		t.Fatal("profiles with distinct affinity tags should not share a slot")
	}
}

func TestAffinityRouter_NoTag_FallsBackToName(t *testing.T) {
	r := NewAffinityRouter(DefaultAffinityConfig())
	p := config.Profile{Name: "staging"}
	slot := r.Assign(p)
	if slot >= r.cfg.SlotCount {
		t.Fatalf("slot %d out of range [0, %d)", slot, r.cfg.SlotCount)
	}
}

func TestAffinityRouter_SlotInRange(t *testing.T) {
	r := NewAffinityRouter(DefaultAffinityConfig())
	profiles := []config.Profile{
		{Name: "dev"}, {Name: "staging"}, {Name: "prod"},
	}
	for _, p := range profiles {
		s := r.Assign(p)
		if s >= r.cfg.SlotCount {
			t.Fatalf("profile %q: slot %d >= SlotCount %d", p.Name, s, r.cfg.SlotCount)
		}
	}
}
