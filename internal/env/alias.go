package env

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// AliasManager maps short alias names to Vault secret paths.
type AliasManager struct {
	dir   string
	aliases map[string]string
}

// NewAliasManager creates an AliasManager backed by a JSON file in dir.
func NewAliasManager(dir string) (*AliasManager, error) {
	am := &AliasManager{dir: dir, aliases: make(map[string]string)}
	if err := am.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("alias: load: %w", err)
	}
	return am, nil
}

func (am *AliasManager) filePath() string {
	return filepath.Join(am.dir, "aliases.json")
}

func (am *AliasManager) load() error {
	data, err := os.ReadFile(am.filePath())
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &am.aliases)
}

func (am *AliasManager) save() error {
	if err := os.MkdirAll(am.dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(am.aliases, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(am.filePath(), data, 0o600)
}

// Set registers or updates an alias.
func (am *AliasManager) Set(alias, path string) error {
	if alias == "" || path == "" {
		return fmt.Errorf("alias: alias and path must not be empty")
	}
	am.aliases[alias] = path
	return am.save()
}

// Get resolves an alias to its Vault path. Returns empty string if not found.
func (am *AliasManager) Get(alias string) (string, bool) {
	v, ok := am.aliases[alias]
	return v, ok
}

// Delete removes an alias.
func (am *AliasManager) Delete(alias string) error {
	if _, ok := am.aliases[alias]; !ok {
		return fmt.Errorf("alias: %q not found", alias)
	}
	delete(am.aliases, alias)
	return am.save()
}

// List returns all aliases sorted by name.
func (am *AliasManager) List() []AliasEntry {
	entries := make([]AliasEntry, 0, len(am.aliases))
	for k, v := range am.aliases {
		entries = append(entries, AliasEntry{Name: k, Path: v})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })
	return entries
}

// AliasEntry holds a single alias name and its resolved path.
type AliasEntry struct {
	Name string
	Path string
}
