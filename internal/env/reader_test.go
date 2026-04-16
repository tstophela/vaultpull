package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeEnvFile: %v", err)
	}
	return p
}

func TestReader_Read_BasicParsing(t *testing.T) {
	p := writeEnvFile(t, "KEY1=value1\nKEY2=value2\n")
	r := NewReader(p)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY1"] != "value1" || got["KEY2"] != "value2" {
		t.Errorf("unexpected values: %v", got)
	}
}

func TestReader_Read_IgnoresComments(t *testing.T) {
	p := writeEnvFile(t, "# comment\nKEY=val\n")
	got, err := NewReader(p).Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["# comment"]; ok {
		t.Error("comment should be ignored")
	}
	if got["KEY"] != "val" {
		t.Errorf("expected val, got %q", got["KEY"])
	}
}

func TestReader_Read_QuotedValues(t *testing.T) {
	p := writeEnvFile(t, `KEY="hello world"`+"\n")
	got, err := NewReader(p).Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", got["KEY"])
	}
}

func TestReader_Read_MissingFile(t *testing.T) {
	r := NewReader("/nonexistent/.env")
	got, err := r.Read()
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}
