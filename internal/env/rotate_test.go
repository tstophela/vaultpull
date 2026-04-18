package env

import (
	"strings"
	"testing"
	"time"
)

func fixedRotateTime() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func newTestRotator() *Rotator {
	r := NewRotator(nil)
	r.now = fixedRotateTime
	return r
}

func TestRotator_Detect_ChangedKeys(t *testing.T) {
	rot := newTestRotator()
	existing := map[string]string{"DB_PASSWORD": "old", "API_KEY": "abc"}
	incoming := map[string]string{"DB_PASSWORD": "new", "API_KEY": "abc"}

	records := rot.Detect(existing, incoming)
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].Key != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD, got %s", records[0].Key)
	}
}

func TestRotator_Detect_NewKey_NotReported(t *testing.T) {
	rot := newTestRotator()
	existing := map[string]string{}
	incoming := map[string]string{"NEW_SECRET": "value"}

	records := rot.Detect(existing, incoming)
	if len(records) != 0 {
		t.Errorf("expected no records for new key, got %d", len(records))
	}
}

func TestRotator_Detect_RedactsSensitive(t *testing.T) {
	rot := newTestRotator()
	existing := map[string]string{"SECRET_TOKEN": "plaintext-old-value"}
	incoming := map[string]string{"SECRET_TOKEN": "plaintext-new-value"}

	records := rot.Detect(existing, incoming)
	if len(records) != 1 {
		t.Fatalf("expected 1 record")
	}
	if records[0].OldValue == "plaintext-old-value" {
		t.Error("expected old value to be redacted")
	}
	if records[0].NewValue == "plaintext-new-value" {
		t.Error("expected new value to be redacted")
	}
}

func TestSummary_Empty(t *testing.T) {
	s := Summary(nil)
	if s != "no secrets rotated" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_WithRecords(t *testing.T) {
	records := []RotationRecord{
		{Key: "API_KEY", OldValue: "***", NewValue: "***", RotatedAt: fixedRotateTime()},
	}
	s := Summary(records)
	if !strings.Contains(s, "1 secret(s) rotated") {
		t.Errorf("unexpected summary: %s", s)
	}
	if !strings.Contains(s, "API_KEY") {
		t.Errorf("expected API_KEY in summary")
	}
}
