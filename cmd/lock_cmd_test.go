package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestLockCmd_AcquireAndRelease(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, ".env")
	lockDir := filepath.Join(dir, "locks")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)

	rootCmd.SetArgs([]string{"lock", "acquire", "--lock-dir", lockDir, target})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("acquire: %v", err)
	}
	if !strings.Contains(buf.String(), "Locked") {
		t.Errorf("expected Locked in output, got: %s", buf.String())
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"lock", "release", "--lock-dir", lockDir, target})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("release: %v", err)
	}
	if !strings.Contains(buf.String(), "Released") {
		t.Errorf("expected Released in output, got: %s", buf.String())
	}
}

func TestLockCmd_DoubleAcquireFails(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, ".env")
	lockDir := filepath.Join(dir, "locks")

	rootCmd.SetArgs([]string{"lock", "acquire", "--lock-dir", lockDir, target})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("first acquire: %v", err)
	}
	defer func() {
		rootCmd.SetArgs([]string{"lock", "release", "--lock-dir", lockDir, target})
		rootCmd.Execute() //nolint
	}()

	rootCmd.SetArgs([]string{"lock", "acquire", "--lock-dir", lockDir, target})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error on double acquire")
	}
}
