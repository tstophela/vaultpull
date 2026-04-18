package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSnapshotManager_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	m := NewSnapshotManager(dir)

	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "secret",
	}

	if err := m.Save("secret/myapp", secrets); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := m.Load("secret/myapp")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if snap == nil {
		t.Fatal("expected snapshot, got nil")
	}
	if snap.Path != "secret/myapp" {
		t.Errorf("path = %q, want %q", snap.Path, "secret/myapp")
	}
	if snap.Secrets["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST = %q, want localhost", snap.Secrets["DB_HOST"])
	}
	if snap.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestSnapshotManager_Load_Missing(t *testing.T) {
	dir := t.TempDir()
	m := NewSnapshotManager(dir)

	snap, err := m.Load("secret/nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap != nil {
		t.Error("expected nil for missing snapshot")
	}
}

func TestSnapshotManager_Save_CreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "snapshots")
	m := NewSnapshotManager(dir)

	if err := m.Save("kv/app", map[string]string{"KEY": "val"}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if _, err := os.Stat(dir); err != nil {
		t.Errorf("expected dir to exist: %v", err)
	}
}

func TestSnapshotManager_Filename_Sanitizes(t *testing.T) {
	dir := t.TempDir()
	m := NewSnapshotManager(dir)

	_ = m.Save("secret/app/prod", map[string]string{"X": "1"})

	matches, _ := filepath.Glob(filepath.Join(dir, "*.snap.json"))
	if len(matches) != 1 {
		t.Fatalf("expected 1 snapshot file, got %d", len(matches))
	}
}
