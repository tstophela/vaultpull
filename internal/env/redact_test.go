package env

import (
	"testing"
)

func TestRedactor_IsSensitive(t *testing.T) {
	r := NewRedactor()
	cases := []struct {
		key      string
		want     bool
	}{
		{"DB_PASSWORD", true},
		{"API_SECRET", true},
		{"AUTH_TOKEN", true},
		{"PRIVATE_KEY", true},
		{"APP_NAME", false},
		{"DATABASE_URL", false},
		{"PORT", false},
	}
	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			got := r.IsSensitive(tc.key)
			if got != tc.want {
				t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.want)
			}
		})
	}
}

func TestRedactor_Redact_MasksLongValues(t *testing.T) {
	r := NewRedactor()
	secrets := map[string]string{
		"API_TOKEN": "supersecretvalue1234",
		"APP_NAME":  "myapp",
	}
	out := r.Redact(secrets)
	if out["API_TOKEN"] != "****1234" {
		t.Errorf("expected ****1234, got %q", out["API_TOKEN"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected myapp, got %q", out["APP_NAME"])
	}
}

func TestRedactor_Redact_MasksShortValues(t *testing.T) {
	r := NewRedactor()
	secrets := map[string]string{
		"DB_PASSWORD": "abc",
	}
	out := r.Redact(secrets)
	if out["DB_PASSWORD"] != "****" {
		t.Errorf("expected ****, got %q", out["DB_PASSWORD"])
	}
}

func TestRedactor_CustomPatterns(t *testing.T) {
	r := NewRedactorWithPatterns([]string{"CUSTOM"})
	if !r.IsSensitive("MY_CUSTOM_FIELD") {
		t.Error("expected MY_CUSTOM_FIELD to be sensitive")
	}
	if r.IsSensitive("API_TOKEN") {
		t.Error("expected API_TOKEN to not be sensitive with custom patterns")
	}
}

func TestRedactValue_ExactlyFour(t *testing.T) {
	got := redactValue("abcd")
	if got != "****" {
		t.Errorf("expected ****, got %q", got)
	}
}
