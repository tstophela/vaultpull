package env

import (
	"fmt"
	"sort"
	"strings"
)

// TagManager manages key-value tags associated with secret keys.
type TagManager struct {
	tags map[string]map[string]string
}

// NewTagManager creates a new TagManager.
func NewTagManager() *TagManager {
	return &TagManager{tags: make(map[string]map[string]string)}
}

// Set assigns a tag to a secret key.
func (t *TagManager) Set(secretKey, tagKey, tagValue string) {
	if _, ok := t.tags[secretKey]; !ok {
		t.tags[secretKey] = make(map[string]string)
	}
	t.tags[secretKey][tagKey] = tagValue
}

// Get returns tags for a given secret key.
func (t *TagManager) Get(secretKey[string]string {
	return t.tags[secretKey]
}

// HasTag checks whether a secret key has a specific tag key/value pair.
func (t *TagManager) HasTag(secretKey, tagKey, tagValue string) bool {
	tags, ok := t.tags[secretKey]
	if !ok {
		return false
	}
	v, ok := tags[tagKey]
	return ok && v == tagValue
}

// FilterByTag returns all secret keys that have the given tag key/value.
func (t *TagManager) FilterByTag(tagKey, tagValue string) []string {
	var result []string
	for secretKey, tags := range t.tags {
		if v, ok := tags[tagKey]; ok && v == tagValue {
			result = append(result, secretKey)
		}
	}
	sort.Strings(result)
	return result
}

// Summary returns a human-readable summary of all tags.
func (t *TagManager) Summary() string {
	if len(t.tags) == 0 {
		return "no tags"
	}
	keys := make([]string, 0, len(t.tags))
	for k := range t.tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	for _, k := range keys {
		tagParts := make([]string, 0)
		for tk, tv := range t.tags[k] {
			tagParts = append(tagParts, fmt.Sprintf("%s=%s", tk, tv))
		}
		sort.Strings(tagParts)
		sb.WriteString(fmt.Sprintf("%s: [%s]\n", k, strings.Join(tagParts, ", ")))
	}
	return strings.TrimRight(sb.String(), "\n")
}
