package sync_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/sync"
	"github.com/your-org/vaultpull/internal/vault"
)

func newMockVault(t *testing.T, response string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
}

func TestSyncer_Sync_Success(t *testing.T) {
	body := `{"data":{"data":{"API_KEY":"abc123","DB_PASS":"secret"}}}`
	srv := newMockVault(t, body)
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, ".env")
	writer := env.NewWriter()

	s := sync.New(client, writer)
	result, err := s.Sync("secret", "myapp", outFile, false)
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}

	if result.Keyssynced != 2 {
		t.Errorf("expected 2 keys synced, got %d", result.Keyssynced)
	}
	if result.BackedUp {
		t.Error("expected no backup")
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}
	if len(data) == 0 {
		t.Error("output file is empty")
	}
}

func TestSyncer_Sync_WithBackup(t *testing.T) {
	body := `{"data":{"data":{"FOO":"bar"}}}`
	srv := newMockVault(t, body)
	defer srv.Close()

	client, _ := vault.NewClient(srv.URL, "tok")
	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, ".env")
	_ = os.WriteFile(outFile, []byte("OLD=value\n"), 0600)

	s := sync.New(client, env.NewWriter())
	result, err := s.Sync("secret", "app", outFile, true)
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if !result.BackedUp {
		t.Error("expected backup to be created")
	}
}
