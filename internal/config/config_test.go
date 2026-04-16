package config

import (
	"os"
	"testing"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	t.Setenv(key, value)
}

func TestLoad_MissingToken(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")
	os.Unsetenv("VAULT_PATH")

	_, err := Load("")
	if err == nil {
		t.Fatal("expected error for missing VAULT_TOKEN, got nil")
	}
}

func TestLoad_MissingPath(t *testing.T) {
	setEnv(t, "VAULT_TOKEN", "test-token")
	os.Unsetenv("VAULT_PATH")

	_, err := Load("")
	if err == nil {
		t.Fatal("expected error for missing VAULT_PATH, got nil")
	}
}

func TestLoad_Defaults(t *testing.T) {
	setEnv(t, "VAULT_TOKEN", "test-token")
	setEnv(t, "VAULT_PATH", "secret/myapp")
	os.Unsetenv("VAULT_ADDR")
	os.Unsetenv("VAULTPULL_OUTPUT")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("expected default VaultAddr, got %q", cfg.VaultAddr)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("expected default OutputFile '.env', got %q", cfg.OutputFile)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	setEnv(t, "VAULT_TOKEN", "s.abc123")
	setEnv(t, "VAULT_PATH", "secret/data/prod")
	setEnv(t, "VAULT_ADDR", "https://vault.example.com")
	setEnv(t, "VAULTPULL_OUTPUT", ".env.prod")
	setEnv(t, "VAULT_NAMESPACE", "admin")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.VaultToken != "s.abc123" {
		t.Errorf("unexpected VaultToken: %q", cfg.VaultToken)
	}
	if cfg.Namespace != "admin" {
		t.Errorf("unexpected Namespace: %q", cfg.Namespace)
	}
	if cfg.OutputFile != ".env.prod" {
		t.Errorf("unexpected OutputFile: %q", cfg.OutputFile)
	}
}
