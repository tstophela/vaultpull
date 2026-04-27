package env

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Namespace represents a named grouping of environment keys with an optional prefix.
type Namespace struct {
	Name   string   `json:"name"`
	Prefix string   `json:"prefix,omitempty"`
	Keys   []string `json:"keys,omitempty"`
}

// NamespaceManager manages named namespaces persisted to disk.
type NamespaceManager struct {
	dir string
}

// NewNamespaceManager returns a NamespaceManager that stores data under dir.
func NewNamespaceManager(dir string) *NamespaceManager {
	return &NamespaceManager{dir: dir}
}

func (m *NamespaceManager) filePath(name string) string {
	return filepath.Join(m.dir, "ns_"+sanitizeFilename(name)+".json")
}

// Save persists a namespace to disk.
func (m *NamespaceManager) Save(ns Namespace) error {
	if err := os.MkdirAll(m.dir, 0o700); err != nil {
		return fmt.Errorf("namespace: create dir: %w", err)
	}
	data, err := json.MarshalIndent(ns, "", "  ")
	if err != nil {
		return fmt.Errorf("namespace: marshal: %w", err)
	}
	return os.WriteFile(m.filePath(ns.Name), data, 0o600)
}

// Get loads a namespace by name.
func (m *NamespaceManager) Get(name string) (Namespace, error) {
	data, err := os.ReadFile(m.filePath(name))
	if err != nil {
		if os.IsNotExist(err) {
			return Namespace{}, fmt.Errorf("namespace %q not found", name)
		}
		return Namespace{}, fmt.Errorf("namespace: read: %w", err)
	}
	var ns Namespace
	if err := json.Unmarshal(data, &ns); err != nil {
		return Namespace{}, fmt.Errorf("namespace: unmarshal: %w", err)
	}
	return ns, nil
}

// Delete removes a namespace by name.
func (m *NamespaceManager) Delete(name string) error {
	err := os.Remove(m.filePath(name))
	if os.IsNotExist(err) {
		return fmt.Errorf("namespace %q not found", name)
	}
	return err
}

// List returns all saved namespaces sorted by name.
func (m *NamespaceManager) List() ([]Namespace, error) {
	pattern := filepath.Join(m.dir, "ns_*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("namespace: glob: %w", err)
	}
	var result []Namespace
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var ns Namespace
		if json.Unmarshal(data, &ns) == nil {
			result = append(result, ns)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result, nil
}

// FilterKeys returns only the env map entries whose keys belong to the namespace.
func (m *NamespaceManager) FilterKeys(ns Namespace, env map[string]string) map[string]string {
	out := make(map[string]string)
	allowed := make(map[string]bool, len(ns.Keys))
	for _, k := range ns.Keys {
		allowed[k] = true
	}
	for k, v := range env {
		if len(ns.Keys) == 0 || allowed[k] {
			out[k] = v
		}
	}
	return out
}
