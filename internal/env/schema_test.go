package env

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeSchemaFile(t *testing.T, dir string, s Schema) string {
	t.Helper()
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshal schema: %v", err)
	}
	p := filepath.Join(dir, "schema.json")
	if err := os.WriteFile(p, data, 0o600); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	return p
}

func TestSchemaManager_ValidEnv(t *testing.T) {
	dir := t.TempDir()
	path := writeSchemaFile(t, dir, Schema{
		"APP_ENV": {Required: true, Pattern: `^(dev|staging|prod)$`},
		"PORT":    {Required: true, Pattern: `^\d+$`},
	})
	sm, err := NewSchemaManager(path)
	if err != nil {
		t.Fatalf("NewSchemaManager: %v", err)
	}
	env := map[string]string{"APP_ENV": "prod", "PORT": "8080"}
	violations := sm.Validate(env)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestSchemaManager_MissingRequired(t *testing.T) {
	dir := t.TempDir()
	path := writeSchemaFile(t, dir, Schema{
		"DB_URL": {Required: true},
	})
	sm, _ := NewSchemaManager(path)
	violations := sm.Validate(map[string]string{})
	if len(violations) != 1 || violations[0].Key != "DB_URL" {
		t.Errorf("expected DB_URL violation, got %v", violations)
	}
}

func TestSchemaManager_PatternMismatch(t *testing.T) {
	dir := t.TempDir()
	path := writeSchemaFile(t, dir, Schema{
		"LOG_LEVEL": {Required: true, Pattern: `^(debug|info|warn|error)$`},
	})
	sm, _ := NewSchemaManager(path)
	violations := sm.Validate(map[string]string{"LOG_LEVEL": "verbose"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "LOG_LEVEL" {
		t.Errorf("unexpected key %q", violations[0].Key)
	}
}

func TestSchemaManager_ApplyDefaults(t *testing.T) {
	dir := t.TempDir()
	path := writeSchemaFile(t, dir, Schema{
		"TIMEOUT": {Default: "30s"},
		"RETRIES": {Default: "3"},
	})
	sm, _ := NewSchemaManager(path)
	env := map[string]string{"RETRIES": "5"}
	out := sm.ApplyDefaults(env)
	if out["TIMEOUT"] != "30s" {
		t.Errorf("expected TIMEOUT=30s, got %q", out["TIMEOUT"])
	}
	if out["RETRIES"] != "5" {
		t.Errorf("expected RETRIES=5 (preserved), got %q", out["RETRIES"])
	}
}

func TestNewSchemaManager_MissingFile(t *testing.T) {
	_, err := NewSchemaManager("/nonexistent/schema.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
