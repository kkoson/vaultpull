package sync

import "sync"

// PinStore records a pinned secret version per profile.
// A pinned profile will not be re-fetched from Vault until unpinned.
type PinStore struct {
	mu   sync.RWMutex
	pins map[string]string
}

// NewPinStore returns an initialised PinStore.
func NewPinStore() *PinStore {
	return &PinStore{pins: make(map[string]string)}
}

// Pin records version v for the given profile name.
// Calling Pin again overwrites any existing pin.
func (s *PinStore) Pin(profile, version string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pins[profile] = version
}

// Get returns (true, version) when profile is pinned, otherwise (false, "").
func (s *PinStore) Get(profile string) (bool, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.pins[profile]
	return ok, v
}

// Unpin removes any pin for profile. It is a no-op if profile is not pinned.
func (s *PinStore) Unpin(profile string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.pins, profile)
}

// Pinned returns a snapshot of all currently pinned profiles.
func (s *PinStore) Pinned() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]string, len(s.pins))
	for k, v := range s.pins {
		out[k] = v
	}
	return out
}
