package env_test

import (
	"os"
	"testing"

	"github.com/your-org/vaultpull/internal/env"
)

func newTestAliasManager(t *testing.T) *env.AliasManager {
	t.Helper()
	dir := t.TempDir()
	am, err := env.NewAliasManager(dir)
	if err != nil {
		t.Fatalf("NewAliasManager: %v", err)
	}
	return am
}

func TestAliasManager_SetAndGet(t *testing.T) {
	am := newTestAliasManager(t)

	if err := am.Set("prod", "secret/prod/app"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	path, ok := am.Get("prod")
	if !ok {
		t.Fatal("expected alias to exist")
	}
	if path != "secret/prod/app" {
		t.Errorf("got %q, want %q", path, "secret/prod/app")
	}
}

func TestAliasManager_Get_Missing(t *testing.T) {
	am := newTestAliasManager(t)
	_, ok := am.Get("nonexistent")
	if ok {
		t.Error("expected alias to be missing")
	}
}

func TestAliasManager_Delete(t *testing.T) {
	am := newTestAliasManager(t)
	_ = am.Set("staging", "secret/staging/app")

	if err := am.Delete("staging"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, ok := am.Get("staging")
	if ok {
		t.Error("expected alias to be deleted")
	}
}

func TestAliasManager_Delete_NonExistent(t *testing.T) {
	am := newTestAliasManager(t)
	if err := am.Delete("ghost"); err == nil {
		t.Error("expected error deleting non-existent alias")
	}
}

func TestAliasManager_List_Sorted(t *testing.T) {
	am := newTestAliasManager(t)
	_ = am.Set("zebra", "secret/z")
	_ = am.Set("alpha", "secret/a")
	_ = am.Set("mango", "secret/m")

	entries := am.List()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Name != "alpha" || entries[1].Name != "mango" || entries[2].Name != "zebra" {
		t.Errorf("unexpected order: %v", entries)
	}
}

func TestAliasManager_Persistence(t *testing.T) {
	dir := t.TempDir()
	am1, _ := env.NewAliasManager(dir)
	_ = am1.Set("dev", "secret/dev/app")

	am2, err := env.NewAliasManager(dir)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	path, ok := am2.Get("dev")
	if !ok || path != "secret/dev/app" {
		t.Errorf("alias not persisted; got %q, ok=%v", path, ok)
	}
}

func TestAliasManager_Set_EmptyAlias(t *testing.T) {
	am := newTestAliasManager(t)
	if err := am.Set("", "secret/path"); err == nil {
		t.Error("expected error for empty alias")
	}
}

func TestAliasManager_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	subdir := dir + "/nested/aliases"
	am, _ := env.NewAliasManager(subdir)
	_ = am.Set("x", "secret/x")
	if _, err := os.Stat(subdir); err != nil {
		t.Errorf("expected dir to be created: %v", err)
	}
}
