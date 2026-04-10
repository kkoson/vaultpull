package vault

import (
	"testing"
)

func TestFlattenData_StringValues(t *testing.T) {
	input := map[string]interface{}{
		"KEY_A": "value_a",
		"KEY_B": "value_b",
	}
	got := flattenData(input)
	if got["KEY_A"] != "value_a" {
		t.Errorf("expected value_a, got %q", got["KEY_A"])
	}
	if got["KEY_B"] != "value_b" {
		t.Errorf("expected value_b, got %q", got["KEY_B"])
	}
}

func TestFlattenData_NilValue(t *testing.T) {
	input := map[string]interface{}{
		"KEY_NULL": nil,
	}
	got := flattenData(input)
	if got["KEY_NULL"] != "" {
		t.Errorf("expected empty string for nil, got %q", got["KEY_NULL"])
	}
}

func TestFlattenData_NonStringValue(t *testing.T) {
	input := map[string]interface{}{
		"KEY_INT": 42,
		"KEY_BOOL": true,
	}
	got := flattenData(input)
	if got["KEY_INT"] != "42" {
		t.Errorf("expected \"42\", got %q", got["KEY_INT"])
	}
	if got["KEY_BOOL"] != "true" {
		t.Errorf("expected \"true\", got %q", got["KEY_BOOL"])
	}
}

func TestFlattenData_KVv2Nested(t *testing.T) {
	// Simulate what ReadSecrets does with KV v2 data before calling flattenData.
	raw := map[string]interface{}{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	got := flattenData(raw)
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
}

func TestFlattenData_Empty(t *testing.T) {
	got := flattenData(map[string]interface{}{})
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}
