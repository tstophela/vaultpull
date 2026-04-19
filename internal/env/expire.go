package env

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ExpiryRecord holds expiry metadata for a secret path.
type ExpiryRecord struct {
	Path      string    `json:"path"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// IsExpired returns true if the record's expiry time has passed.
func (r ExpiryRecord) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// ExpiryManager manages TTL records for secret paths.
type ExpiryManager struct {
	dir string
}

// NewExpiryManager creates an ExpiryManager storing records under dir.
func NewExpiryManager(dir string) *ExpiryManager {
	return &ExpiryManager{dir: dir}
}

func (m *ExpiryManager) filename(path string) string {
	return filepath.Join(m.dir, sanitizeFilename(path)+".expiry.json")
}

// Set stores an expiry record for the given secret path with a TTL duration.
func (m *ExpiryManager) Set(path string, ttl time.Duration) error {
	if err := os.MkdirAll(m.dir, 0700); err != nil {
		return fmt.Errorf("expiry: mkdir: %w", err)
	}
	rec := ExpiryRecord{
		Path:      path,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
	}
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.filename(path), data, 0600)
}

// Get retrieves the expiry record for a path. Returns error if not found.
func (m *ExpiryManager) Get(path string) (ExpiryRecord, error) {
	data, err := os.ReadFile(m.filename(path))
	if err != nil {
		return ExpiryRecord{}, fmt.Errorf("expiry: not found: %w", err)
	}
	var rec ExpiryRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return ExpiryRecord{}, fmt.Errorf("expiry: parse: %w", err)
	}
	return rec, nil
}

// Delete removes the expiry record for a path.
func (m *ExpiryManager) Delete(path string) error {
	err := os.Remove(m.filename(path))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
