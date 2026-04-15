package sync

import (
	"testing"
)

func TestFilter_NilAllowsAll(t *testing.T) {
	var f *Filter
	if !f.Allow("production") {
		t.Error("nil filter should allow all profiles")
	}
}

func TestFilter_EmptyAllowsAll(t *testing.T) {
	f := NewFilter(nil, nil)
	profiles := []string{"dev", "staging", "prod"}
	allowed := f.AllowedProfiles(profiles)
	if len(allowed) != len(profiles) {
		t.Errorf("expected %d profiles, got %d", len(profiles), len(allowed))
	}
}

func TestFilter_IncludeList(t *testing.T) {
	f := NewFilter([]string{"dev", "staging"}, nil)

	if !f.Allow("dev") {
		t.Error("expected dev to be allowed")
	}
	if !f.Allow("staging") {
		t.Error("expected staging to be allowed")
	}
	if f.Allow("prod") {
		t.Error("expected prod to be excluded")
	}
}

func TestFilter_ExcludeList(t *testing.T) {
	f := NewFilter(nil, []string{"prod"})

	if !f.Allow("dev") {
		t.Error("expected dev to be allowed")
	}
	if f.Allow("prod") {
		t.Error("expected prod to be excluded")
	}
}

func TestFilter_ExcludeTakesPrecedence(t *testing.T) {
	f := NewFilter([]string{"dev", "prod"}, []string{"prod"})

	if !f.Allow("dev") {
		t.Error("expected dev to be allowed")
	}
	if f.Allow("prod") {
		t.Error("exclude should take precedence over include for prod")
	}
}

func TestFilter_CaseInsensitive(t *testing.T) {
	f := NewFilter([]string{"Dev"}, []string{"PROD"})

	if !f.Allow("dev") {
		t.Error("expected case-insensitive match for include")
	}
	if f.Allow("prod") {
		t.Error("expected case-insensitive match for exclude")
	}
}

func TestFilter_AllowedProfiles(t *testing.T) {
	f := NewFilter([]string{"dev", "staging"}, []string{"staging"})
	profiles := []string{"dev", "staging", "prod"}
	allowed := f.AllowedProfiles(profiles)

	if len(allowed) != 1 || allowed[0] != "dev" {
		t.Errorf("expected only [dev], got %v", allowed)
	}
}

func TestFilter_AllowedProfiles_EmptyInput(t *testing.T) {
	f := NewFilter([]string{"dev"}, nil)
	allowed := f.AllowedProfiles([]string{})

	if len(allowed) != 0 {
		t.Errorf("expected empty result for empty input, got %v", allowed)
	}
}
