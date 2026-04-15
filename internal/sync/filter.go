package sync

import "strings"

// Filter controls which profiles are included or excluded during a sync run.
type Filter struct {
	// Include is an optional list of profile names to sync. If empty, all profiles are included.
	Include []string
	// Exclude is an optional list of profile names to skip.
	Exclude []string
}

// NewFilter returns a Filter with the given include and exclude lists.
func NewFilter(include, exclude []string) *Filter {
	return &Filter{
		Include: include,
		Exclude: exclude,
	}
}

// Allow reports whether the given profile name should be processed.
// Exclusions take precedence over inclusions.
func (f *Filter) Allow(name string) bool {
	if f == nil {
		return true
	}

	name = strings.TrimSpace(name)

	for _, ex := range f.Exclude {
		if strings.EqualFold(ex, name) {
			return false
		}
	}

	if len(f.Include) == 0 {
		return true
	}

	for _, inc := range f.Include {
		if strings.EqualFold(inc, name) {
			return true
		}
	}

	return false
}

// AllowedProfiles filters a slice of profile names using the filter rules.
func (f *Filter) AllowedProfiles(names []string) []string {
	out := make([]string, 0, len(names))
	for _, n := range names {
		if f.Allow(n) {
			out = append(out, n)
		}
	}
	return out
}

// IsEmpty reports whether the filter has no include or exclude rules defined.
// An empty filter allows all profiles to pass through.
func (f *Filter) IsEmpty() bool {
	if f == nil {
		return true
	}
	return len(f.Include) == 0 && len(f.Exclude) == 0
}
