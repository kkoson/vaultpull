package sync

import (
	"context"
	"fmt"

	"github.com/yourusername/vaultpull/internal/audit"
	"github.com/yourusername/vaultpull/internal/config"
)

// Runner orchestrates syncing secrets for one or more profiles.
type Runner struct {
	cfg     *config.Config
	newSync func(profile *config.Profile, opts Options) (*Syncer, error)
	logger  *audit.Logger
}

// NewRunner creates a Runner backed by the given config and audit logger.
func NewRunner(cfg *config.Config, logger *audit.Logger) *Runner {
	return &Runner{
		cfg:     cfg,
		newSync: defaultNewSyncer,
		logger:  logger,
	}
}

// RunProfile syncs secrets for a single named profile.
func (r *Runner) RunProfile(ctx context.Context, name string, opts Options) error {
	profile, err := r.cfg.GetProfile(name)
	if err != nil {
		return fmt.Errorf("runner: profile %q not found: %w", name, err)
	}

	s, err := r.newSync(profile, opts)
	if err != nil {
		return fmt.Errorf("runner: failed to build syncer for profile %q: %w", name, err)
	}

	if err := s.Run(ctx); err != nil {
		r.logger.Write(audit.Entry{Profile: name, Error: err})
		return fmt.Errorf("runner: sync failed for profile %q: %w", name, err)
	}

	r.logger.Write(audit.Entry{Profile: name})
	return nil
}

// RunAll syncs secrets for every profile defined in the config.
func (r *Runner) RunAll(ctx context.Context, opts Options) error {
	var errs []error
	for _, p := range r.cfg.Profiles {
		if err := r.RunProfile(ctx, p.Name, opts); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("runner: %d profile(s) failed: %v", len(errs), errs)
	}
	return nil
}

func defaultNewSyncer(profile *config.Profile, opts Options) (*Syncer, error) {
	return New(profile, opts)
}
