package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/vaultpull/internal/vault"
)

// newMockVault creates a test HTTP server that responds to all requests with
// the given JSON body. The server is not path-specific; callers should verify
// path handling at a higher level if needed.
func newMockVault(t *testing.T, path string, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(body); err != nil {
			t.Fatalf("mock encode: %v", err)
		}
	}))
}

// newVaultClient is a test helper that creates a vault.Client pointed at the
// given server URL, failing the test immediately on error.
func newVaultClient(t *testing.T, serverURL string) *vault.Client {
	t.Helper()
	client, err := vault.NewClient(vault.Config{
		Address: serverURL,
		Token:   "test-token",
		Timeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
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

	client := newVaultClient(t, srv.URL)

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

	client := newVaultClient(t, srv.URL)

	_, err := client.ReadSecrets("secret/data/missing")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}
