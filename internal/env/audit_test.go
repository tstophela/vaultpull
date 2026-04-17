package env

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func fixedTime() time.Time {
	t, _ := time.Parse(time.RFC3339, "2024-01-15T10:00:00Z")
	return t
}

func TestAuditLogger_Log_AllCategories(t *testing.T) {
	var buf bytes.Buffer
	logger := NewAuditLogger(&buf)

	entry := AuditEntry{
		Timestamp: fixedTime(),
		Path:      "secret/myapp",
		Added:     []string{"NEW_KEY"},
		Updated:   []string{"CHANGED_KEY"},
		Unchanged: []string{"STABLE_KEY"},
	}
	logger.Log(entry)

	out := buf.String()
	if !strings.Contains(out, "secret/myapp") {
		t.Error("expected path in output")
	}
	if !strings.Contains(out, "+ NEW_KEY") {
		t.Error("expected added key")
	}
	if !strings.Contains(out, "~ CHANGED_KEY") {
		t.Error("expected updated key")
	}
	if !strings.Contains(out, "= STABLE_KEY") {
		t.Error("expected unchanged key")
	}
	if !strings.Contains(out, "1 added, 1 updated, 1 unchanged") {
		t.Error("expected summary line")
	}
}

func TestAuditLogger_Log_EmptyEntry(t *testing.T) {
	var buf bytes.Buffer
	logger := NewAuditLogger(&buf)

	entry := AuditEntry{
		Timestamp: fixedTime(),
		Path:      "secret/empty",
	}
	logger.Log(entry)

	out := buf.String()
	if !strings.Contains(out, "0 added, 0 updated, 0 unchanged") {
		t.Error("expected zero summary")
	}
}

func TestNewAuditLogger_NilUsesStdout(t *testing.T) {
	logger := NewAuditLogger(nil)
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
}
