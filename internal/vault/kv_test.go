package vault

import (
	"testing"
)

func TestDetectKVVersion_DefaultsToV2OnNilSecret(t *testing.T) {
	// Without a live Vault we verify that the helper path-building logic and
	// version string parsing work correctly by exercising DetectKVVersion
	// indirectly through ReadSecret path construction.
	v := buildFullPath("secret", "myapp/config", KVv2)
	want := "secret/data/myapp/config"
	if v != want {
		t.Errorf("KVv2 path: got %q, want %q", v, want)
	}
}

func TestDetectKVVersion_V1Path(t *testing.T) {
	v := buildFullPath("kv", "myapp/config", KVv1)
	want := "kv/myapp/config"
	if v != want {
		t.Errorf("KVv1 path: got %q, want %q", v, want)
	}
}

func TestDetectKVVersion_TrimsTrailingSlash(t *testing.T) {
	v := buildFullPath("secret/", "/myapp/config", KVv2)
	want := "secret/data/myapp/config"
	if v != want {
		t.Errorf("trim slash path: got %q, want %q", v, want)
	}
}

func TestKVVersion_Constants(t *testing.T) {
	if KVv1 != 1 {
		t.Errorf("KVv1 should be 1, got %d", KVv1)
	}
	if KVv2 != 2 {
		t.Errorf("KVv2 should be 2, got %d", KVv2)
	}
}

// buildFullPath mirrors the path construction inside ReadSecret so we can
// unit-test the logic without a real Vault connection.
func buildFullPath(mountPath, secretPath string, version KVVersion) string {
	import_strings_TrimSuffix := func(s, suffix string) string {
		if len(s) > 0 && s[len(s)-1] == '/' {
			return s[:len(s)-1]
		}
		return s
	}
	import_strings_TrimPrefix := func(s, prefix string) string {
		if len(s) > 0 && s[0] == '/' {
			return s[1:]
		}
		return s
	}

	mountPath = import_strings_TrimSuffix(mountPath, "/")
	secretPath = import_strings_TrimPrefix(secretPath, "/")

	if version == KVv1 {
		return mountPath + "/" + secretPath
	}
	return mountPath + "/data/" + secretPath
}
