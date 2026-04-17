package env

import "strings"

// SensitivePatterns are key substrings that trigger redaction.
var SensitivePatterns = []string{
	"PASSWORD", "SECRET", "TOKEN", "KEY", "PRIVATE", "CREDENTIAL", "AUTH",
}

// Redactor masks sensitive values in a secrets map for safe display.
type Redactor struct {
	patterns []string
}

// NewRedactor returns a Redactor using the default sensitive patterns.
func NewRedactor() *Redactor {
	return &Redactor{patterns: SensitivePatterns}
}

// NewRedactorWithPatterns returns a Redactor with custom patterns.
func NewRedactorWithPatterns(patterns []string) *Redactor {
	return &Redactor{patterns: patterns}
}

// IsSensitive returns true if the key matches any sensitive pattern.
func (r *Redactor) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range r.patterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// Redact returns a copy of secrets with sensitive values masked.
func (r *Redactor) Redact(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if r.IsSensitive(k) {
			out[k] = redactValue(v)
		} else {
			out[k] = v
		}
	}
	return out
}

// redactValue masks all but the last 4 characters of a value.
func redactValue(v string) string {
	if len(v) <= 4 {
		return "****"
	}
	return "****" + v[len(v)-4:]
}
