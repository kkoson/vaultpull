package sync

import (
	"errors"
	"testing"
)

func TestNewInvariantChecker_NotNil(t *testing.T) {
	c := NewInvariantChecker()
	if c == nil {
		t.Fatal("expected non-nil checker")
	}
}

func TestInvariantChecker_Len_Empty(t *testing.T) {
	c := NewInvariantChecker()
	if c.Len() != 0 {
		t.Fatalf("expected 0, got %d", c.Len())
	}
}

func TestInvariantChecker_Register_NilIgnored(t *testing.T) {
	c := NewInvariantChecker()
	c.Register(InvariantRule{Name: "noop", Check: nil})
	if c.Len() != 0 {
		t.Fatal("nil check should be ignored")
	}
}

func TestInvariantChecker_Register_IncrementsLen(t *testing.T) {
	c := NewInvariantChecker()
	c.Register(RequireKeys("FOO"))
	c.Register(ForbidKeys("BAR"))
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}

func TestInvariantChecker_Check_AllPass(t *testing.T) {
	c := NewInvariantChecker()
	c.Register(RequireKeys("DB_URL"))
	secrets := map[string]string{"DB_URL": "postgres://localhost"}
	if err := c.Check(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInvariantChecker_Check_RequireKeys_Missing(t *testing.T) {
	c := NewInvariantChecker()
	c.Register(RequireKeys("DB_URL", "API_KEY"))
	secrets := map[string]string{"DB_URL": "postgres://localhost"}
	err := c.Check(secrets)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if !contains(err.Error(), "API_KEY") {
		t.Fatalf("expected error to mention API_KEY, got: %v", err)
	}
}

func TestInvariantChecker_Check_ForbidKeys_Present(t *testing.T) {
	c := NewInvariantChecker()
	c.Register(ForbidKeys("DEBUG"))
	secrets := map[string]string{"DEBUG": "true"}
	err := c.Check(secrets)
	if err == nil {
		t.Fatal("expected error for forbidden key")
	}
}

func TestInvariantChecker_Check_CustomRule(t *testing.T) {
	c := NewInvariantChecker()
	c.Register(InvariantRule{
		Name: "non-empty-token",
		Check: func(s map[string]string) error {
			if s["TOKEN"] == "" {
				return errors.New("TOKEN must not be empty")
			}
			return nil
		},
	})
	if err := c.Check(map[string]string{"TOKEN": ""}); err == nil {
		t.Fatal("expected error")
	}
	if err := c.Check(map[string]string{"TOKEN": "abc"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInvariantChecker_Check_StopsAtFirstViolation(t *testing.T) {
	c := NewInvariantChecker()
	calls := 0
	c.Register(InvariantRule{Name: "fail-first", Check: func(_ map[string]string) error {
		calls++
		return errors.New("first fails")
	}})
	c.Register(InvariantRule{Name: "second", Check: func(_ map[string]string) error {
		calls++
		return nil
	}})
	c.Check(map[string]string{})
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
