package sync

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/vaultpull/internal/env"
)

type mockVault struct {
	secrets map[string]string
	err     error
}

func (m *mockVault) ReadSecrets(_ context.Context, _ string) (map[string]string, error) {
	return m.secrets, m.err
}

func newMockVault(secrets map[string]string) *mockVault {
	return &mockVault{secrets: secrets}
}

func TestSyncer_Sync_Success(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, ".env")

	v := newMockVault(map[string]string{"FOO": "bar", "BAZ": "qux"})
	w := env.NewWriter(false)
	r := env.NewReader()
	logger := env.NewAuditLogger(nil)

	s := New(v, w, r, logger)
	if err := s.Sync(context.Background(), "secret/app", target); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}
	if !strings.Contains(string(data), "FOO=") {
		t.Error("expected FOO in output")
	}
}

func TestSyncer_Sync_WithBackup(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, ".env")
	_ = os.WriteFile(target, []byte("EXISTING=yes\n"), 0600)

	v := newMockVault(map[string]string{"NEW_KEY": "value"})
	w := env.NewWriter(true)
	r := env.NewReader()

	s := New(v, w, r, nil)
	if err := s.Sync(context.Background(), "secret/app", target); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := os.ReadDir(dir)
	var hasBackup bool
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".bak") {
			hasBackup = true
		}
	}
	if !hasBackup {
		t.Error("expected backup file")
	}
}

func TestSyncer_Sync_AuditOutput(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, ".env")

	v := newMockVault(map[string]string{"ALPHA": "1"})
	w := env.NewWriter(false)
	r := env.NewReader()

	var buf bytes.Buffer
	logger := env.NewAuditLogger(&buf)

	s := New(v, w, r, logger)
	_ = s.Sync(context.Background(), "secret/myapp", target)

	if !strings.Contains(buf.String(), "secret/myapp") {
		t.Error("expected audit entry for path")
	}
	if !strings.Contains(buf.String(), "+ ALPHA") {
		t.Error("expected ALPHA as added key in audit")
	}
}
