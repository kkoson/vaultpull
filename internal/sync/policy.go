package sync

import "time"

// RetentionPolicy defines how long synced secret data should be considered
// valid before a forced re-sync is required, independent of drift detection.
type RetentionPolicy struct {
	// MaxAge is the maximum age of a successful sync before it is considered stale.
	MaxAge time.Duration

	// EnforceOnFailure, when true, blocks writes if the last successful sync
	// exceeds MaxAge, preventing stale data propagation.
	EnforceOnFailure bool
}

// DefaultRetentionPolicy returns a RetentionPolicy with sensible defaults.
func DefaultRetentionPolicy() RetentionPolicy {
	return RetentionPolicy{
		MaxAge:           24 * time.Hour,
		EnforceOnFailure: false,
	}
}

// PolicyEnforcer checks whether a profile sync is permitted under the
// configured RetentionPolicy.
type PolicyEnforcer struct {
	policy    RetentionPolicy
	lastSyncs map[string]time.Time
}

// NewPolicyEnforcer creates a PolicyEnforcer with the given RetentionPolicy.
// If policy.MaxAge is zero, DefaultRetentionPolicy values are used.
func NewPolicyEnforcer(policy RetentionPolicy) *PolicyEnforcer {
	if policy.MaxAge <= 0 {
		policy.MaxAge = DefaultRetentionPolicy().MaxAge
	}
	return &PolicyEnforcer{
		policy:    policy,
		lastSyncs: make(map[string]time.Time),
	}
}

// Record marks a successful sync for the given profile at the current time.
func (e *PolicyEnforcer) Record(profile string) {
	e.lastSyncs[profile] = time.Now()
}

// Allow returns true if the profile is permitted to sync under the policy.
// A profile is allowed if it has never been synced or its last sync is within MaxAge.
func (e *PolicyEnforcer) Allow(profile string) bool {
	t, ok := e.lastSyncs[profile]
	if !ok {
		return true
	}
	return time.Since(t) <= e.policy.MaxAge
}

// Age returns the duration since the last successful sync for the given profile.
// Returns -1 if the profile has never been synced.
func (e *PolicyEnforcer) Age(profile string) time.Duration {
	t, ok := e.lastSyncs[profile]
	if !ok {
		return -1
	}
	return time.Since(t)
}
