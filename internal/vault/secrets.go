package vault

import (
	"fmt"
	"strings"
)

// ReadSecrets reads a KV secret at the given path and returns a map of
// key/value string pairs. Supports both KV v1 and KV v2 mounts.
func (c *Client) ReadSecrets(path string) (map[string]string, error) {
	secret, err := c.v.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("vault: read %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("vault: no secret found at %q", path)
	}

	data := secret.Data

	// KV v2 wraps values under a "data" key.
	if nested, ok := data["data"]; ok {
		if nestedMap, ok := nested.(map[string]interface{}); ok {
			data = nestedMap
		}
	}

	return flattenData(data), nil
}

// flattenData converts map[string]interface{} to map[string]string.
func flattenData(raw map[string]interface{}) map[string]string {
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		switch val := v.(type) {
		case string:
			out[k] = val
		case nil:
			out[k] = ""
		default:
			out[k] = strings.TrimSpace(fmt.Sprintf("%v", val))
		}
	}
	return out
}
