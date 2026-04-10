package vault

import (
	"testing"
)

func TestNew_NoAuthMethod(t *testing.T) {
	_, err := New(Config{
		Address: "http://127.0.0.1:8200",
	})
	if err == nil {
		t.Fatal("expected error when no auth method is provided")
	}
}

func TestNew_WithToken(t *testing.T) {
	client, err := New(Config{
		Address: "http://127.0.0.1:8200",
		Token:   "root",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.v == nil {
		t.Fatal("expected underlying vault client to be set")
	}
}

func TestNew_AppRoleMissingSecretID(t *testing.T) {
	// RoleID provided but no SecretID — should fall through to error.
	_, err := New(Config{
		Address: "http://127.0.0.1:8200",
		RoleID:  "my-role",
	})
	if err == nil {
		t.Fatal("expected error when SecretID is missing")
	}
}
