package sync_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yourusername/vaultpull/internal/sync"
)

type mockSecretsClient struct {
	secrets map[string]string
	err     error
}

func (m *mockSecretsClient) GetSecrets(path string) (map[string]string, error) {
	return m.secrets, m.err
}

type mockEnvWriter struct {
	written map[string]string
	err     error
}

func (m *mockEnvWriter) Write(path string, secrets map[string]string) error {
	m.written = secrets
	return m.err
}

func TestSyncer_Run_Success(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, ".env")

	client := &mockSecretsClient{
		secrets: map[string]string{
			"DB_HOST": "localhost",
			"DB_PORT": "5432",
		},
	}
	writer := &mockEnvWriter{}

	s := sync.New(client, writer)
	err := s.Run("secret/myapp", outPath)

	require.NoError(t, err)
	assert.Equal(t, "localhost", writer.written["DB_HOST"])
	assert.Equal(t, "5432", writer.written["DB_PORT"])
}

func TestSyncer_Run_ClientError(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, ".env")

	client := &mockSecretsClient{
		err: os.ErrNotExist,
	}
	writer := &mockEnvWriter{}

	s := sync.New(client, writer)
	err := s.Run("secret/missing", outPath)

	require.Error(t, err)
	assert.Nil(t, writer.written)
}

func TestSyncer_Run_WriterError(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, ".env")

	client := &mockSecretsClient{
		secrets: map[string]string{"KEY": "value"},
	}
	writer := &mockEnvWriter{
		err: os.ErrPermission,
	}

	s := sync.New(client, writer)
	err := s.Run("secret/myapp", outPath)

	require.Error(t, err)
}
