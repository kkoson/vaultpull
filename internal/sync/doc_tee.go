// Package sync provides synchronisation primitives for vaultpull.
//
// # Tee Writer
//
// TeeWriter fans out a secret map write to multiple env.Writer targets.
// It is useful when the same set of secrets must be written to more than
// one output file (e.g. a shared .env and a service-specific override).
//
//	primary, _ := env.NewWriter(".env")
//	secondary, _ := env.NewWriter(".env.local")
//	tee := sync.NewTeeWriter(primary, secondary)
//	err := tee.Write(secrets)
//
// All writers are attempted; the first error encountered is returned.
package sync
