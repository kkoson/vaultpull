package sync

import (
	"errors"
	"testing"
)

func TestWithCachePolicy_PanicsOnNilPolicy(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil policy")
		}
	}()
	WithCachePolicy(nil, nil, nil)
}

func TestWithCachePolicy_FetchesWhenNotFresh(t *testing.T) {
	policy := NewCachePolicy(DefaultCachePolicyConfig())
	fetched := false
	fetch := func(p string) (map[string]string, error) {
		fetched = true
		return map[string]string{"KEY": "val"}, nil
	}
	consumed := map[string]string{}
	consume := func(p string, s map[string]string) error {
		for k, v := range s {
			consumed[k] = v
		}
		return nil
	}
	fn := WithCachePolicy(policy, fetch, consume)
	if err := fn("dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fetched {
		t.Fatal("expected fetch to be called")
	}
	if consumed["KEY"] != "val" {
		t.Fatalf("expected val, got %q", consumed["KEY"])
	}
}

func TestWithCachePolicy_ServesFromCacheWhenFresh(t *testing.T) {
	policy := NewCachePolicy(DefaultCachePolicyConfig())
	policy.Set("dev", map[string]string{"KEY": "cached"})
	fetchCalled := false
	fetch := func(p string) (map[string]string, error) {
		fetchCalled = true
		return nil, nil
	}
	consumed := map[string]string{}
	consume := func(p string, s map[string]string) error {
		for k, v := range s {
			consumed[k] = v
		}
		return nil
	}
	fn := WithCachePolicy(policy, fetch, consume)
	if err := fn("dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fetchCalled {
		t.Fatal("expected fetch NOT to be called when cache is fresh")
	}
	if consumed["KEY"] != "cached" {
		t.Fatalf("expected cached, got %q", consumed["KEY"])
	}
}

func TestWithCachePolicy_FetchError_Propagated(t *testing.T) {
	policy := NewCachePolicy(DefaultCachePolicyConfig())
	fetchErr := errors.New("vault unavailable")
	fetch := func(p string) (map[string]string, error) { return nil, fetchErr }
	consume := func(p string, s map[string]string) error { return nil }
	fn := WithCachePolicy(policy, fetch, consume)
	if err := fn("dev"); !errors.Is(err, fetchErr) {
		t.Fatalf("expected fetchErr, got %v", err)
	}
}

func TestWithCachePolicy_WritesBackInWriteThrough(t *testing.T) {
	cfg := CachePolicyConfig{Mode: CachePolicyWriteThrough, TTL: DefaultCachePolicyConfig().TTL}
	policy := NewCachePolicy(cfg)
	fetch := func(p string) (map[string]string, error) {
		return map[string]string{"X": "y"}, nil
	}
	consume := func(p string, s map[string]string) error { return nil }
	fn := WithCachePolicy(policy, fetch, consume)
	if err := fn("prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !policy.IsFresh("prod") {
		t.Fatal("expected cache to be populated after write-through fetch")
	}
}
