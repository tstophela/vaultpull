package env

import (
	"testing"
	"os"
)

func newTestPinManager(t *testing.T) (*PinManager, string) {
	t.Helper()
	dir := t.TempDir()
	return NewPinManager(dir), dir
}

func TestPinManager_PinAndLoad(t *testing.T) {
	pm, _ := newTestPinManager(t)
	if err := pm.Pin("dev", "DB_PASS", 3, "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pins, err := pm.Load("dev")
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	e, ok := pins["DB_PASS"]
	if !ok {
		t.Fatal("expected pin for DB_PASS")
	}
	if e.Version != 3 || e.PinnedBy != "alice" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestPinManager_Unpin(t *testing.T) {
	pm, _ := newTestPinManager(t)
	_ = pm.Pin("dev", "API_KEY", 1, "bob")
	if err := pm.Unpin("dev", "API_KEY"); err != nil {
		t.Fatalf("unpin error: %v", err)
	}
	pins, _ := pm.Load("dev")
	if _, ok := pins["API_KEY"]; ok {
		t.Error("expected API_KEY to be removed")
	}
}

func TestPinManager_Load_Missing(t *testing.T) {
	pm, _ := newTestPinManager(t)
	pins, err := pm.Load("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pins) != 0 {
		t.Errorf("expected empty pins, got %d", len(pins))
	}
}

func TestPinManager_MultipleKeys(t *testing.T) {
	pm, _ := newTestPinManager(t)
	_ = pm.Pin("prod", "SECRET_A", 2, "carol")
	_ = pm.Pin("prod", "SECRET_B", 5, "carol")
	pins, _ := pm.Load("prod")
	if len(pins) != 2 {
		t.Errorf("expected 2 pins, got %d", len(pins))
	}
}

func TestPinManager_CreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := base + "/nested/pins"
	pm := NewPinManager(dir)
	if err := pm.Pin("dev", "KEY", 1, "user"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("expected dir to be created: %v", err)
	}
}
