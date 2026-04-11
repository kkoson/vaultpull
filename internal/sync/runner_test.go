package sync

import (
	"context"
	"errors"
	"testing"

	"github.com/yourusername/vaultpull/internal/audit"
	"github.com/yourusername/vaultpull/internal/config"
)

func buildTestRunner(cfg *config.Config) (*Runner, *audit.Logger) {
	logger := audit.NewLogger(nil)
	r := NewRunner(cfg, logger)
	return r, logger
}

func TestRunner_RunProfile_ProfileNotFound(t *testing.T) {
	cfg := &config.Config{}
	r, _ := buildTestRunner(cfg)

	err := r.RunProfile(context.Background(), "missing", DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing profile, got nil")
	}
}

func TestRunner_RunProfile_SyncerError(t *testing.T) {
	cfg := &config.Config{
		Profiles: []config.Profile{
			{Name: "dev", VaultPath: "secret/dev", OutputFile: ".env"},
		},
	}
	r, _ := buildTestRunner(cfg)
	syncErr := errors.New("vault unavailable")
	r.newSync = func(p *config.Profile, opts Options) (*Syncer, error) {
		return nil, syncErr
	}

	err := r.RunProfile(context.Background(), "dev", DefaultOptions())
	if err == nil {
		t.Fatal("expected error from syncer constructor, got nil")
	}
}

func TestRunner_RunAll_AllSuccess(t *testing.T) {
	cfg := &config.Config{
		Profiles: []config.Profile{
			{Name: "dev", VaultPath: "secret/dev", OutputFile: ".env.dev"},
			{Name: "prod", VaultPath: "secret/prod", OutputFile: ".env.prod"},
		},
	}
	r, _ := buildTestRunner(cfg)
	r.newSync = func(p *config.Profile, opts Options) (*Syncer, error) {
		s := &Syncer{runFn: func(ctx context.Context) error { return nil }}
		return s, nil
	}

	if err := r.RunAll(context.Background(), DefaultOptions()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunner_RunAll_PartialFailure(t *testing.T) {
	cfg := &config.Config{
		Profiles: []config.Profile{
			{Name: "dev", VaultPath: "secret/dev", OutputFile: ".env.dev"},
			{Name: "prod", VaultPath: "secret/prod", OutputFile: ".env.prod"},
		},
	}
	r, _ := buildTestRunner(cfg)
	r.newSync = func(p *config.Profile, opts Options) (*Syncer, error) {
		if p.Name == "prod" {
			return nil, errors.New("prod vault error")
		}
		s := &Syncer{runFn: func(ctx context.Context) error { return nil }}
		return s, nil
	}

	err := r.RunAll(context.Background(), DefaultOptions())
	if err == nil {
		t.Fatal("expected partial failure error, got nil")
	}
}
