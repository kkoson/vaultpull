// Package token provides Vault token caching to avoid repeated logins.
package token

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// ErrCacheMiss is returned when no valid cached token is found.
var ErrCacheMiss = errors.New("token: cache miss")

// Entry holds a cached Vault token and its expiry.
type Entry struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IsExpired reports whether the token has expired.
func (e Entry) IsExpired() bool {
	return !e.ExpiresAt.IsZero() && time.Now().After(e.ExpiresAt)
}

// Cache persists and retrieves Vault tokens from disk.
type Cache struct {
	path string
}

// NewCache creates a Cache that stores its file at the given path.
// If path is empty, a default under the user cache dir is used.
func NewCache(path string) (*Cache, error) {
	if path == "" {
		dir, err := os.UserCacheDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(dir, "vaultpull", "token.json")
	}
	return &Cache{path: path}, nil
}

// Load reads the cached token. Returns ErrCacheMiss if absent or expired.
func (c *Cache) Load() (Entry, error) {
	data, err := os.ReadFile(c.path)
	if errors.Is(err, os.ErrNotExist) {
		return Entry{}, ErrCacheMiss
	}
	if err != nil {
		return Entry{}, err
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return Entry{}, ErrCacheMiss
	}
	if e.IsExpired() {
		_ = c.Clear()
		return Entry{}, ErrCacheMiss
	}
	return e, nil
}

// Save writes the token entry to disk, creating parent directories as needed.
func (c *Cache) Save(e Entry) error {
	if err := os.MkdirAll(filepath.Dir(c.path), 0o700); err != nil {
		return err
	}
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0o600)
}

// Clear removes the cached token file.
func (c *Cache) Clear() error {
	err := os.Remove(c.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
