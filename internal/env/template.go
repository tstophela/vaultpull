package env

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TemplateRenderer renders .env templates by substituting {{VAR}} placeholders
// with values from a provided secrets map.
type TemplateRenderer struct {
	placeholder *regexp.Regexp
	strict      bool
}

// NewTemplateRenderer creates a renderer. If strict is true, missing keys cause errors.
func NewTemplateRenderer(strict bool) *TemplateRenderer {
	return &TemplateRenderer{
		placeholder: regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`),
		strict:      strict,
	}
}

// RenderFile reads a template file and substitutes placeholders from secrets.
func (r *TemplateRenderer) RenderFile(path string, secrets map[string]string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("template: read %s: %w", path, err)
	}
	return r.Render(string(data), secrets)
}

// Render processes template content and returns resolved key=value pairs.
func (r *TemplateRenderer) Render(content string, secrets map[string]string) (map[string]string, error) {
	result := make(map[string]string)
	var missing []string

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := r.placeholder.ReplaceAllStringFunc(parts[1], func(match string) string {
			sub := r.placeholder.FindStringSubmatch(match)
			if len(sub) < 2 {
				return match
			}
			v, ok := secrets[sub[1]]
			if !ok {
				missing = append(missing, sub[1])
				return match
			}
			return v
		})
		result[key] = val
	}

	if r.strict && len(missing) > 0 {
		return nil, fmt.Errorf("template: unresolved placeholders: %s", strings.Join(missing, ", "))
	}
	return result, nil
}
