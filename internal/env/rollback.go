package env

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// RollbackPoint represents a saved state that can be restored.
type RollbackPoint struct {
	ID        string
	CreatedAt time.Time
	Label     string
	Secrets   map[string]string
}

// RollbackManager manages rollback points for an env file.
type RollbackManager struct {
	dir      string
	envFile  string
	now      func() time.Time
}

// NewRollbackManager creates a RollbackManager storing points under dir.
func NewRollbackManager(dir, envFile string) *RollbackManager {
	return &RollbackManager{
		dir:     dir,
		envFile: envFile,
		now:     time.Now,
	}
}

// Save creates a new rollback point from the given secrets map.
func (r *RollbackManager) Save(label string, secrets map[string]string) (*RollbackPoint, error) {
	if err := os.MkdirAll(r.dir, 0700); err != nil {
		return nil, fmt.Errorf("rollback: create dir: %w", err)
	}
	ts := r.now()
	id := fmt.Sprintf("%d", ts.UnixNano())
	rp := &RollbackPoint{ID: id, CreatedAt: ts, Label: label, Secrets: secrets}
	path := filepath.Join(r.dir, id+".json")
	if err := writeJSON(path, rp); err != nil {
		return nil, fmt.Errorf("rollback: save: %w", err)
	}
	return rp, nil
}

// List returns all rollback points sorted newest-first.
func (r *RollbackManager) List() ([]*RollbackPoint, error) {
	entries, err := os.ReadDir(r.dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("rollback: list: %w", err)
	}
	var points []*RollbackPoint
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		var rp RollbackPoint
		if err := readJSON(filepath.Join(r.dir, e.Name()), &rp); err != nil {
			continue
		}
		points = append(points, &rp)
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].CreatedAt.After(points[j].CreatedAt)
	})
	return points, nil
}

// Restore writes the secrets from the given rollback point ID to the env file.
func (r *RollbackManager) Restore(id string) (*RollbackPoint, error) {
	path := filepath.Join(r.dir, id+".json")
	var rp RollbackPoint
	if err := readJSON(path, &rp); err != nil {
		return nil, fmt.Errorf("rollback: load point %s: %w", id, err)
	}
	w := NewWriter(r.envFile)
	if err := w.Write(rp.Secrets); err != nil {
		return nil, fmt.Errorf("rollback: restore write: %w", err)
	}
	return &rp, nil
}
