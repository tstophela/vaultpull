package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriter_Write_BasicSecrets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := NewWriter(path, false)
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "abc123",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}

	content := string(data)
	for key, val := range secrets {
		expected := key + "=" + val
		if !strings.Contains(content, expected) {
			t.Errorf("expected %q in output, got:\n%s", expected, content)
		}
	}
}

func TestWriter_Write_QuotesSpecialValues(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := NewWriter(path, false)
	secrets := map[string]string{
		"NOTE": "hello world",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), `NOTE="hello world"`) {
		t.Errorf("expected quoted value, got: %s", string(data))
	}
}

func TestWriter_Write_BackupCreated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	// Write an initial file to back up.
	if err := os.WriteFile(path, []byte("OLD=value\n"), 0600); err != nil {
		t.Fatal(err)
	}

	w := NewWriter(path, true)
	if err := w.Write(map[string]string{"NEW": "value"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path + ".bak"); os.IsNotExist(err) {
		t.Error("expected backup file to exist")
	}

	bak, _ := os.ReadFile(path + ".bak")
	if !strings.Contains(string(bak), "OLD=value") {
		t.Errorf("backup should contain original content, got: %s", string(bak))
	}
}

func TestWriter_Write_SortedOutput(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := NewWriter(path, false)
	secrets := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	if err := w.Write(secrets); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 || !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected sorted output, got: %v", lines)
	}
}
