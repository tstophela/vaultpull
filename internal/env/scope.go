package env

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Scope represents a named environment scope (e.g. "dev", "staging", "prod").
type Scope struct {
	Name string            `json:"name"`
	Path string            `json:"path"`
	Meta map[string]string `json:"meta,omitempty"`
}

// ScopeManager manages named scopes stored in a JSON index file.
type ScopeManager struct {
	indexPath string
}

// NewScopeManager creates a ScopeManager backed by the given index file path.
func NewScopeManager(indexPath string) *ScopeManager {
	return &ScopeManager{indexPath: indexPath}
}

func (m *ScopeManager) load() (map[string]Scope, error) {
	data, err := os.ReadFile(m.indexPath)
	if os.IsNotExist(err) {
		return map[string]Scope{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scope: read index: %w", err)
	}
	var scopes map[string]Scope
	if err := json.Unmarshal(data, &scopes); err != nil {
		return nil, fmt.Errorf("scope: parse index: %w", err)
	}
	return scopes, nil
}

func (m *ScopeManager) save(scopes map[string]Scope) error {
	if err := os.MkdirAll(filepath.Dir(m.indexPath), 0o700); err != nil {
		return fmt.Errorf("scope: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(scopes, "", "  ")
	if err != nil {
		return fmt.Errorf("scope: marshal: %w", err)
	}
	return os.WriteFile(m.indexPath, data, 0o600)
}

// Register adds or updates a scope entry.
func (m *ScopeManager) Register(s Scope) error {
	scopes, err := m.load()
	if err != nil {
		return err
	}
	scopes[s.Name] = s
	return m.save(scopes)
}

// Get returns the scope with the given name.
func (m *ScopeManager) Get(name string) (Scope, bool, error) {
	scopes, err := m.load()
	if err != nil {
		return Scope{}, false, err
	}
	s, ok := scopes[name]
	return s, ok, nil
}

// List returns all registered scopes sorted by name.
func (m *ScopeManager) List() ([]Scope, error) {
	scopes, err := m.load()
	if err != nil {
		return nil, err
	}
	out := make([]Scope, 0, len(scopes))
	for _, s := range scopes {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// Remove deletes a scope entry by name.
func (m *ScopeManager) Remove(name string) error {
	scopes, err := m.load()
	if err != nil {
		return err
	}
	delete(scopes, name)
	return m.save(scopes)
}
