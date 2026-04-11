package rotate_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/audit"
	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/rotate"
	"github.com/your-org/vaultpull/internal/token"
)

func tmpCachePath(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "cache-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestRotator_NonApproleProfile(t *testing.T) {
	profile := &config.Profile{
		Name: "dev",
		Auth: config.Auth{Method: "token"},
	}
	cache := token.NewCache(tmpCachePath(t))
	logger := audit.NewLogger(nil)
	r := rotate.New(profile, cache_, err := r.Rotate(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for non-approle profile")
	}
}

func TestRotator_MissingRoleID(t *testing.T) {
	profile := &config.Profile{
		Name: "dev",
		Auth: config.Auth{Method: "approle", RoleID: ""},
	}
	cache := token.NewCache(tmpCachePath(t))
	logger := audit.NewLogger(nil)
	r := rotate.New(profile, cache, logger)

	_, err := r.Rotate(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for missing role_id")
	}
}

func TestRotator_CacheExpiry_AfterRotate(t *testing.T) {
	path := tmpCachePath(t)
	cache := token.NewCache(path)

	// Pre-populate cache with an expired entry.
	past := time.Now().Add(-1 * time.Hour)
	_ = cache.Save("dev", "old-secret", past)

	_, loaded, err := cache.Load("dev")
	if err == nil && loaded != "" {
		t.Fatal("expected expired cache entry to be absent")
	}
}
