package sync

import "fmt"

// HookEvent represents the lifecycle event that triggered a hook.
type HookEvent string

const (
	// HookPreSync fires before secrets are fetched from Vault.
	HookPreSync HookEvent = "pre_sync"
	// HookPostSync fires after secrets have been written to the output file.
	HookPostSync HookEvent = "post_sync"
)

// HookFunc is a callback invoked at a specific point in the sync lifecycle.
// profile is the name of the profile being synced.
// event identifies when the hook is firing.
// err is non-nil when the hook fires after a failed operation (post_sync only).
type HookFunc func(profile string, event HookEvent, err error) error

// Hooks holds optional callbacks for sync lifecycle events.
type Hooks struct {
	PreSync  HookFunc
	PostSync HookFunc
}

// runPreSync invokes the PreSync hook if one is registered.
// It wraps any returned error with context about the profile and event.
func (h *Hooks) runPreSync(profile string) error {
	if h == nil || h.PreSync == nil {
		return nil
	}
	if err := h.PreSync(profile, HookPreSync, nil); err != nil {
		return fmt.Errorf("pre_sync hook failed for profile %q: %w", profile, err)
	}
	return nil
}

// runPostSync invokes the PostSync hook if one is registered.
// syncErr is the error (if any) produced by the sync operation itself.
func (h *Hooks) runPostSync(profile string, syncErr error) error {
	if h == nil || h.PostSync == nil {
		return nil
	}
	if err := h.PostSync(profile, HookPostSync, syncErr); err != nil {
		return fmt.Errorf("post_sync hook failed for profile %q: %w", profile, err)
	}
	return nil
}
