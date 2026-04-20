package sync

import "strings"

// Suppressor holds a set of profile names that should be skipped during sync.
type Suppressor struct {
	suppressed map[string]struct{}
}

// NewSuppressor creates a Suppressor from a list of profile names.
// If names is nil or empty, no profiles are suppressed.
func NewSuppressor(names []string) *Suppressor {
	m := make(map[string]struct{}, len(names))
	for _, n := range names {
		key := strings.TrimSpace(strings.ToLower(n))
		if key != "" {
			m[key] = struct{}{}
		}
	}
	return &Suppressor{suppressed: m}
}

// IsSuppressed reports whether the given profile name is suppressed.
func (s *Suppressor) IsSuppressed(name string) bool {
	if s == nil || len(s.suppressed) == 0 {
		return false
	}
	key := strings.TrimSpace(strings.ToLower(name))
	_, ok := s.suppressed[key]
	return ok
}

// Add adds a profile name to the suppression list.
func (s *Suppressor) Add(name string) {
	key := strings.TrimSpace(strings.ToLower(name))
	if key != "" {
		s.suppressed[key] = struct{}{}
	}
}

// Remove removes a profile name from the suppression list.
func (s *Suppressor) Remove(name string) {
	key := strings.TrimSpace(strings.ToLower(name))
	delete(s.suppressed, key)
}

// Len returns the number of suppressed profiles.
func (s *Suppressor) Len() int {
	if s == nil {
		return 0
	}
	return len(s.suppressed)
}
