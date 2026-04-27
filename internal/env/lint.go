package env

import (
	"fmt"
	"strings"
)

// LintRule represents a single linting rule applied to env keys/values.
type LintRule string

const (
	RuleNoLowercase  LintRule = "no_lowercase"
	RuleNoSpaces     LintRule = "no_spaces"
	RuleNoDuplicates LintRule = "no_duplicates"
	RuleNoEmptyKeys  LintRule = "no_empty_keys"
)

// LintIssue describes a single lint violation.
type LintIssue struct {
	Key     string
	Rule    LintRule
	Message string
}

func (i LintIssue) String() string {
	return fmt.Sprintf("[%s] %s: %s", i.Rule, i.Key, i.Message)
}

// Linter checks env maps for common issues.
type Linter struct {
	rules []LintRule
}

// NewLinter creates a Linter with the given rules enabled.
// If no rules are provided, all default rules are used.
func NewLinter(rules ...LintRule) *Linter {
	if len(rules) == 0 {
		rules = []LintRule{RuleNoLowercase, RuleNoSpaces, RuleNoDuplicates, RuleNoEmptyKeys}
	}
	return &Linter{rules: rules}
}

// Lint runs all enabled rules against the provided env map and returns any issues found.
func (l *Linter) Lint(env map[string]string) []LintIssue {
	var issues []LintIssue
	seen := make(map[string]bool)

	for key, val := range env {
		for _, rule := range l.rules {
			switch rule {
			case RuleNoEmptyKeys:
				if strings.TrimSpace(key) == "" {
					issues = append(issues, LintIssue{Key: key, Rule: rule, Message: "key must not be empty"})
				}
			case RuleNoLowercase:
				if key != strings.ToUpper(key) {
					issues = append(issues, LintIssue{Key: key, Rule: rule, Message: "key should be uppercase"})
				}
			case RuleNoSpaces:
				if strings.Contains(key, " ") {
					issues = append(issues, LintIssue{Key: key, Rule: rule, Message: "key must not contain spaces"})
				}
				if strings.HasPrefix(val, " ") || strings.HasSuffix(val, " ") {
					issues = append(issues, LintIssue{Key: key, Rule: rule, Message: "value has leading or trailing spaces"})
				}
			case RuleNoDuplicates:
				norm := strings.ToUpper(key)
				if seen[norm] {
					issues = append(issues, LintIssue{Key: key, Rule: rule, Message: "duplicate key (case-insensitive)"})
				}
				seen[norm] = true
			}
		}
	}
	return issues
}
