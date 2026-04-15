package sync

import "strings"

// Labeler attaches computed labels to a profile based on its properties.
// Labels can be used downstream for routing, filtering, or observability.
type Labeler struct {
	rules []LabelRule
}

// LabelRule maps a predicate over a profile name/tags to a label value.
type LabelRule struct {
	Key       string
	Predicate func(profile string, tags []string) bool
	Value     string
}

// NewLabeler constructs a Labeler with the provided rules.
// If rules is nil an empty Labeler is returned that produces no labels.
func NewLabeler(rules []LabelRule) *Labeler {
	if rules == nil {
		rules = []LabelRule{}
	}
	return &Labeler{rules: rules}
}

// Label evaluates all rules against the given profile name and tags,
// returning a map of key→value pairs for every matching rule.
// Later rules overwrite earlier rules that share the same key.
func (l *Labeler) Label(profile string, tags []string) map[string]string {
	out := make(map[string]string)
	for _, r := range l.rules {
		if r.Predicate(profile, tags) {
			out[r.Key] = r.Value
		}
	}
	return out
}

// HasTag is a helper predicate that returns true when the given tag is present.
func HasTag(tag string) func(string, []string) bool {
	return func(_ string, tags []string) bool {
		for _, t := range tags {
			if t == tag {
				return true
			}
		}
		return false
	}
}

// ProfilePrefix is a helper predicate that returns true when the profile name
// starts with the given prefix.
func ProfilePrefix(prefix string) func(string, []string) bool {
	return func(profile string, _ []string) bool {
		return strings.HasPrefix(profile, prefix)
	}
}
