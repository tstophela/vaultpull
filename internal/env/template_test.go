package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTemplateRenderer_Render_Basic(t *testing.T) {
	r := NewTemplateRenderer(false)
	content := `
DB_HOST={{ DB_HOST }}
DB_PASS={{DB_PASS}}
`
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PASS": "s3cr3t"}
	out, err := r.Render(content, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", out["DB_HOST"])
	}
	if out["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected s3cr3t, got %s", out["DB_PASS"])
	}
}

func TestTemplateRenderer_Render_MissingNonStrict(t *testing.T) {
	r := NewTemplateRenderer(false)
	content := "API_KEY={{MISSING_KEY}}\n"
	out, err := r.Render(content, map[string]string{})
	if err != nil {
		t.Fatalf("expected no error in non-strict mode, got %v", err)
	}
	if out["API_KEY"] != "{{MISSING_KEY}}" {
		t.Errorf("expected placeholder preserved, got %s", out["API_KEY"])
	}
}

func TestTemplateRenderer_Render_MissingStrict(t *testing.T) {
	r := NewTemplateRenderer(true)
	content := "TOKEN={{SECRET_TOKEN}}\n"
	_, err := r.Render(content, map[string]string{})
	if err == nil {
		t.Fatal("expected error in strict mode for missing key")
	}
}

func TestTemplateRenderer_Render_IgnoresComments(t *testing.T) {
	r := NewTemplateRenderer(true)
	content := "# this is a comment\nFOO=bar\n"
	out, err := r.Render(content, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected bar, got %s", out["FOO"])
	}
	if _, ok := out["# this is a comment"]; ok {
		t.Error("comment line should not be parsed as key")
	}
}

func TestTemplateRenderer_RenderFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.tpl")
	_ = os.WriteFile(path, []byte("SVC_URL={{BASE_URL}}/api\n"), 0600)

	r := NewTemplateRenderer(true)
	out, err := r.RenderFile(path, map[string]string{"BASE_URL": "https://example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SVC_URL"] != "https://example.com/api" {
		t.Errorf("unexpected value: %s", out["SVC_URL"])
	}
}

func TestTemplateRenderer_RenderFile_Missing(t *testing.T) {
	r := NewTemplateRenderer(false)
	_, err := r.RenderFile("/nonexistent/.env.tpl", nil)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
