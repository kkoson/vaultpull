package sync

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
)

func buildValidatorConfig(profiles []config.Profile) *config.Config {
	return &config.Config{Profiles: profiles}
}

func TestNewValidator_NotNil(t *testing.T) {
	v := NewValidator(&config.Config{})
	if v == nil {
		t.Fatal("expected non-nil Validator")
	}
}

func TestValidateAll_NilConfig(t *testing.T) {
	v := &Validator{cfg: nil}
	_, err := v.ValidateAll()
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestValidateAll_NoProfiles(t *testing.T) {
	v := NewValidator(buildValidatorConfig(nil))
	results, err := v.ValidateAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestValidateAll_ValidProfile(t *testing.T) {
	profiles := []config.Profile{
		{Name: "prod", VaultPath: "secret/prod", OutputFile: ".env"},
	}
	v := NewValidator(buildValidatorConfig(profiles))
	results, err := v.ValidateAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Valid() {
		t.Errorf("expected valid profile, got errors: %v", results[0].Errors)
	}
}

func TestValidateAll_InvalidProfile_MissingOutputFile(t *testing.T) {
	profiles := []config.Profile{
		{Name: "dev", VaultPath: "secret/dev", OutputFile: ""},
	}
	v := NewValidator(buildValidatorConfig(profiles))
	results, err := v.ValidateAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Valid() {
		t.Error("expected invalid result for missing output_file")
	}
}

func TestValidateProfile_NotFound(t *testing.T) {
	v := NewValidator(buildValidatorConfig(nil))
	_, err := v.ValidateProfile("ghost")
	if err == nil {
		t.Fatal("expected error for unknown profile")
	}
}

func TestValidateProfile_Found_Valid(t *testing.T) {
	profiles := []config.Profile{
		{Name: "staging", VaultPath: "secret/staging", OutputFile: ".env.staging"},
	}
	v := NewValidator(buildValidatorConfig(profiles))
	result, err := v.ValidateProfile("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Valid() {
		t.Errorf("expected valid result, got errors: %v", result.Errors)
	}
	if result.Profile != "staging" {
		t.Errorf("expected profile name 'staging', got %q", result.Profile)
	}
}
