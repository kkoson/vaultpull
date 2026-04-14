package sync

import (
	"testing"
)

func TestDefaultRedactConfig_Values(t *testing.T) {
	cfg := DefaultRedactConfig()
	if cfg.Placeholder != "[REDACTED]" {
		t.Fatalf("expected [REDACTED], got %q", cfg.Placeholder)
	}
	if len(cfg.Patterns) == 0 {
		t.Fatal("expected non-empty patterns")
	}
}

func TestNewRedactor_NotNil(t *testing.T) {
	r := NewRedactor(DefaultRedactConfig())
	if r == nil {
		t.Fatal("expected non-nil Redactor")
	}
}

func TestNewRedactor_EmptyPlaceholderFallback(t *testing.T) {
	r := NewRedactor(RedactConfig{Placeholder: ""})
	if r.cfg.Placeholder != "[REDACTED]" {
		t.Fatalf("expected fallback placeholder, got %q", r.cfg.Placeholder)
	}
}

func TestRedactor_IsSensitive_MatchesPassword(t *testing.T) {
	r := NewRedactor(DefaultRedactConfig())
	if !r.IsSensitive("DB_PASSWORD") {
		t.Fatal("expected DB_PASSWORD to be sensitive")
	}
}

func TestRedactor_IsSensitive_MatchesToken(t *testing.T) {
	r := NewRedactor(DefaultRedactConfig())
	if !r.IsSensitive("VAULT_TOKEN") {
		t.Fatal("expected VAULT_TOKEN to be sensitive")
	}
}

func TestRedactor_IsSensitive_SafeKey(t *testing.T) {
	r := NewRedactor(DefaultRedactConfig())
	if r.IsSensitive("DATABASE_HOST") {
		t.Fatal("expected DATABASE_HOST to not be sensitive")
	}
}

func TestRedactor_RedactValue_SensitiveKey(t *testing.T) {
	r := NewRedactor(DefaultRedactConfig())
	got := r.RedactValue("API_SECRET", "super-secret")
	if got != "[REDACTED]" {
		t.Fatalf("expected [REDACTED], got %q", got)
	}
}

func TestRedactor_RedactValue_SafeKey(t *testing.T) {
	r := NewRedactor(DefaultRedactConfig())
	got := r.RedactValue("APP_ENV", "production")
	if got != "production" {
		t.Fatalf("expected production, got %q", got)
	}
}

func TestRedactor_Redact_MixedMap(t *testing.T) {
	r := NewRedactor(DefaultRedactConfig())
	input := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "hunter2",
		"API_KEY":     "abc123",
	}
	out := r.Redact(input)
	if out["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be unchanged")
	}
	if out["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("DB_PASSWORD should be redacted")
	}
	if out["API_KEY"] != "[REDACTED]" {
		t.Errorf("API_KEY should be redacted")
	}
}

func TestRedactor_Redact_EmptyMap(t *testing.T) {
	r := NewRedactor(DefaultRedactConfig())
	out := r.Redact(map[string]string{})
	if len(out) != 0 {
		t.Fatal("expected empty output map")
	}
}

func TestRedactor_CustomPattern(t *testing.T) {
	r := NewRedactor(RedactConfig{
		Patterns:    []string{`(?i)^my_custom_key$`},
		Placeholder: "***",
	})
	if !r.IsSensitive("MY_CUSTOM_KEY") {
		t.Fatal("expected MY_CUSTOM_KEY to match custom pattern")
	}
	if r.IsSensitive("OTHER_KEY") {
		t.Fatal("expected OTHER_KEY to not match")
	}
}
