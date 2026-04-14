package sync

import (
	"regexp"
	"strings"
)

// RedactConfig controls which secret keys are redacted in logs and output.
type RedactConfig struct {
	// Patterns is a list of regex patterns matching key names to redact.
	Patterns []string
	// Placeholder is substituted for redacted values.
	Placeholder string
}

// DefaultRedactConfig returns a RedactConfig with sensible defaults.
func DefaultRedactConfig() RedactConfig {
	return RedactConfig{
		Patterns:    []string{`(?i)password`, `(?i)secret`, `(?i)token`, `(?i)api_?key`, `(?i)private`},
		Placeholder: "[REDACTED]",
	}
}

// Redactor masks secret values whose keys match configured patterns.
type Redactor struct {
	cfg      RedactConfig
	compiled []*regexp.Regexp
}

// NewRedactor compiles the patterns in cfg and returns a Redactor.
// If cfg.Placeholder is empty it falls back to "[REDACTED]".
func NewRedactor(cfg RedactConfig) *Redactor {
	if cfg.Placeholder == "" {
		cfg.Placeholder = "[REDACTED]"
	}
	var compiled []*regexp.Regexp
	for _, p := range cfg.Patterns {
		if re, err := regexp.Compile(p); err == nil {
			compiled = append(compiled, re)
		}
	}
	return &Redactor{cfg: cfg, compiled: compiled}
}

// Redact returns a copy of secrets where matched keys have their values replaced.
func (r *Redactor) Redact(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if r.matches(k) {
			out[k] = r.cfg.Placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

// RedactValue returns the placeholder if key matches, otherwise the original value.
func (r *Redactor) RedactValue(key, value string) string {
	if r.matches(key) {
		return r.cfg.Placeholder
	}
	return value
}

// IsSensitive reports whether a key name matches any redaction pattern.
func (r *Redactor) IsSensitive(key string) bool {
	return r.matches(key)
}

func (r *Redactor) matches(key string) bool {
	for _, re := range r.compiled {
		if re.MatchString(strings.TrimSpace(key)) {
			return true
		}
	}
	return false
}
