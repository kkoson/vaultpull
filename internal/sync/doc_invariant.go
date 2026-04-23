// Package sync provides the InvariantChecker for asserting postconditions
// on synced secret maps.
//
// An invariant is a named rule applied after secrets are fetched from Vault
// and before they are written to the .env file. Rules can enforce required
// keys, forbid sensitive keys from leaking into certain profiles, or apply
// arbitrary custom validation logic.
//
// Example:
//
//	checker := sync.NewInvariantChecker()
//	checker.Register(sync.RequireKeys("DB_URL", "API_KEY"))
//	checker.Register(sync.ForbidKeys("DEBUG"))
//
//	if err := checker.Check(secrets); err != nil {
//		log.Fatalf("invariant violation: %v", err)
//	}
package sync
