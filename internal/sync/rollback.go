package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RollbackStore manages backup copies of .env files so that a failed sync
// can be reverted to the last known-good state.
type RollbackStore struct {
	dir string
}

// NewRollbackStore creates a RollbackStore that persists backups under dir.
// The directory is created if it does not exist.
func NewRollbackStore(dir string) (*RollbackStore, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("rollback: create dir: %w", err)
	}
	return &RollbackStore{dir: dir}, nil
}

// Backup copies the file at src into the store, keyed by profile name.
// If src does not exist the call is a no-op (nothing to back up).
func (r *RollbackStore) Backup(profile, src string) error {
	data, err := os.ReadFile(src)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("rollback: read source: %w", err)
	}
	dest := r.backupPath(profile)
	if err := os.WriteFile(dest, data, 0o600); err != nil {
		return fmt.Errorf("rollback: write backup: %w", err)
	}
	return nil
}

// Restore copies the stored backup for profile back to dest.
// Returns an error if no backup exists.
func (r *RollbackStore) Restore(profile, dest string) error {
	src := r.backupPath(profile)
	data, err := os.ReadFile(src)
	if os.IsNotExist(err) {
		return fmt.Errorf("rollback: no backup found for profile %q", profile)
	}
	if err != nil {
		return fmt.Errorf("rollback: read backup: %w", err)
	}
	if err := os.WriteFile(dest, data, 0o600); err != nil {
		return fmt.Errorf("rollback: restore file: %w", err)
	}
	return nil
}

// Clear removes the stored backup for profile.
func (r *RollbackStore) Clear(profile string) error {
	path := r.backupPath(profile)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("rollback: clear: %w", err)
	}
	return nil
}

// HasBackup reports whether a backup exists for profile.
func (r *RollbackStore) HasBackup(profile string) bool {
	_, err := os.Stat(r.backupPath(profile))
	return err == nil
}

func (r *RollbackStore) backupPath(profile string) string {
	safe := filepath.Base(profile) // prevent path traversal
	return filepath.Join(r.dir, safe+".bak")
}

// BackupTimestamp returns the modification time of the backup file, or the
// zero time if no backup exists.
func (r *RollbackStore) BackupTimestamp(profile string) time.Time {
	info, err := os.Stat(r.backupPath(profile))
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}
