package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeImportFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestImporter_BasicImport(t *testing.T) {
	dir := t.TempDir()
	src := writeImportFile(t, dir, "source.env", "FOO=bar\nBAZ=qux\n")

	im := NewImporter(mustParseStrategy("overwrite"), nil)
	result, res, err := im.ImportFile(src, map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if result["FOO"] != "bar" || result["BAZ"] != "qux" {
		t.Errorf("unexpected result: %v", result)
	}
	if res.Imported != 2 {
		t.Errorf("expected 2 imported, got %d", res.Imported)
	}
}

func TestImporter_PreserveStrategy(t *testing.T) {
	dir := t.TempDir()
	src := writeImportFile(t, dir, "source.env", "FOO=new\nBAR=added\n")

	existing := map[string]string{"FOO": "old"}
	im := NewImporter(mustParseStrategy("preserve"), nil)
	result, res, err := im.ImportFile(src, existing)
	if err != nil {
		t.Fatal(err)
	}
	if result["FOO"] != "old" {
		t.Errorf("preserve: FOO should remain 'old', got %s", result["FOO"])
	}
	if result["BAR"] != "added" {
		t.Errorf("preserve: BAR should be added, got %s", result["BAR"])
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.Skipped)
	}
}

func TestImporter_WithFilter(t *testing.T) {
	dir := t.TempDir()
	src := writeImportFile(t, dir, "source.env", "APP_FOO=1\nDB_BAR=2\n")

	f := NewFilter([]string{"APP_"}, nil)
	im := NewImporter(mustParseStrategy("overwrite"), f)
	result, res, err := im.ImportFile(src, map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := result["DB_BAR"]; ok {
		t.Error("DB_BAR should have been filtered out")
	}
	if res.Imported != 1 {
		t.Errorf("expected 1 imported, got %d", res.Imported)
	}
}

func TestImporter_MissingFile(t *testing.T) {
	im := NewImporter(mustParseStrategy("overwrite"), nil)
	_, _, err := im.ImportFile("/nonexistent/path.env", map[string]string{})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func mustParseStrategy(s string) Strategy {
	st, err := ParseStrategy(s)
	if err != nil {
		panic(err)
	}
	return st
}
