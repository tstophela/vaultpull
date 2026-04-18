package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/vaultpull/internal/env"
)

func TestHistoryCmd_NoHistory(t *testing.T) {
	dir := t.TempDir()
	historyDir = dir
	historyEnvFile = ".env"

	buf := &bytes.Buffer{}
	historyCmd.SetOut(buf)
	err := historyCmd.RunE(historyCmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No history") {
		t.Errorf("expected 'No history', got: %s", buf.String())
	}
}

func TestHistoryCmd_ShowsEntries(t *testing.T) {
	dir := t.TempDir()
	hm := env.NewHistoryManager(dir)
	_ = hm.Append(".env", []string{"FOO", "BAR"}, []string{"BAZ"}, nil, nil)

	historyDir = dir
	historyEnvFile = ".env"
	historyLimit = 10

	buf := &bytes.Buffer{}
	historyCmd.SetOut(buf)
	err := historyCmd.RunE(historyCmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "TIMESTAMP") {
		t.Errorf("expected header row")
	}
	if !strings.Contains(out, "2") {
		t.Errorf("expected added count 2")
	}
}

func TestHistoryCmd_LimitEntries(t *testing.T) {
	dir := t.TempDir()
	hm := env.NewHistoryManager(filepath.Join(dir))
	for i := 0; i < 5; i++ {
		_ = hm.Append(".env", []string{"K"}, nil, nil, nil)
	}

	historyDir = dir
	historyEnvFile = ".env"
	historyLimit = 2

	buf := &bytes.Buffer{}
	historyCmd.SetOut(buf)
	err := historyCmd.RunE(historyCmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + 2 data lines
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header+2), got %d", len(lines))
	}
}
