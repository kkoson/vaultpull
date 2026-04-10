package config

import (
	"testing"
)

func TestProfile_Validate_Valid(t *testing.T) {
	p := &Profile{
		Name:       "dev",
		VaultPath:  "apps/dev/myapp",
		OutputFile: ".env",
	}
	if err := p.Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestProfile_Validate_MissingName(t *testing.T) {
	p := &Profile{
		VaultPath:  "apps/dev/myapp",
		OutputFile: ".env",
	}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing name, got nil")
	}
}

func TestProfile_Validate_MissingVaultPath(t *testing.T) {
	p := &Profile{
		Name:       "dev",
		OutputFile: ".env",
	}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing vault_path, got nil")
	}
}

func TestProfile_Validate_MissingOutputFile(t *testing.T) {
	p := &Profile{
		Name:      "dev",
		VaultPath: "apps/dev/myapp",
	}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing output_file, got nil")
	}
}

func TestProfile_DefaultMountPath_Empty(t *testing.T) {
	p := &Profile{MountPath: ""}
	if got := p.DefaultMountPath(); got != "secret" {
		t.Fatalf("expected \"secret\", got %q", got)
	}
}

func TestProfile_DefaultMountPath_Custom(t *testing.T) {
	p := &Profile{MountPath: "kv"}
	if got := p.DefaultMountPath(); got != "kv" {
		t.Fatalf("expected \"kv\", got %q", got)
	}
}

func TestGetProfile_Found(t *testing.T) {
	cfg := &Config{
		Profiles: []Profile{
			{Name: "staging", VaultPath: "apps/staging", OutputFile: ".env.staging"},
			{Name: "prod", VaultPath: "apps/prod", OutputFile: ".env.prod"},
		},
	}
	p, err := cfg.GetProfile("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "prod" {
		t.Fatalf("expected profile \"prod\", got %q", p.Name)
	}
}

func TestGetProfile_Missing(t *testing.T) {
	cfg := &Config{
		Profiles: []Profile{
			{Name: "dev", VaultPath: "apps/dev", OutputFile: ".env"},
		},
	}
	_, err := cfg.GetProfile("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing profile, got nil")
	}
}
