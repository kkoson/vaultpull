// Package sync orchestrates fetching secrets from Vault and writing them to
// a local .env file.
package sync

import "fmt"

// SecretsClient defines the interface for retrieving secrets from a remote
// secrets backend.
type SecretsClient interface {
	GetSecrets(path string) (map[string]string, error)
}

// EnvWriter defines the interface for persisting secrets to a local file.
type EnvWriter interface {
	Write(path string, secrets map[string]string) error
}

// Syncer coordinates pulling secrets from Vault and writing them to disk.
type Syncer struct {
	client SecretsClient
	writer EnvWriter
}

// New creates a new Syncer with the provided SecretsClient and EnvWriter.
func New(client SecretsClient, writer EnvWriter) *Syncer {
	return &Syncer{
		client: client,
		writer: writer,
	}
}

// Run fetches secrets from the given Vault path and writes them to outPath.
// It returns an error if either the fetch or the write operation fails.
func (s *Syncer) Run(vaultPath, outPath string) error {
	secrets, err := s.client.GetSecrets(vaultPath)
	if err != nil {
		return fmt.Errorf("sync: failed to fetch secrets from %q: %w", vaultPath, err)
	}

	if err := s.writer.Write(outPath, secrets); err != nil {
		return fmt.Errorf("sync: failed to write secrets to %q: %w", outPath, err)
	}

	return nil
}
