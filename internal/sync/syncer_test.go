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
	w := env.NewWriter(target, false)

	s := New(v, w, nil, nil, nil)
	if err := s.Sync(context.Background(), "secret", "app", target, nil, false); err != nil {
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
	w := env.NewWriter(target, true)

	s := New(v, w, nil, nil, nil)
	if err := s.Sync(context.Background(), "secret", "app", target, nil, false); err != nil {
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
	w := env.NewWriter(target, false)

	var buf bytes.Buffer
	logger := env.NewAuditLogger(&buf)

	s := New(v, w, logger, nil, nil)
	_ = s.Sync(context.Background(), "secret", "myapp", target, nil, false)

	output := buf.String()
	if !strings.Contains(output, "myapp") {
		t.Error("expected myapp in audit output")
	}
	if !strings.Contains(output, "ALPHA") {
		t.Error("expected ALPHA as added key in audit")
	}
}

func TestSyncer_Sync_DryRun(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, ".env")

	v := newMockVault(map[string]string{"FOO": "bar", "BAZ": "qux"})
	w := env.NewWriter(target, false)

	var out bytes.Buffer
	s := New(v, w, nil, nil, nil)
	if err := s.Sync(context.Background(), "secret", "app", target, &out, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify file was NOT created
	if _, err := os.Stat(target); err == nil {
		t.Error("expected file not to be created in dry-run mode")
	}

	// Verify dry-run output contains the expected markers
	output := out.String()
	if !strings.Contains(output, "[DRY-RUN]") {
		t.Error("expected [DRY-RUN] marker in output")
	}
	if !strings.Contains(output, "Would sync") {
		t.Error("expected 'Would sync' in output")
	}
	if !strings.Contains(output, "Added") {
		t.Error("expected 'Added' section in dry-run output")
	}
}

func TestSyncer_Sync_DryRun_WithUpdates(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, ".env")
	_ = os.WriteFile(target, []byte("EXISTING=old\nKEEP=yes\n"), 0600)

	v := newMockVault(map[string]string{"EXISTING": "new", "NEW_KEY": "added", "KEEP": "yes"})
	w := env.NewWriter(target, false)

	var out bytes.Buffer
	s := New(v, w, nil, nil, nil)
	if err := s.Sync(context.Background(), "secret", "app", target, &out, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Added") || !strings.Contains(output, "NEW_KEY") {
		t.Error("expected 'Added' section with NEW_KEY in dry-run output")
	}
	if !strings.Contains(output, "Updated") || !strings.Contains(output, "EXISTING") {
		t.Error("expected 'Updated' section with EXISTING in dry-run output")
	}
	if !strings.Contains(output, "Unchanged") || !strings.Contains(output, "KEEP") {
		t.Error("expected 'Unchanged' section with KEEP in dry-run output")
	}
}
