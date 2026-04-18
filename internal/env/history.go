package env

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// HistoryEntry records a single sync event.
type HistoryEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Path      string            `json:"path"`
	Added     []string          `json:"added,omitempty"`
	Updated   []string          `json:"updated,omitempty"`
	Removed   []string          `json:"removed,omitempty"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// HistoryManager persists sync history to a JSON file.
type HistoryManager struct {
	dir string
	now func() time.Time
}

// NewHistoryManager creates a HistoryManager storing entries under dir.
func NewHistoryManager(dir string) *HistoryManager {
	return &HistoryManager{dir: dir, now: time.Now}
}

func (h *HistoryManager) file(envPath string) string {
	safe := sanitizeFilename(envPath)
	return filepath.Join(h.dir, safe+".history.json")
}

// Append adds an entry for the given env file path.
func (h *HistoryManager) Append(envPath string, added, updated, removed []string, meta map[string]string) error {
	entries, _ := h.Load(envPath)
	entries = append(entries, HistoryEntry{
		Timestamp: h.now().UTC(),
		Path:      envPath,
		Added:     added,
		Updated:   updated,
		Removed:   removed,
		Meta:      meta,
	})
	if err := os.MkdirAll(h.dir, 0o700); err != nil {
		return fmt.Errorf("history: mkdir: %w", err)
	}
	f, err := os.Create(h.file(envPath))
	if err != nil {
		return fmt.Errorf("history: create: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

// Load returns all history entries for the given env file path.
func (h *HistoryManager) Load(envPath string) ([]HistoryEntry, error) {
	data, err := os.ReadFile(h.file(envPath))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("history: read: %w", err)
	}
	var entries []HistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("history: parse: %w", err)
	}
	return entries, nil
}

func sanitizeFilename(p string) string {
	safe := make([]byte, len(p))
	for i := 0; i < len(p); i++ {
		c := p[i]
		if c == '/' || c == '\\' || c == ':' {
			c = '_'
		}
		safe[i] = c
	}
	return string(safe)
}
