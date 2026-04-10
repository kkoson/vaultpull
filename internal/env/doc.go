// Package env provides utilities for reading and writing .env files.
//
// Writer creates or overwrites a .env file from a map of secret key-value
// pairs sourced from HashiCorp Vault. Reader parses an existing .env file
// back into a map, enabling diff and merge workflows before writing.
package env
