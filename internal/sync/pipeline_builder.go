package sync

import (
	"context"

	"github.com/yourusername/vaultpull/internal/config"
)

// PipelineBuilder constructs a standard sync Pipeline for a given profile
// using the provided Runner. It wires together the canonical stages:
// validate → fetch → write.
type PipelineBuilder struct {
	runner *Runner
	opts   Options
}

// NewPipelineBuilder creates a PipelineBuilder backed by runner and opts.
func NewPipelineBuilder(runner *Runner, opts Options) *PipelineBuilder {
	return &PipelineBuilder{runner: runner, opts: opts}
}

// Build returns a Pipeline for the supplied profile.
func (b *PipelineBuilder) Build(profile config.Profile) *Pipeline {
	return NewPipeline(
		Stage{
			Name: "validate",
			Run: func(ctx context.Context, name string) error {
				return profile.Validate()
			},
		},
		Stage{
			Name: "sync",
			Run: func(ctx context.Context, name string) error {
				return b.runner.RunProfile(ctx, name, b.opts)
			},
		},
	)
}
