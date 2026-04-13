package sync

// TagFilter selects profiles based on user-defined tags attached to each profile.
// A profile matches if it carries at least one of the required tags, or if no
// tags are required (the filter is open).
type TagFilter struct {
	required map[string]struct{}
}

// NewTagFilter returns a TagFilter that passes only profiles whose tags
// overlap with the provided list. An empty or nil list means all profiles pass.
func NewTagFilter(tags []string) *TagFilter {
	required := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		if t != "" {
			required[t] = struct{}{}
		}
	}
	return &TagFilter{required: required}
}

// Allow returns true when the profile should be processed.
// It returns true unconditionally when no tags were configured.
func (f *TagFilter) Allow(profileTags []string) bool {
	if len(f.required) == 0 {
		return true
	}
	for _, t := range profileTags {
		if _, ok := f.required[t]; ok {
			return true
		}
	}
	return false
}

// RequiredTags returns a copy of the tag set the filter was built with.
func (f *TagFilter) RequiredTags() []string {
	out := make([]string, 0, len(f.required))
	for t := range f.required {
		out = append(out, t)
	}
	return out
}
