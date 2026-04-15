package sync

import (
	"testing"
	"time"
)

func TestDefaultAdmissionPolicy_Values(t *testing.T) {
	p := DefaultAdmissionPolicy()
	if p.maxFailures != 3 {
		t.Errorf("expected maxFailures=3, got %d", p.maxFailures)
	}
	if p.window != 2*time.Minute {
		t.Errorf("expected window=2m, got %s", p.window)
	}
}

func TestNewAdmissionPolicy_ZeroFallsBackToDefaults(t *testing.T) {
	p := NewAdmissionPolicy(0, 0)
	if p.maxFailures != 3 {
		t.Errorf("expected maxFailures=3, got %d", p.maxFailures)
	}
	if p.window != 2*time.Minute {
		t.Errorf("expected window=2m, got %s", p.window)
	}
}

func TestAdmissionPolicy_Admit_AllowsInitially(t *testing.T) {
	p := DefaultAdmissionPolicy()
	if err := p.Admit("dev"); err != nil {
		t.Errorf("expected nil error on fresh profile, got %v", err)
	}
}

func TestAdmissionPolicy_Admit_DeniedAfterMaxFailures(t *testing.T) {
	p := NewAdmissionPolicy(2, time.Minute)
	p.RecordFailure("staging")
	p.RecordFailure("staging")

	if err := p.Admit("staging"); err == nil {
		t.Error("expected admission to be denied after max failures")
	}
}

func TestAdmissionPolicy_Admit_AllowsAfterWindowExpires(t *testing.T) {
	p := NewAdmissionPolicy(1, 50*time.Millisecond)
	p.RecordFailure("prod")

	if err := p.Admit("prod"); err == nil {
		t.Error("expected admission denied immediately after failure")
	}

	time.Sleep(60 * time.Millisecond)

	if err := p.Admit("prod"); err != nil {
		t.Errorf("expected admission allowed after window expired, got %v", err)
	}
}

func TestAdmissionPolicy_Reset_ClearsFailures(t *testing.T) {
	p := NewAdmissionPolicy(1, time.Minute)
	p.RecordFailure("qa")

	if err := p.Admit("qa"); err == nil {
		t.Error("expected admission denied before reset")
	}

	p.Reset("qa")

	if err := p.Admit("qa"); err != nil {
		t.Errorf("expected admission allowed after reset, got %v", err)
	}
}

func TestAdmissionPolicy_IsolatedPerProfile(t *testing.T) {
	p := NewAdmissionPolicy(1, time.Minute)
	p.RecordFailure("alpha")

	if err := p.Admit("beta"); err != nil {
		t.Errorf("expected beta to be unaffected by alpha failures, got %v", err)
	}
}
