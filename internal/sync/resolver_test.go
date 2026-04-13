package sync

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
)

func buildResolverConfig(names ...string) *config.Config {
	profiles := make([]config.Profile, 0, len(names))
	for _, n := range names {
		profiles = append(profiles, config.Profile{
			Name:       n,
			VaultPath:  "secret/" + n,
			OutputFile: n + ".env",
		})
	}
	return &config.Config{Profiles: profiles}
}

func TestNewResolver_NotNil(t *testing.T) {
	r := NewResolver(buildResolverConfig("dev"), ResolveModeExact)
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}

func TestResolver_Exact_Match(t *testing.T) {
	r := NewResolver(buildResolverConfig("dev", "prod"), ResolveModeExact)
	got, err := r.Resolve("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Name != "dev" {
		t.Fatalf("expected [dev], got %v", got)
	}
}

func TestResolver_Exact_NoMatch(t *testing.T) {
	r := NewResolver(buildResolverConfig("dev"), ResolveModeExact)
	_, err := r.Resolve("staging")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestResolver_Prefix_MultipleMatch(t *testing.T) {
	r := NewResolver(buildResolverConfig("prod-us", "prod-eu", "dev"), ResolveModePrefix)
	got, err := r.Resolve("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(got))
	}
}

func TestResolver_Prefix_NoMatch(t *testing.T) {
	r := NewResolver(buildResolverConfig("dev"), ResolveModePrefix)
	_, err := r.Resolve("prod")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestResolver_NilConfig_ReturnsError(t *testing.T) {
	r := NewResolver(nil, ResolveModeExact)
	_, err := r.Resolve("dev")
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestResolver_ResolveAll_Success(t *testing.T) {
	r := NewResolver(buildResolverConfig("dev", "prod"), ResolveModeExact)
	got, err := r.ResolveAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(got))
	}
}

func TestResolver_ResolveAll_EmptyConfig(t *testing.T) {
	r := NewResolver(&config.Config{}, ResolveModeExact)
	_, err := r.ResolveAll()
	if err == nil {
		t.Fatal("expected error for empty config")
	}
}
