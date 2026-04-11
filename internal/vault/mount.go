package vault

import (
	"context"
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// MountInfo holds metadata about a Vault secret mount.
type MountInfo struct {
	Path    string
	Type    string
	Version string
}

// ListMounts returns all KV mounts accessible to the client.
func ListMounts(ctx context.Context, client *vaultapi.Client) ([]MountInfo, error) {
	mounts, err := client.Sys().ListMountsWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing mounts: %w", err)
	}

	var result []MountInfo
	for path, mount := range mounts {
		if mount.Type != "kv" {
			continue
		}
		version := mount.Options["version"]
		if version == "" {
			version = "1"
		}
		result = append(result, MountInfo{
			Path:    strings.TrimSuffix(path, "/"),
			Type:    mount.Type,
			Version: version,
		})
	}
	return result, nil
}

// FindMount returns the MountInfo for the mount that best matches the given
// secret path, or an error if no KV mount covers it.
func FindMount(ctx context.Context, client *vaultapi.Client, secretPath string) (*MountInfo, error) {
	mounts, err := ListMounts(ctx, client)
	if err != nil {
		return nil, err
	}

	var best *MountInfo
	for i := range mounts {
		m := &mounts[i]
		prefix := m.Path + "/"
		if strings.HasPrefix(secretPath, prefix) || secretPath == m.Path {
			if best == nil || len(m.Path) > len(best.Path) {
				best = m
			}
		}
	}
	if best == nil {
		return nil, fmt.Errorf("no KV mount found for path %q", secretPath)
	}
	return best, nil
}
