// Package vault provides a thin wrapper around the HashiCorp Vault HTTP API
// client for use within vaultpull.
//
// It supports two authentication methods:
//   - Token-based: set Config.Token to a valid Vault token.
//   - AppRole:     set both Config.RoleID and Config.SecretID; the package
//     performs the login exchange and stores the resulting token.
//
// Secret retrieval is handled by [Client.ReadSecrets], which transparently
// supports both KV v1 and KV v2 secret engines by detecting the presence of
// the nested "data" envelope returned by v2 mounts.
package vault
