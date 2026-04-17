package env

import (
	"fmt"
	"strings"
	"unicode"
)

// Formatter controls how key=value pairs are serialized to .env format.
type Formatter struct {
	QuoteAll bool
}

// NewFormatter returns a Formatter with default settings.
func NewFormatter(quoteAll bool) *Formatter {
	return &Formatter{QuoteAll: quoteAll}
}

// FormatLine serializes a single key-value pair into a .env line.
func (f *Formatter) FormatLine(key, value string) string {
	if f.QuoteAll || needsQuoting(value) {
		escaped := strings.ReplaceAll(value, `"`, `\"`)
		return fmt.Sprintf(`%s="%s"`, key, escaped)
	}
	return fmt.Sprintf("%s=%s", key, value)
}

// FormatMap serializes a map of secrets into a slice of .env lines, sorted by key.
func (f *Formatter) FormatMap(secrets map[string]string, keys []string) []string {
	lines := make([]string, 0, len(keys))
	for _, k := range keys {
		v, ok := secrets[k]
		if !ok {
			continue
		}
		lines = append(lines, f.FormatLine(k, v))
	}
	return lines
}

// needsQuoting returns true if the value contains whitespace, quotes, or special shell chars.
func needsQuoting(v string) bool {
	if v == "" {
		return false
	}
	for _, r := range v {
		if unicode.IsSpace(r) {
			return true
		}
		switch r {
		case '#', '$', '&', '*', '(', ')', '{', '}', '|', ';', '<', '>', '`', '!', '\\', '\'', '"':
			return true
		}
	}
	return false
}
