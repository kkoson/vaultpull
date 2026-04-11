package vault

import (
	"context"
	"fmt"
	"strings"
)

// KVVersion represents the KV secrets engine version.
type KVVersion int

const (
	KVv1 KVVersion = 1
	KVv2 KVVersion = 2
)

// DetectKVVersion probes the mount path to determine whether it is a KV v1 or
// KV v2 secrets engine. It returns KVv2 by default when detection is
// ambiguous.
func DetectKVVersion(ctx context.Context, c *Client, mountPath string) (KVVersion, error) {
	mountPath = strings.TrimSuffix(mountPath, "/")

	// Vault returns mount metadata under sys/mounts/<mount>.
	path := fmt.Sprintf("sys/mounts/%s", mountPath)
	secret, err := c.vault.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return KVv2, fmt.Errorf("detect kv version: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return KVv2, nil
	}

	options, ok := secret.Data["options"].(map[string]interface{})
	if !ok {
		return KVv2, nil
	}

	version, ok := options["version"].(string)
	if !ok {
		return KVv2, nil
	}

	if version == "1" {
		return KVv1, nil
	}
	return KVv2, nil
}

// ReadSecret reads a secret from the appropriate KV path depending on the
// engine version. For KV v2 it injects the "data" segment automatically.
func ReadSecret(ctx context.Context, c *Client, mountPath, secretPath string, version KVVersion) (map[string]interface{}, error) {
	mountPath = strings.TrimSuffix(mountPath, "/")
	secretPath = strings.TrimPrefix(secretPath, "/")

	var fullPath string
	switch version {
	case KVv1:
		fullPath = fmt.Sprintf("%s/%s", mountPath, secretPath)
	default:
		fullPath = fmt.Sprintf("%s/data/%s", mountPath, secretPath)
	}

	secret, err := c.vault.Logical().ReadWithContext(ctx, fullPath)
	if err != nil {
		return nil, fmt.Errorf("read secret %q: %w", fullPath, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret not found at path %q", fullPath)
	}

	return flattenData(secret.Data), nil
}
