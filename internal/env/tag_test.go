package env

import (
	"testing"
)

func TestTagManager_SetAndGet(t *testing.T) {
	tm := NewTagManager()
	tm.Set("DB_PASSWORD", "env", "production")
	tm.Set("DB_PASSWORD", "owner", "infra")

	tags := tm.Get("DB_PASSWORD")
	if tags["env"] != "production" {
		t.Errorf("expected production, got %s", tags["env"])
	}
	if tags["owner"] != "infra" {
		t.Errorf("expected infra, got %s", tags["owner"])
	}
}

func TestTagManager_HasTag(t *testing.T) {
	tm := NewTagManager()
	tm.Set("API_KEY", "env", "staging")

	if !tm.HasTag("API_KEY", "env", "staging") {
		t.Error("expected HasTag to return true")
	}
	if tm.HasTag("API_KEY", "env", "production") {
		t.Error("expected HasTag to return false for wrong value")
	}
	if tm.HasTag("MISSING", "env", "staging") {
		t.Error("expected HasTag to return false for missing key")
	}
}

func TestTagManager_FilterByTag(t *testing.T) {
	tm := NewTagManager()
	tm.Set("DB_PASSWORD", "env", "production")
	tm.Set("API_KEY", "env", "production")
	tm.Set("DEV_TOKEN", "env", "staging")

	result := tm.FilterByTag("env", "production")
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
	if result[0] != "API_KEY" || result[1] != "DB_PASSWORD" {
		t.Errorf("unexpected order or values: %v", result)
	}
}

func TestTagManager_Summary_Empty(t *testing.T) {
	tm := NewTagManager()
	if tm.Summary() != "no tags" {
		t.Errorf("expected 'no tags', got %s", tm.Summary())
	}
}

func TestTagManager_Summary_WithTags(t *testing.T) {
	tm := NewTagManager()
	tm.Set("API_KEY", "env", "prod")

	s := tm.Summary()
	if s == "" || s == "no tags" {
		t.Error("expected non-empty summary")
	}
	if !containsStr(s, "API_KEY") {
		t.Errorf("expected API_KEY in summary, got: %s", s)
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
