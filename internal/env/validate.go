package env

import (
	"fmt"
	"regexp"
	"strings"
)

var validKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ValidationError holds all issues found during validation.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s", strings.Join(e.Issues, "; "))
}

func (e *ValidationError) HasIssues() bool {
	return len(e.Issues) > 0
}

// Validator checks a map of env secrets for common issues.
type Validator struct {
	warnEmpty bool
}

// NewValidator creates a Validator. If warnEmpty is true, empty values are reported.
func NewValidator(warnEmpty bool) *Validator {
	return &Validator{warnEmpty: warnEmpty}
}

// Validate checks all keys and optionally values, returning a ValidationError if any issues exist.
func (v *Validator) Validate(secrets map[string]string) error {
	ve := &ValidationError{}
	for k, val := range secrets {
		if !validKeyPattern.MatchString(k) {
			ve.Issues = append(ve.Issues, fmt.Sprintf("invalid key %q: must match [A-Za-z_][A-Za-z0-9_]*", k))
		}
		if v.warnEmpty && val == "" {
			ve.Issues = append(ve.Issues, fmt.Sprintf("key %q has empty value", k))
		}
	}
	if ve.HasIssues() {
		return ve
	}
	return nil
}
