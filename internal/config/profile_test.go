package config

import "testing"

func TestProfile_Validate_Valid(t *testing.T) {
	p := &Profile{
		Name:       "dev",
		VaultPath:  "secret/dev",
		OutputFile: ".env",
	}
	if err := p.Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestProfile_Validate_MissingName(t *testing.T) {
	p := &Profile{VaultPath: "secret/dev", OutputFile: ".env"}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestProfile_Validate_MissingVaultPath(t *testing.T) {
	p := &Profile{Name: "dev", OutputFile: ".env"}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing vault_path")
	}
}

func TestProfile_Validate_MissingOutputFile(t *testing.T) {
	p := &Profile{Name: "dev", VaultPath: "secret/dev"}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing output_file")
	}
}

func TestProfile_DefaultMountPath_Empty(t *testing.T) {
	p := &Profile{}
	if got := p.DefaultMountPath(); got != "secret" {
		t.Fatalf("expected \"secret\", got %q", got)
	}
}

func TestProfile_DefaultMountPath_Set(t *testing.T) {
	p := &Profile{MountPath: "kv"}
	if got := p.DefaultMountPath(); got != "kv" {
		t.Fatalf("expected \"kv\", got %q", got)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	cfg := &Config{Profiles: []Profile{{Name: "dev", VaultPath: "v", OutputFile: "o"}}}
	if _, err := cfg.GetProfile("prod"); err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestGetProfile_Default(t *testing.T) {
	cfg := &Config{Profiles: []Profile{{Name: "dev", VaultPath: "v", OutputFile: "o"}}}
	p, err := cfg.GetDefaultProfile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "dev" {
		t.Fatalf("expected \"dev\", got %q", p.Name)
	}
}

func TestGetProfile_DefaultMultiple(t *testing.T) {
	cfg := &Config{
		Profiles: []Profile{
			{Name: "dev", VaultPath: "v", OutputFile: "o"},
			{Name: "prod", VaultPath: "v2", OutputFile: "o2"},
		},
	}
	if _, err := cfg.GetDefaultProfile(); err == nil {
		t.Fatal("expected error when multiple profiles defined")
	}
}
