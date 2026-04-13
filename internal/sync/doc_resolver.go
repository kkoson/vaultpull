// Package sync provides primitives for synchronising Vault secrets into
// local .env files.
//
// # Resolver
//
// Resolver translates a user-supplied profile name (or prefix) into one or
// more concrete config.Profile values that can be passed to a Runner or
// Syncer.
//
// Two resolution modes are supported:
//
//   - ResolveModeExact  – the supplied name must match a profile name exactly.
//   - ResolveModePrefix – any profile whose name starts with the supplied
//     string is included in the result set.
//
// Example:
//
//	resolver := sync.NewResolver(cfg, sync.ResolveModePrefix)
//	profiles, err := resolver.Resolve("prod")
//	// profiles contains every profile whose name begins with "prod"
package sync
