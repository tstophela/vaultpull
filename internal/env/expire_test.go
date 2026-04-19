package env

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExpiryManager_SetAndGet(t *testing.T) {
	dir := t.TempDir()
	m := NewExpiryManager(dir)

	err := m.Set("secret/myapp", 10*time.Minute)
	if err != nil {
		t.Fatalf("Set: %v", err)
	}

	rec, err := m.Get("secret/myapp")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.Path != "secret/myapp" {
		t.Errorf("expected path secret/myapp, got %s", rec.Path)
	}
	if rec.IsExpired() {
		t.Error("expected record to not be expired")
	}
}

func TestExpiryManager_IsExpired(t *testing.T) {
	dir := t.TempDir()
	m := NewExpiryManager(dir)

	err := m.Set("secret/old", -1*time.Second)
	if err != nil {
		t.Fatalf("Set: %v", err)
	}

	rec, err := m.Get("secret/old")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !rec.IsExpired() {
		t.Error("expected record to be expired")
	}
}

func TestExpiryManager_Get_Missing(t *testing.T) {
	dir := t.TempDir()
	m := NewExpiryManager(dir)

	_, err := m.Get("secret/nonexistent")
	if err == nil {
		t.Error("expected error for missing record")
	}
}

func TestExpiryManager_Delete(t *testing.T) {
	dir := t.TempDir()
	m := NewExpiryManager(dir)

	_ = m.Set("secret/tmp", time.Minute)
	if err := m.Delete("secret/tmp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	files, _ := filepath.Glob(filepath.Join(dir, "*.expiry.json"))
	if len(files) != 0 {
		t.Errorf("expected no files after delete, found %d", len(files))
	}
}

func TestExpiryManager_Delete_NonExistent(t *testing.T) {
	dir := t.TempDir()
	m := NewExpiryManager(dir)

	if err := m.Delete("secret/ghost"); err != nil {
		t.Errorf("expected no error deleting non-existent, got %v", err)
	}
}

func TestExpiryManager_CreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "expiry", "nested")
	m := NewExpiryManager(dir)

	if err := m.Set("secret/x", time.Minute); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("expected dir to be created: %v", err)
	}
}
