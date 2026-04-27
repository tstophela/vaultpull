package env

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// SchemaField defines expectations for a single env key.
type SchemaField struct {
	Required bool   `json:"required"`
	Pattern  string `json:"pattern,omitempty"`
	Default  string `json:"default,omitempty"`
}

// Schema maps key names to their field definitions.
type Schema map[string]SchemaField

// SchemaViolation describes a single schema validation failure.
type SchemaViolation struct {
	Key     string
	Message string
}

func (v SchemaViolation) Error() string {
	return fmt.Sprintf("schema violation [%s]: %s", v.Key, v.Message)
}

// SchemaManager loads and applies a JSON schema against env maps.
type SchemaManager struct {
	schema Schema
}

// NewSchemaManager loads a schema from a JSON file.
func NewSchemaManager(path string) (*SchemaManager, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("reading schema file: %w", err)
	}
	var s Schema
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing schema file: %w", err)
	}
	return &SchemaManager{schema: s}, nil
}

// Validate checks env against the schema and returns all violations.
func (m *SchemaManager) Validate(env map[string]string) []SchemaViolation {
	var violations []SchemaViolation
	for key, field := range m.schema {
		val, exists := env[key]
		if !exists || val == "" {
			if field.Required {
				violations = append(violations, SchemaViolation{
					Key:     key,
					Message: "required key is missing or empty",
				})
			}
			continue
		}
		if field.Pattern != "" {
			re, err := regexp.Compile(field.Pattern)
			if err != nil {
				violations = append(violations, SchemaViolation{
					Key:     key,
					Message: fmt.Sprintf("invalid pattern %q: %v", field.Pattern, err),
				})
				continue
			}
			if !re.MatchString(val) {
				violations = append(violations, SchemaViolation{
					Key:     key,
					Message: fmt.Sprintf("value %q does not match pattern %q", val, field.Pattern),
				})
			}
		}
	}
	return violations
}

// ApplyDefaults fills in default values for missing keys defined in the schema.
func (m *SchemaManager) ApplyDefaults(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	for key, field := range m.schema {
		if _, exists := out[key]; !exists && field.Default != "" {
			out[key] = field.Default
		}
	}
	return out
}
