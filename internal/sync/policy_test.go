package sync

import (
	"testing"
	"time"
)

func TestDefaultRetentionPolicy_Values(t *testing.T) {
	p := DefaultRetentionPolicy()
	if p.MaxAge != 24*time.Hour {
		t.Errorf("expected 24h, got %v", p.MaxAge)
	}
	if p.EnforceOnFailure {
		t.Error("expected EnforceOnFailure to be false")
	}
}

func TestNewPolicyEnforcer_ZeroMaxAgeUsesDefault(t *testing.T) {
	e := NewPolicyEnforcer(RetentionPolicy{})
	if e.policy.MaxAge != 24*time.Hour {
		t.Errorf("expected default MaxAge, got %v", e.policy.MaxAge)
	}
}

func TestPolicyEnforcer_Allow_NeverSynced(t *testing.T) {
	e := NewPolicyEnforcer(DefaultRetentionPolicy())
	if !e.Allow("prod") {
		t.Error("expected Allow=true for never-synced profile")
	}
}

func TestPolicyEnforcer_Allow_RecentSync(t *testing.T) {
	e := NewPolicyEnforcer(RetentionPolicy{MaxAge: time.Hour})
	e.Record("prod")
	if !e.Allow("prod") {
		t.Error("expected Allow=true for recently synced profile")
	}
}

func TestPolicyEnforcer_Allow_ExpiredSync(t *testing.T) {
	e := NewPolicyEnforcer(RetentionPolicy{MaxAge: time.Millisecond})
	e.Record("prod")
	time.Sleep(5 * time.Millisecond)
	if e.Allow("prod") {
		t.Error("expected Allow=false for expired profile")
	}
}

func TestPolicyEnforcer_Age_NeverSynced(t *testing.T) {
	e := NewPolicyEnforcer(DefaultRetentionPolicy())
	if e.Age("prod") != -1 {
		t.Error("expected Age=-1 for never-synced profile")
	}
}

func TestPolicyEnforcer_Age_AfterRecord(t *testing.T) {
	e := NewPolicyEnforcer(DefaultRetentionPolicy())
	e.Record("prod")
	age := e.Age("prod")
	if age < 0 {
		t.Errorf("expected non-negative age, got %v", age)
	}
	if age > time.Second {
		t.Errorf("expected age < 1s, got %v", age)
	}
}

func TestPolicyEnforcer_Record_Overwrite(t *testing.T) {
	e := NewPolicyEnforcer(RetentionPolicy{MaxAge: time.Millisecond})
	e.Record("prod")
	time.Sleep(5 * time.Millisecond)
	e.Record("prod") // refresh
	if !e.Allow("prod") {
		t.Error("expected Allow=true after re-record")
	}
}
