package env

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures the state of secrets at a point in time.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Path      string            `json:"path"`
	Secrets   map[string]string `json:"secrets"`
}

// SnapshotManager handles saving and loading snapshots.
type SnapshotManager struct {
	dir string
}

// NewSnapshotManager creates a manager that stores snapshots in dir.
func NewSnapshotManager(dir string) *SnapshotManager {
	return &SnapshotManager{dir: dir}
}

// Save writes a snapshot to disk as JSON.
func (m *SnapshotManager) Save(path string, secrets map[string]string) error {
	if err := os.MkdirAll(m.dir, 0700); err != nil {
		return fmt.Errorf("snapshot: mkdir: %w", err)
	}
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Path:      path,
		Secrets:   secrets,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	file := m.filename(path)
	return os.WriteFile(file, data, 0600)
}

// Load reads the most recent snapshot for the given vault path.
func (m *SnapshotManager) Load(path string) (*Snapshot, error) {
	file := m.filename(path)
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("snapshot: read: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &snap, nil
}

func (m *SnapshotManager) filename(path string) string {
	safe := ""
	for _, c := range path {
		if c == '/' || c == '\\' {
			safe += "_"
		} else {
			safe += string(c)
		}
	}
	return fmt.Sprintf("%s/%s.snap.json", m.dir, safe)
}
