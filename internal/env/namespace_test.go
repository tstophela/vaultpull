package env

import (
	"os"
	"testing"
)

func newTestNamespaceManager(t *testing.T) *NamespaceManager {
	t.Helper()
	dir := t.TempDir()
	return NewNamespaceManager(dir)
}

func TestNamespaceManager_SaveAndGet(t *testing.T) {
	m := newTestNamespaceManager(t)
	ns := Namespace{Name: "production", Prefix: "PROD_", Keys: []string{"DB_URL", "API_KEY"}}

	if err := m.Save(ns); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := m.Get("production")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != ns.Name || got.Prefix != ns.Prefix {
		t.Errorf("got %+v, want %+v", got, ns)
	}
	if len(got.Keys) != 2 || got.Keys[0] != "DB_URL" {
		t.Errorf("unexpected keys: %v", got.Keys)
	}
}

func TestNamespaceManager_Get_Missing(t *testing.T) {
	m := newTestNamespaceManager(t)
	_, err := m.Get("ghost")
	if err == nil {
		t.Fatal("expected error for missing namespace")
	}
}

func TestNamespaceManager_Delete(t *testing.T) {
	m := newTestNamespaceManager(t)
	ns := Namespace{Name: "staging"}
	_ = m.Save(ns)

	if err := m.Delete("staging"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := m.Get("staging"); err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestNamespaceManager_Delete_NonExistent(t *testing.T) {
	m := newTestNamespaceManager(t)
	if err := m.Delete("nope"); err == nil {
		t.Fatal("expected error deleting non-existent namespace")
	}
}

func TestNamespaceManager_List_Sorted(t *testing.T) {
	m := newTestNamespaceManager(t)
	for _, name := range []string{"zebra", "alpha", "mango"} {
		_ = m.Save(Namespace{Name: name})
	}

	list, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3, got %d", len(list))
	}
	if list[0].Name != "alpha" || list[1].Name != "mango" || list[2].Name != "zebra" {
		t.Errorf("wrong order: %v", list)
	}
}

func TestNamespaceManager_List_Empty(t *testing.T) {
	m := newTestNamespaceManager(t)
	list, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %v", list)
	}
}

func TestNamespaceManager_FilterKeys_WithKeys(t *testing.T) {
	m := newTestNamespaceManager(t)
	ns := Namespace{Name: "limited", Keys: []string{"DB_URL", "PORT"}}
	env := map[string]string{"DB_URL": "postgres://", "PORT": "5432", "SECRET": "hidden"}

	out := m.FilterKeys(ns, env)
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d: %v", len(out), out)
	}
	if _, ok := out["SECRET"]; ok {
		t.Error("SECRET should have been filtered out")
	}
}

func TestNamespaceManager_FilterKeys_NoKeys_ReturnsAll(t *testing.T) {
	m := newTestNamespaceManager(t)
	ns := Namespace{Name: "open"}
	env := map[string]string{"A": "1", "B": "2"}

	out := m.FilterKeys(ns, env)
	if len(out) != 2 {
		t.Errorf("expected all keys, got %d", len(out))
	}
}

func TestNamespaceManager_Save_CreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := base + "/deep/ns"
	m := NewNamespaceManager(dir)

	if err := m.Save(Namespace{Name: "test"}); err != nil {
		t.Fatalf("Save with new dir: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("dir not created: %v", err)
	}
}
