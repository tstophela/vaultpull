package env

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func fixedRollbackTime() func() time.Time {
	t := time.Date(2024, 6, 1, 12, 0, 0, 42, time.UTC)
	return func() time.Time { return t }
}

func newTestRollback(t *testing.T) (*RollbackManager, string) {
	t.Helper()
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	rm := NewRollbackManager(filepath.Join(dir, "rollbacks"), envFile)
	rm.now = fixedRollbackTime()
	return rm, envFile
}

func TestRollback_SaveAndList(t *testing.T) {
	rm, _ := newTestRollback(t)
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	rp, err := rm.Save("initial", secrets)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if rp.Label != "initial" {
		t.Errorf("expected label 'initial', got %q", rp.Label)
	}
	points, err := rm.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(points) != 1 {
		t.Fatalf("expected 1 point, got %d", len(points))
	}
	if points[0].ID != rp.ID {
		t.Errorf("ID mismatch: %s vs %s", points[0].ID, rp.ID)
	}
}

func TestRollback_List_Empty(t *testing.T) {
	rm, _ := newTestRollback(t)
	points, err := rm.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(points) != 0 {
		t.Errorf("expected 0 points, got %d", len(points))
	}
}

func TestRollback_Restore(t *testing.T) {
	rm, envFile := newTestRollback(t)
	secrets := map[string]string{"KEY": "value123"}
	rp, err := rm.Save("test", secrets)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	restored, err := rm.Restore(rp.ID)
	if err != nil {
		t.Fatalf("Restore: %v", err)
	}
	if restored.Label != "test" {
		t.Errorf("expected label 'test', got %q", restored.Label)
	}
	data, err := os.ReadFile(envFile)
	if err != nil {
		t.Fatalf("read env file: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty env file after restore")
	}
}

func TestRollback_Restore_Missing(t *testing.T) {
	rm, _ := newTestRollback(t)
	_, err := rm.Restore("nonexistent")
	if err == nil {
		t.Error("expected error restoring missing point")
	}
}

func TestRollback_List_NewestFirst(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	rm := NewRollbackManager(filepath.Join(dir, "rollbacks"), envFile)
	times := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
	}
	for i, ts := range times {
		captured := ts
		rm.now = func() time.Time { return captured }
		_ = i
		_, _ = rm.Save("label", map[string]string{"K": "v"})
	}
	points, _ := rm.List()
	if len(points) != 3 {
		t.Fatalf("expected 3 points, got %d", len(points))
	}
	if !points[0].CreatedAt.After(points[1].CreatedAt) {
		t.Error("expected newest-first ordering")
	}
}
