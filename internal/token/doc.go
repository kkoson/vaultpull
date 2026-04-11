// Package token manages caching of HashiCorp Vault authentication tokens.
//
// It persists tokens to a JSON file under the OS user cache directory
// (or a caller-supplied path) and handles expiry checks automatically.
// Use NewCache to obtain a Cache, then Save/Load/Clear as needed.
package token
