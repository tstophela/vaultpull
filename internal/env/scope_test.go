package env

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestScopeManager(t *testing.T) *ScopeManager {
	t.Helper()
	dir := t.TempDir()
	return NewScopeManager(filepath.Join(dir, "scopes.json"))
}

func TestScopeManager_RegisterAndGet(t *testing.T) {
	m := newTestScopeManager(t)

	s := Scope{Name: "dev", Path: ".env.dev", Meta: map[string]string{"region": "us-east-1"}}
	if err := m.Register(s); err != nil {
		t.Fatalf("Register: %v", err)
	}

	got, ok, err := m.Get("dev")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !ok {
		t.Fatal("expected scope to exist")
	}
	if got.Path != ".env.dev" {
		t.Errorf("path = %q, want .env.dev", got.Path)
	}
	if got.Meta["region"] != "us-east-1" {
		t.Errorf("meta region = %q, want us-east-1", got.Meta["region"])
	}
}

func TestScopeManager_Get_Missing(t *testing.T) {
	m := newTestScopeManager(t)
	_, ok, err := m.Get("prod")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if ok {
		t.Fatal("expected scope to be missing")
	}
}

func TestScopeManager_List_Sorted(t *testing.T) {
	m := newTestScopeManager(t)
	for _, name := range []string{"staging", "dev", "prod"} {
		if err := m.Register(Scope{Name: name, Path: ".env." + name}); err != nil {
			t.Fatalf("Register %s: %v", name, err)
		}
	}

	list, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("len = %d, want 3", len(list))
	}
	expected := []string{"dev", "prod", "staging"}
	for i, s := range list {
		if s.Name != expected[i] {
			t.Errorf("list[%d] = %q, want %q", i, s.Name, expected[i])
		}
	}
}

func TestScopeManager_Remove(t *testing.T) {
	m := newTestScopeManager(t)
	if err := m.Register(Scope{Name: "dev", Path: ".env.dev"}); err != nil {
		t.Fatalf("Register: %v", err)
	}
	if err := m.Remove("dev"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	_, ok, err := m.Get("dev")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if ok {
		t.Fatal("expected scope to be removed")
	}
}

func TestScopeManager_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "nested", "deep", "scopes.json")
	m := NewScopeManager(indexPath)
	if err := m.Register(Scope{Name: "dev", Path: ".env.dev"}); err != nil {
		t.Fatalf("Register: %v", err)
	}
	if _, err := os.Stat(indexPath); err != nil {
		t.Errorf("index file not created: %v", err)
	}
}
