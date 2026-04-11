// Package diff provides functionality to compare secret maps,
// typically between the current local .env file contents and
// secrets fetched from HashiCorp Vault.
//
// It classifies each key as Added, Updated, Removed, or Unchanged
// and exposes a Summary for reporting purposes.
package diff
