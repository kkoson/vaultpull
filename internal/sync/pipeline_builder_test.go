package sync

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
)

func TestNewPipelineBuilder_NotNil(t *testing.T) {
	r := buildTestRunner(t, nil)
	b := NewPipelineBuilder(r, DefaultOptions())
	if b == nil {
		t.Fatal("expected non-nil PipelineBuilder")
	}
}

func TestPipelineBuilder_Build_ReturnsTwoStages(t *testing.T) {
	r := buildTestRunner(t, nil)
	b := NewPipelineBuilder(r, DefaultOptions())
	profile := config.Profile{
		Name:       "dev",
		VaultPath:  "secret/dev",
		OutputFile: ".env",
	}
	p := b.Build(profile)
	if p.Len() != 2 {
		t.Fatalf("expected 2 stages, got %d", p.Len())
	}
}

func TestPipelineBuilder_Build_ValidateStageFailsOnInvalidProfile(t *testing.T) {
	r := buildTestRunner(t, nil)
	b := NewPipelineBuilder(r, DefaultOptions())
	// Profile missing required fields — Validate should fail.
	invalid := config.Profile{Name: "bad"}
	p := b.Build(invalid)
	// Only run the validate stage by building a one-stage pipeline copy.
	validateStage := Stage{
		Name: "validate",
		Run:  p.stages[0].Run,
	}
	sp := NewPipeline(validateStage)
	if err := sp.Execute(t.Context(), "bad"); err == nil {
		t.Fatal("expected validation error for incomplete profile")
	}
}
