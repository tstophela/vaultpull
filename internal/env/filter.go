package env

import "strings"

// Filter holds rules for including or excluding secret keys.
type Filter struct {
	prefixes []string
	excludes map[string]struct{}
}

// NewFilter creates a Filter from include-prefix and exclude-key lists.
func NewFilter(prefixes []string, excludes []string) *Filter {
	exMap := make(map[string]struct{}, len(excludes))
	for _, e := range excludes {
		exMap[strings.ToUpper(e)] = struct{}{}
	}
	return &Filter{prefixes: prefixes, excludes: exMap}
}

// Allow returns true when the key should be included in the output.
func (f *Filter) Allow(key string) bool {
	upper := strings.ToUpper(key)

	if _, excluded := f.excludes[upper]; excluded {
		return false
	}

	if len(f.prefixes) == 0 {
		return true
	}

	for _, p := range f.prefixes {
		if strings.HasPrefix(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

// Apply filters a secrets map, returning only allowed keys.
func (f *Filter) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range secrets {
		if f.Allow(k) {
			out[k] = v
		}
	}
	return out
}
