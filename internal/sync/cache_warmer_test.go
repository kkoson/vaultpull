package sync

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/token"
)

func buildWarmer(t *testing.T, profiles []config.Profile) (*CacheWarmer, string) {
	t.Helper()
	dir := t.TempDir()
	cachePath := filepath.Join(dir, "token_cache.json")
	cache := token.NewCache(cachePath)
	cfg := &config.Config{Profiles: profiles}
	return NewCacheWarmer(cfg, cache, 2*time.Second), cachePath
}

func TestCacheWarmer_NoProfiles(t *testing.T) {
	w, _ := buildWarmer(t, nil)
	results := w.WarmAll(context.Background())
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestCacheWarmer_CacheMiss_NoError(t *testing.T) {
	w, _ := buildWarmer(t, []config.Profile{{Name: "dev"}})
	results := w.WarmAll(context.Background())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Profile != "dev" {
		t.Errorf("expected profile 'dev', got %q", r.Profile)
	}
	if r.Hit {
		t.Error("expected cache miss (Hit=false)")
	}
	if r.Err != nil {
		t.Errorf("unexpected error on miss: %v", r.Err)
	}
}

func TestCacheWarmer_CacheHit(t *testing.T) {
	dir := t.TempDir()
	cachePath := filepath.Join(dir, "token_cache.json")
	cache := token.NewCache(cachePath)

	// Pre-populate cache with a valid entry.
	_ = cache.Save("prod", "s.testtoken", time.Now().Add(1*time.Hour))

	cfg := &config.Config{Profiles: []config.Profile{{Name: "prod"}}}
	w := NewCacheWarmer(cfg, cache, 2*time.Second)

	results := w.WarmAll(context.Background())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Hit {
		t.Error("expected cache hit")
	}
	if results[0].Err != nil {
		t.Errorf("unexpected error: %v", results[0].Err)
	}
}

func TestCacheWarmer_ContextCancelled(t *testing.T) {
	w, _ := buildWarmer(t, []config.Profile{{Name: "staging"}})
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	results := w.WarmAll(ctx)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	// With an already-cancelled context the warmer may time out or succeed from
	// cache; either way we should get a result without a panic.
	_ = results[0]
}

func TestNewCacheWarmer_DefaultTimeout(t *testing.T) {
	dir := t.TempDir()
	cache := token.NewCache(filepath.Join(dir, "c.json"))
	cfg := &config.Config{}
	w := NewCacheWarmer(cfg, cache, 0)
	if w.timeout != 10*time.Second {
		t.Errorf("expected default 10s timeout, got %v", w.timeout)
	}
	_ = os.Remove(filepath.Join(dir, "c.json"))
}
