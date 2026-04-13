package sync

import (
	"testing"
)

func TestNewTagFilter_NilAllowsAll(t *testing.T) {
	f := NewTagFilter(nil)
	if !f.Allow([]string{"prod", "eu"}) {
		t.Fatal("expected nil tag filter to allow everything")
	}
}

func TestNewTagFilter_EmptyAllowsAll(t *testing.T) {
	f := NewTagFilter([]string{})
	if !f.Allow([]string{"staging"}) {
		t.Fatal("expected empty tag filter to allow everything")
	}
}

func TestTagFilter_Allow_MatchingTag(t *testing.T) {
	f := NewTagFilter([]string{"prod", "eu"})
	if !f.Allow([]string{"eu", "critical"}) {
		t.Fatal("expected profile with matching tag 'eu' to be allowed")
	}
}

func TestTagFilter_Allow_NoMatchingTag(t *testing.T) {
	f := NewTagFilter([]string{"prod"})
	if f.Allow([]string{"staging", "eu"}) {
		t.Fatal("expected profile without matching tag to be denied")
	}
}

func TestTagFilter_Allow_EmptyProfileTags(t *testing.T) {
	f := NewTagFilter([]string{"prod"})
	if f.Allow([]string{}) {
		t.Fatal("expected profile with no tags to be denied when filter is set")
	}
}

func TestTagFilter_Allow_NilProfileTags(t *testing.T) {
	f := NewTagFilter([]string{"prod"})
	if f.Allow(nil) {
		t.Fatal("expected profile with nil tags to be denied when filter is set")
	}
}

func TestTagFilter_RequiredTags_ReturnsAll(t *testing.T) {
	f := NewTagFilter([]string{"a", "b", "c"})
	tags := f.RequiredTags()
	if len(tags) != 3 {
		t.Fatalf("expected 3 required tags, got %d", len(tags))
	}
}

func TestTagFilter_RequiredTags_EmptyWhenNilInput(t *testing.T) {
	f := NewTagFilter(nil)
	if len(f.RequiredTags()) != 0 {
		t.Fatal("expected no required tags for nil input")
	}
}

func TestTagFilter_IgnoresBlankTags(t *testing.T) {
	f := NewTagFilter([]string{"", "prod", ""})
	if len(f.RequiredTags()) != 1 {
		t.Fatalf("expected blank tags to be ignored, got %d required", len(f.RequiredTags()))
	}
	if !f.Allow([]string{"prod"}) {
		t.Fatal("expected 'prod' to be allowed")
	}
}
