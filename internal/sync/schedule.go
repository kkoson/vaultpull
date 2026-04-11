package sync

import (
	"context"
	"fmt"
	"time"
)

// Schedule defines a recurring sync interval for a set of profiles.
type Schedule struct {
	// Interval is how often the sync should run.
	Interval time.Duration
	// Profiles lists profile names to sync; empty means all profiles.
	Profiles []string
	// runner executes syncs on each tick.
	runner *Runner
	// opts are applied to every sync run.
	opts Options
}

// NewSchedule creates a Schedule that will sync using the given Runner.
func NewSchedule(r *Runner, interval time.Duration, opts Options, profiles ...string) (*Schedule, error) {
	if r == nil {
		return nil, fmt.Errorf("schedule: runner must not be nil")
	}
	if interval <= 0 {
		return nil, fmt.Errorf("schedule: interval must be positive, got %s", interval)
	}
	return &Schedule{
		Interval: interval,
		Profiles: profiles,
		runner:   r,
		opts:     opts,
	}, nil
}

// Start blocks and runs syncs on every tick until ctx is cancelled.
// It performs an initial sync immediately before waiting for the first tick.
func (s *Schedule) Start(ctx context.Context) error {
	if err := s.tick(ctx); err != nil {
		return err
	}
	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.tick(ctx); err != nil {
				return err
			}
		}
	}
}

func (s *Schedule) tick(ctx context.Context) error {
	if len(s.Profiles) == 0 {
		return s.runner.RunAll(ctx, s.opts)
	}
	for _, name := range s.Profiles {
		if err := s.runner.RunProfile(ctx, name, s.opts); err != nil {
			return fmt.Errorf("schedule: profile %q: %w", name, err)
		}
	}
	return nil
}
