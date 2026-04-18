package env

import (
	"strings"
	"testing"
)

func TestValidator_ValidKeys(t *testing.T) {
	v := NewValidator(false)
	secrets := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"_PRIVATE":     "value",
		"key123":       "val",
	}
	if err := v.Validate(secrets); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidator_InvalidKey(t *testing.T) {
	v := NewValidator(false)
	secrets := map[string]string{
		"123INVALID": "value",
	}
	err := v.Validate(secrets)
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
	if !strings.Contains(err.Error(), "invalid key") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidator_EmptyValue_WarnEnabled(t *testing.T) {
	v := NewValidator(true)
	secrets := map[string]string{
		"GOOD_KEY": "",
	}
	err := v.Validate(secrets)
	if err == nil {
		t.Fatal("expected error for empty value")
	}
	if !strings.Contains(err.Error(), "empty value") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidator_EmptyValue_WarnDisabled(t *testing.T) {
	v := NewValidator(false)
	secrets := map[string]string{
		"GOOD_KEY": "",
	}
	if err := v.Validate(secrets); err != nil {
		t.Fatalf("expected no error when warnEmpty=false, got: %v", err)
	}
}

func TestValidationError_HasIssues(t *testing.T) {
	ve := &ValidationError{Issues: []string{"issue1", "issue2"}}
	if !ve.HasIssues() {
		t.Error("expected HasIssues to return true")
	}
	empty := &ValidationError{}
	if empty.HasIssues() {
		t.Error("expected HasIssues to return false for empty issues")
	}
}
