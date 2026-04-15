package sync

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
)

// FingerprintStore tracks per-profile secret fingerprints to detect
// whether secrets have changed since the last sync.
type FingerprintStore struct {
	mu    sync.RWMutex
	store map[string]string
}

// NewFingerprintStore returns an initialised FingerprintStore.
func NewFingerprintStore() *FingerprintStore {
	return &FingerprintStore{
		store: make(map[string]string),
	}
}

// Compute returns a stable SHA-256 hex fingerprint for the given secret map.
// Keys are sorted before hashing so the result is order-independent.
func (f *FingerprintStore) Compute(secrets map[string]string) string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s;", k, secrets[k])
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Changed returns true when the fingerprint for profile differs from the
// last recorded value, or when no fingerprint has been recorded yet.
func (f *FingerprintStore) Changed(profile string, secrets map[string]string) bool {
	current := f.Compute(secrets)
	f.mu.RLock()
	prev, ok := f.store[profile]
	f.mu.RUnlock()
	if !ok {
		return true
	}
	return current != prev
}

// Record stores the fingerprint for profile so future calls to Changed
// can compare against it.
func (f *FingerprintStore) Record(profile string, secrets map[string]string) {
	fp := f.Compute(secrets)
	f.mu.Lock()
	f.store[profile] = fp
	f.mu.Unlock()
}

// Clear removes the stored fingerprint for profile.
func (f *FingerprintStore) Clear(profile string) {
	f.mu.Lock()
	delete(f.store, profile)
	f.mu.Unlock()
}
