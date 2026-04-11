// Package vault provides utilities for interacting with HashiCorp Vault,
// including authentication, KV secret reading, version detection, and
// mount discovery.
//
// The mount sub-feature (mount.go) allows callers to list all KV-type
// secret engine mounts and resolve which mount owns a given secret path.
// This is used during sync to automatically determine the KV version
// without requiring explicit configuration.
package vault
