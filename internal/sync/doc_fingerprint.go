// Package sync provides synchronisation primitives for vaultpull.
//
// # Fingerprint
//
// FingerprintStore tracks a SHA-256 digest of the secret map for each
// profile. It allows the sync pipeline to skip writing .env files when
// the secrets fetched from Vault are identical to the previous run,
// reducing unnecessary disk I/O and downstream change noise.
//
// Usage:
//
//	fs := sync.NewFingerprintStore()
//
//	if fs.Changed(profile.Name, secrets) {
//		// write .env and record new fingerprint
//		fs.Record(profile.Name, secrets)
//	}
package sync
