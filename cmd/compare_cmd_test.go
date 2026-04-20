package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func writeCompareEnvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeCompareEnvFile: %v", err)
	}
	return p
}

func TestCompareCmd_AddedAndRemoved(t *testing.T) {
	dir := t.TempDir()
	writeCompareEnvFile(t, dir, "a.env", "FOO=bar\nBAZ=qux\n")
	writeCompareEnvFile(t, dir, "b.env", "FOO=bar\nNEW=value\n")

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	rootCmd.SetArgs([]string{
		"compare",
		"--file-a", filepath.Join(dir, "a.env"),
		"--file-b", filepath.Join(dir, "b.env"),
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if out == "" {
		t.Error("expected non-empty compare output")
	}
}

func TestCompareCmd_Identical(t *testing.T) {
	dir := t.TempDir()
	writeCompareEnvFile(t, dir, "a.env", "FOO=bar\nBAZ=qux\n")
	writeCompareEnvFile(t, dir, "b.env", "FOO=bar\nBAZ=qux\n")

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	rootCmd.SetArgs([]string{
		"compare",
		"--file-a", filepath.Join(dir, "a.env"),
		"--file-b", filepath.Join(dir, "b.env"),
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCompareCmd_MissingFileA(t *testing.T) {
	dir := t.TempDir()
	writeCompareEnvFile(t, dir, "b.env", "FOO=bar\n")

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	rootCmd.SetArgs([]string{
		"compare",
		"--file-a", filepath.Join(dir, "missing.env"),
		"--file-b", filepath.Join(dir, "b.env"),
	})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for missing file-a")
	}
}
