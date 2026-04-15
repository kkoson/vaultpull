package sync

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/config"
)

func buildPolicyProfile(name string) config.Profile {
	return config.Profile{Name: name, VaultPath: "secret/data/test", OutputFile: ".env"}
}

func TestWithRetentionPolicy_PanicsOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on nil enforcer")
		}
	}()
	WithRetentionPolicy(nil)
}

func TestWithRetentionPolicy_AllowsNeverSynced(t *testing.T) {
	enforcer := NewPolicyEnforcer(DefaultRetentionPolicy())
	called := false
	stage := WithRetentionPolicy(enforcer)(func(_ context.Context, _ config.Profile) error {
		called = true
		return nil
	})
	if err := stage(context.Background(), buildPolicyProfile("prod")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected inner stage to be called")
	}
}

func TestWithRetentionPolicy_DeniesExpiredSync(t *testing.T) {
	enforcer := NewPolicyEnforcer(RetentionPolicy{MaxAge: time.Millisecond})
	enforcer.Record("prod")
	time.Sleep(5 * time.Millisecond)

	stage := WithRetentionPolicy(enforcer)(func(_ context.Context, _ config.Profile) error {
		return nil
	})
	err := stage(context.Background(), buildPolicyProfile("prod"))
	if !errors.Is(err, ErrPolicyDenied) {
		t.Errorf("expected ErrPolicyDenied, got %v", err)
	}
}

func TestWithRetentionPolicy_RecordsOnSuccess(t *testing.T) {
	enforcer := NewPolicyEnforcer(DefaultRetentionPolicy())
	stage := WithRetentionPolicy(enforcer)(func(_ context.Context, _ config.Profile) error {
		return nil
	})
	_ = stage(context.Background(), buildPolicyProfile("staging"))
	if enforcer.Age("staging") == -1 {
		t.Error("expected age to be recorded after successful sync")
	}
}

func TestWithRetentionPolicy_DoesNotRecordOnError(t *testing.T) {
	enforcer := NewPolicyEnforcer(DefaultRetentionPolicy())
	synErr := errors.New("vault error")
	stage := WithRetentionPolicy(enforcer)(func(_ context.Context, _ config.Profile) error {
		return synErr
	})
	_ = stage(context.Background(), buildPolicyProfile("dev"))
	if enforcer.Age("dev") != -1 {
		t.Error("expected no age recorded after failed sync")
	}
}
