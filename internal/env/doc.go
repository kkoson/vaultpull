// Package env provides utilities for writing Vault secrets into
// local .env files.
//
// The Writer type serialises a flat map[string]string of key/value
// pairs into the standard KEY=VALUE format understood by tools such
// as direnv, docker-compose, and the godotenv library.
//
// Values that contain whitespace, hash characters, or quotes are
// automatically wrapped in double-quotes so that consumers can parse
// the file unambiguously.
package env
