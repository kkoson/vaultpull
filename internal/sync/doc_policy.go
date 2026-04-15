// Package sync provides the RetentionPolicy and PolicyEnforcer types for
// controlling how long synced secrets remain valid before a forced re-sync
// is required.
//
// # Overview
//
// A RetentionPolicy defines a MaxAge duration. Once a profile's last
// successful sync exceeds MaxAge, the PolicyEnforcer will deny further
// syncs until the operator manually refreshes or the policy is relaxed.
//
// # Usage
//
//	enforcer := sync.NewPolicyEnforcer(sync.RetentionPolicy{
//		MaxAge:           12 * time.Hour,
//		EnforceOnFailure: true,
//	})
//
//	stage := sync.WithRetentionPolicy(enforcer)(innerStage)
//
// The middleware records a successful sync automatically, so the caller
// does not need to call enforcer.Record manually.
package sync
