package env

import (
	"testing"
)

func TestLinter_NoIssues(t *testing.T) {
	l := NewLinter()
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_KEY":      "abc123",
	}
	issues := l.Lint(env)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestLinter_LowercaseKey(t *testing.T) {
	l := NewLinter(RuleNoLowercase)
	env := map[string]string{
		"db_host": "localhost",
	}
	issues := l.Lint(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Rule != RuleNoLowercase {
		t.Errorf("expected rule %s, got %s", RuleNoLowercase, issues[0].Rule)
	}
}

func TestLinter_SpaceInKey(t *testing.T) {
	l := NewLinter(RuleNoSpaces)
	env := map[string]string{
		"MY KEY": "value",
	}
	issues := l.Lint(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Rule != RuleNoSpaces {
		t.Errorf("unexpected rule: %s", issues[0].Rule)
	}
}

func TestLinter_TrailingSpaceInValue(t *testing.T) {
	l := NewLinter(RuleNoSpaces)
	env := map[string]string{
		"API_KEY": "secret ",
	}
	issues := l.Lint(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestLinter_EmptyKey(t *testing.T) {
	l := NewLinter(RuleNoEmptyKeys)
	env := map[string]string{
		"": "orphan",
	}
	issues := l.Lint(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Rule != RuleNoEmptyKeys {
		t.Errorf("unexpected rule: %s", issues[0].Rule)
	}
}

func TestLinter_MultipleRules(t *testing.T) {
	l := NewLinter(RuleNoLowercase, RuleNoSpaces)
	env := map[string]string{
		"my key": "value ",
	}
	issues := l.Lint(env)
	// lowercase key + space in key + trailing space in value = 3
	if len(issues) != 3 {
		t.Fatalf("expected 3 issues, got %d: %v", len(issues), issues)
	}
}

func TestLintIssue_String(t *testing.T) {
	i := LintIssue{Key: "foo", Rule: RuleNoLowercase, Message: "key should be uppercase"}
	s := i.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}
