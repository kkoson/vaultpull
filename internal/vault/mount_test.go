package vault

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func newTestVaultClient(t *testing.T, handler http.Handler) *vaultapi.Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	cfg := vaultapi.DefaultConfig()
	cfg.Address = ts.URL
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create vault client: %v", err)
	}
	client.SetToken("test-token")
	return client
}

func TestFindMount_MatchesLongestPrefix(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": {
				"secret/": {"type": "kv", "options": {"version": "1"}},
				"secret/app/": {"type": "kv", "options": {"version": "2"}},
				"sys/": {"type": "system", "options": {}}
			}
		}`))
	})
	client := newTestVaultClient(t, handler)

	mount, err := FindMount(context.Background(), client, "secret/app/myservice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mount.Path != "secret/app" {
		t.Errorf("expected mount path 'secret/app', got %q", mount.Path)
	}
	if mount.Version != "2" {
		t.Errorf("expected version '2', got %q", mount.Version)
	}
}

func TestFindMount_NoMatchReturnsError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data": {"sys/": {"type": "system", "options": {}}}}`))
	})
	client := newTestVaultClient(t, handler)

	_, err := FindMount(context.Background(), client, "secret/missing")
	if err == nil {
		t.Fatal("expected error for unmatched path, got nil")
	}
}

func TestListMounts_FiltersNonKV(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": {
				"kv/": {"type": "kv", "options": {"version": "2"}},
				"pki/": {"type": "pki", "options": {}},
				"transit/": {"type": "transit", "options": {}}
			}
		}`))
	})
	client := newTestVaultClient(t, handler)

	mounts, err := ListMounts(context.Background(), client)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mounts) != 1 {
		t.Errorf("expected 1 KV mount, got %d", len(mounts))
	}
	if mounts[0].Path != "kv" {
		t.Errorf("expected path 'kv', got %q", mounts[0].Path)
	}
}
