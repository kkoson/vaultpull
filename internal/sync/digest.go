package sync

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// DigestConfig controls how secret digests are computed.
type DigestConfig struct {
	// Algorithm is the hash algorithm to use (currently only "sha256" supported).
	Algorithm string
}

// DefaultDigestConfig returns a DigestConfig with sensible defaults.
func DefaultDigestConfig() DigestConfig {
	return DigestConfig{
		Algorithm: "sha256",
	}
}

// Digester computes a stable digest over a map of secret key/value pairs.
// The digest is order-independent: keys are sorted before hashing so that
// two maps with identical contents always produce the same digest.
type Digester struct {
	cfg DigestConfig
}

// NewDigester creates a Digester using the provided config.
// If cfg is zero-valued, defaults are applied.
func NewDigester(cfg DigestConfig) *Digester {
	if cfg.Algorithm == "" {
		cfg = DefaultDigestConfig()
	}
	return &Digester{cfg: cfg}
}

// Compute returns a hex-encoded digest of the supplied secrets map.
// Keys are sorted lexicographically before hashing to ensure stability.
func (d *Digester) Compute(secrets map[string]string) string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, secrets[k])
	}

	sum := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(sum[:])
}

// Equal returns true when two secrets maps produce the same digest.
func (d *Digester) Equal(a, b map[string]string) bool {
	return d.Compute(a) == d.Compute(b)
}
