package env

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func fixedHistoryTime() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func newTestHistory(t *testing.T) *HistoryManager {
	t.Helper()
	dir := t.TempDir()
	hm := NewHistoryManager(dir)
	hm.now = fixedHistoryTime
	return hm
}

func TestHistory_AppendAndLoad(t *testing.T) {
	hm := newTestHistory(t)
	err := hm.Append(".env", []string{"FOO"}, []string{"BAR"}, nil, nil)
	if err != nil {
		t.Fatalf("Append: %v", err)
	}
	entries, err := hm.Load(".env")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Added[0] != "FOO" {
		t.Errorf("expected FOO in added")
	}
	if entries[0].Timestamp != fixedHistoryTime() {
		t.Errorf("unexpected timestamp")
	}
}

func TestHistory_MultipleAppends(t *testing.T) {
	hm := newTestHistory(t)
	_ = hm.Append(".env", []string{"A"}, nil, nil, nil)
	_ = hm.Append(".env", []string{"B"}, nil, nil, nil)
	entries, _ := hm.Load(".env")
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestHistory_Load_Missing(t *testing.T) {
	hm := newTestHistory(t)
	entries, err := hm.Load("nonexistent.env")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries")
	}
}

func TestHistory_CreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "nested", "history")
	hm := NewHistoryManager(dir)
	hm.now = fixedHistoryTime
	if err := hm.Append(".env", nil, nil, nil, nil); err != nil {
		t.Fatalf("Append: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("expected dir to exist: %v", err)
	}
}

func TestHistory_SanitizesFilename(t *testing.T) {
	hm := newTestHistory(t)
	_ = hm.Append("path/to/.env", []string{"X"}, nil, nil, nil)
	entries, err := hm.Load("path/to/.env")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
}
