package sync

import (
	"context"
	"fmt"
)

// Stage represents a single step in a sync pipeline.
type Stage struct {
	Name string
	Run  func(ctx context.Context, profileName string) error
}

// Pipeline executes a sequence of stages for each profile sync.
type Pipeline struct {
	stages []Stage
}

// NewPipeline constructs a Pipeline with the provided stages.
func NewPipeline(stages ...Stage) *Pipeline {
	return &Pipeline{stages: stages}
}

// AddStage appends a stage to the pipeline.
func (p *Pipeline) AddStage(s Stage) {
	p.stages = append(p.stages, s)
}

// Execute runs all stages in order for the given profile.
// It stops and returns the first error encountered, annotated with the stage name.
func (p *Pipeline) Execute(ctx context.Context, profileName string) error {
	for _, s := range p.stages {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("pipeline cancelled before stage %q: %w", s.Name, err)
		}
		if err := s.Run(ctx, profileName); err != nil {
			return fmt.Errorf("stage %q failed for profile %q: %w", s.Name, profileName, err)
		}
	}
	return nil
}

// Len returns the number of stages in the pipeline.
func (p *Pipeline) Len() int {
	return len(p.stages)
}
