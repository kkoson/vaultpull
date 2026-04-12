package sync

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/fmdunlap/vaultpull/internal/config"
)

func buildPoolRunner(profiles []config.Profile, syncErr error) *Runner {
	r := buildTestRunner(profiles)
	if syncErr != nil {
		r.newSyncer = func(_ *config.Profile, _ Options) (syncer, error) {
			return &fakeSyncer{err: syncErr}, nil
		}
	}
	return r
}

func TestRunAllConcurrent_NoProfiles(t *testing.T) {
	r := buildPoolRunner(nil, nil)
	if err := r.RunAllConcurrent(context.Background(), 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunAllConcurrent_AllSuccess(t *testing.T) {
	profiles := []config.Profile{
		{Name: "dev", VaultPath: "secret/dev", OutputFile: ".env.dev"},
		{Name: "prod", VaultPath: "secret/prod", OutputFile: ".env.prod"},
	}
	r := buildPoolRunner(profiles, nil)
	if err := r.RunAllConcurrent(context.Background(), 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunAllConcurrent_PartialFailure(t *testing.T) {
	profiles := []config.Profile{
		{Name: "dev", VaultPath: "secret/dev", OutputFile: ".env.dev"},
		{Name: "prod", VaultPath: "secret/prod", OutputFile: ".env.prod"},
	}
	wantErr := errors.New("vault down")
	r := buildPoolRunner(profiles, wantErr)
	err := r.RunAllConcurrent(context.Background(), 2)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "concurrent sync errors") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunAllConcurrent_WorkerOne(t *testing.T) {
	profiles := []config.Profile{
		{Name: "a", VaultPath: "secret/a", OutputFile: ".env.a"},
		{Name: "b", VaultPath: "secret/b", OutputFile: ".env.b"},
		{Name: "c", VaultPath: "secret/c", OutputFile: ".env.c"},
	}
	r := buildPoolRunner(profiles, nil)
	if err := r.RunAllConcurrent(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error with single worker: %v", err)
	}
}
