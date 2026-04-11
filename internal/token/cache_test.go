package token_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/token"
)

func tmpCachePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "token.json")
}

func TestCache_SaveAndLoad(t *testing.T) {
	c, _ := token.NewCache(tmpCachePath(t))
	e := token.Entry{
		Token:     "s.abc123",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	if err := c.Save(e); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := c.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Token != e.Token {
		t.Errorf("token: got %q, want %q", got.Token, e.Token)
	}
}

func TestCache_Load_CacheMiss_NoFile(t *testing.T) {
	c, _ := token.NewCache(tmpCachePath(t))
	_, err := c.Load()
	if err != token.ErrCacheMiss {
		t.Errorf("expected ErrCacheMiss, got %v", err)
	}
}

func TestCache_Load_ExpiredToken(t *testing.T) {
	c, _ := token.NewCache(tmpCachePath(t))
	e := token.Entry{
		Token:     "s.expired",
		ExpiresAt: time.Now().Add(-time.Minute),
	}
	_ = c.Save(e)
	_, err := c.Load()
	if err != token.ErrCacheMiss {
		t.Errorf("expected ErrCacheMiss for expired token, got %v", err)
	}
}

func TestCache_Clear(t *testing.T) {
	p := tmpCachePath(t)
	c, _ := token.NewCache(p)
	_ = c.Save(token.Entry{Token: "s.x", ExpiresAt: time.Now().Add(time.Hour)})
	if err := c.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Error("expected file to be removed after Clear")
	}
}

func TestCache_Clear_NoFile(t *testing.T) {
	c, _ := token.NewCache(tmpCachePath(t))
	if err := c.Clear(); err != nil {
		t.Errorf("Clear on missing file should not error, got %v", err)
	}
}

func TestEntry_IsExpired(t *testing.T) {
	past := token.Entry{ExpiresAt: time.Now().Add(-time.Second)}
	if !past.IsExpired() {
		t.Error("expected past entry to be expired")
	}
	future := token.Entry{ExpiresAt: time.Now().Add(time.Hour)}
	if future.IsExpired() {
		t.Error("expected future entry to not be expired")
	}
	zero := token.Entry{}
	if zero.IsExpired() {
		t.Error("expected zero-time entry to not be expired")
	}
}
