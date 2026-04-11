package rotate

import (
	"context"
	"fmt"
	"time"

	"github.com/your-org/vaultpull/internal/vault"
)

const defaultSecretIDTTL = 24 * time.Hour

// SecretIDResponse holds the new secret-id and its TTL.
type SecretIDResponse struct {
	SecretID string
	TTL      time.Duration
}

// generateSecretID calls the Vault API to generate a new AppRole secret-id.
func generateSecretID(ctx context.Context, client *vault.Client, roleID string) (string, time.Duration, error) {
	path := fmt.Sprintf("auth/approle/role/%s/secret-id", roleID)

	secret, err := client.Logical().WriteWithContext(ctx, path, nil)
	if err != nil {
		return "", 0, fmt.Errorf("generateSecretID: vault write failed: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return "", 0, fmt.Errorf("generateSecretID: empty response from vault")
	}

	secretID, ok := secret.Data["secret_id"].(string)
	if !ok || secretID == "" {
		return "", 0, fmt.Errorf("generateSecretID: missing secret_id in response")
	}

	ttl := defaultSecretIDTTL
	if ttlRaw, ok := secret.Data["secret_id_ttl"]; ok {
		if ttlSec, ok := ttlRaw.(float64); ok && ttlSec > 0 {
			ttl = time.Duration(ttlSec) * time.Second
		}
	}

	return secretID, ttl, nil
}
