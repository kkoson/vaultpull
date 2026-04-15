package sync

import (
	"testing"
)

func TestNewLabeler_NotNil(t *testing.T) {
	l := NewLabeler(nil)
	if l == nil {
		t.Fatal("expected non-nil Labeler")
	}
}

func TestLabeler_NoRules_EmptyMap(t *testing.T) {
	l := NewLabeler(nil)
	labels := l.Label("prod-api", []string{"critical"})
	if len(labels) != 0 {
		t.Fatalf("expected empty labels, got %v", labels)
	}
}

func TestLabeler_ProfilePrefix_Matches(t *testing.T) {
	l := NewLabeler([]LabelRule{
		{Key: "env", Predicate: ProfilePrefix("prod-"), Value: "production"},
	})
	labels := l.Label("prod-api", nil)
	if labels["env"] != "production" {
		t.Fatalf("expected env=production, got %q", labels["env"])
	}
}

func TestLabeler_ProfilePrefix_NoMatch(t *testing.T) {
	l := NewLabeler([]LabelRule{
		{Key: "env", Predicate: ProfilePrefix("prod-"), Value: "production"},
	})
	labels := l.Label("staging-api", nil)
	if _, ok := labels["env"]; ok {
		t.Fatal("expected no env label for non-matching profile")
	}
}

func TestLabeler_HasTag_Matches(t *testing.T) {
	l := NewLabeler([]LabelRule{
		{Key: "critical", Predicate: HasTag("critical"), Value: "true"},
	})
	labels := l.Label("any", []string{"critical", "db"})
	if labels["critical"] != "true" {
		t.Fatalf("expected critical=true, got %q", labels["critical"])
	}
}

func TestLabeler_HasTag_NoMatch(t *testing.T) {
	l := NewLabeler([]LabelRule{
		{Key: "critical", Predicate: HasTag("critical"), Value: "true"},
	})
	labels := l.Label("any", []string{"db"})
	if _, ok := labels["critical"]; ok {
		t.Fatal("expected no critical label")
	}
}

func TestLabeler_LaterRuleOverwrites(t *testing.T) {
	l := NewLabeler([]LabelRule{
		{Key: "env", Predicate: ProfilePrefix(""), Value: "default"},
		{Key: "env", Predicate: ProfilePrefix("prod-"), Value: "production"},
	})
	labels := l.Label("prod-api", nil)
	if labels["env"] != "production" {
		t.Fatalf("expected env=production after overwrite, got %q", labels["env"])
	}
}

func TestLabeler_MultipleRules_AllMatch(t *testing.T) {
	l := NewLabeler([]LabelRule{
		{Key: "env", Predicate: ProfilePrefix("prod-"), Value: "production"},
		{Key: "critical", Predicate: HasTag("critical"), Value: "true"},
	})
	labels := l.Label("prod-api", []string{"critical"})
	if labels["env"] != "production" {
		t.Fatalf("expected env=production, got %q", labels["env"])
	}
	if labels["critical"] != "true" {
		t.Fatalf("expected critical=true, got %q", labels["critical"])
	}
}
