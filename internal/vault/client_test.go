package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/vaultpull/internal/vault"
)

func newMockVault(t *testing.T, path string, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(body); err != nil {
			t.Fatalf("mock encode: %v", err)
		}
	}))
}

func TestReadSecrets_KVv2(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]interface{}{
				"DB_PASS": "secret",
				"API_KEY": "abc123",
			},
		},
	}
	srv := newMockVault(t, "/v1/secret/data/myapp", payload)
	defer srv.Close()

	client, err := vault.NewClient(vault.Config{
		Address: srv.URL,
		Token:   "test-token",
		Timeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	secrets, err := client.ReadSecrets("secret/data/myapp")
	if err != nil {
		t.Fatalf("ReadSecrets: %v", err)
	}

	if secrets["DB_PASS"] != "secret" {
		t.Errorf("expected DB_PASS=secret, got %q", secrets["DB_PASS"])
	}
	if secrets["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", secrets["API_KEY"])
	}
}

func TestReadSecrets_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client, err := vault.NewClient(vault.Config{
		Address: srv.URL,
		Token:   "test-token",
		Timeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = client.ReadSecrets("secret/data/missing")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}
