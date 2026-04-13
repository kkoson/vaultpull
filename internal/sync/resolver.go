package sync

import (
	"fmt"
	"strings"

	"github.com/yourusername/vaultpull/internal/config"
)

// ResolveMode controls how profile names are resolved.
type ResolveMode int

const (
	// ResolveModeExact requires an exact match on profile name.
	ResolveModeExact ResolveMode = iota
	// ResolveModePrefix matches any profile whose name starts with the given prefix.
	ResolveModePrefix
)

// Resolver resolves profile names from a config to concrete Profile values.
type Resolver struct {
	cfg  *config.Config
	mode ResolveMode
}

// NewResolver returns a Resolver backed by cfg.
// mode controls exact vs prefix matching.
func NewResolver(cfg *config.Config, mode ResolveMode) *Resolver {
	return &Resolver{cfg: cfg, mode: mode}
}

// Resolve returns all profiles matching name according to the resolver's mode.
// Returns an error if no profiles match.
func (r *Resolver) Resolve(name string) ([]config.Profile, error) {
	if r.cfg == nil {
		return nil, fmt.Errorf("resolver: nil config")
	}

	var matched []config.Profile
	for _, p := range r.cfg.Profiles {
		switch r.mode {
		case ResolveModeExact:
			if p.Name == name {
				matched = append(matched, p)
			}
		case ResolveModePrefix:
			if strings.HasPrefix(p.Name, name) {
				matched = append(matched, p)
			}
		}
	}

	if len(matched) == 0 {
		return nil, fmt.Errorf("resolver: no profiles matched %q", name)
	}
	return matched, nil
}

// ResolveAll returns every profile in the config.
func (r *Resolver) ResolveAll() ([]config.Profile, error) {
	if r.cfg == nil {
		return nil, fmt.Errorf("resolver: nil config")
	}
	if len(r.cfg.Profiles) == 0 {
		return nil, fmt.Errorf("resolver: config contains no profiles")
	}
	return r.cfg.Profiles, nil
}
