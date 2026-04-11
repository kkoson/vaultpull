package sync

import (
	"context"
	"fmt"

	"github.com/yourusername/vaultpull/internal/audit"
	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/diff"
	"github.com/yourusername/vaultpull/internal/env"
)

// SecretReader fetches key/value pairs from a remote secrets store.
type SecretReader interface {
	GetSecrets(ctx context.Context, mountPath, secretPath string) (map[string]string, error)
}

// Syncer orchestrates pulling secrets and writing them to a local env file.
type Syncer struct {
	client  SecretReader
	writer  *env.Merger
	printer *diff.Printer
	audit   *audit.Logger
	opts    Options
}

// New creates a Syncer with the provided dependencies.
func New(client SecretReader, writer *env.Merger, printer *diff.Printer, auditLog *audit.Logger, opts Options) *Syncer {
	return &Syncer{
		client:  client,
		writer:  writer,
		printer: printer,
		audit:   auditLog,
		opts:    opts,
	}
}

// Run pulls secrets for the given profile and merges them into the output file.
func (s *Syncer) Run(ctx context.Context, profile config.Profile) error {
	secrets, err := s.client.GetSecrets(ctx, profile.MountPath, profile.VaultPath)
	if err != nil {
		if s.audit != nil {
			_ = s.audit.Write(audit.Entry{
				Profile:    profile.Name,
				VaultPath:  profile.VaultPath,
				OutputFile: profile.OutputFile,
				DryRun:     s.opts.DryRun,
				Error:      err.Error(),
			})
		}
		return fmt.Errorf("fetch secrets: %w", err)
	}

	changes := diff.Compare(nil, secrets)
	if s.printer != nil && s.opts.Verbose {
		s.printer.Print(changes)
	}

	if !s.opts.DryRun {
		if err := s.writer.Merge(profile.OutputFile, secrets, s.opts.OverwriteExisting); err != nil {
			return fmt.Errorf("write env file: %w", err)
		}
	}

	if s.audit != nil {
		counts := countChanges(changes)
		_ = s.audit.Write(audit.Entry{
			Profile:    profile.Name,
			VaultPath:  profile.VaultPath,
			OutputFile: profile.OutputFile,
			Added:      counts.added,
			Updated:    counts.updated,
			Removed:    counts.removed,
			Unchanged:  counts.unchanged,
			DryRun:     s.opts.DryRun,
		})
	}
	return nil
}

type changeCounts struct{ added, updated, removed, unchanged int }

func countChanges(changes []diff.Change) changeCounts {
	var c changeCounts
	for _, ch := range changes {
		switch ch.Type {
		case diff.Added:
			c.added++
		case diff.Updated:
			c.updated++
		case diff.Removed:
			c.removed++
		case diff.Unchanged:
			c.unchanged++
		}
	}
	return c
}
