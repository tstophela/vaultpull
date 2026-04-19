package env

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// PinEntry records a pinned secret version for a given key.
type PinEntry struct {
	Key       string    `json:"key"`
	Version   int       `json:"version"`
	PinnedAt  time.Time `json:"pinned_at"`
	PinnedBy  string    `json:"pinned_by"`
}

// PinManager manages pinned secret versions.
type PinManager struct {
	dir string
}

func NewPinManager(dir string) *PinManager {
	return &PinManager{dir: dir}
}

func (p *PinManager) pinFile(env string) string {
	return filepath.Join(p.dir, fmt.Sprintf("%s.pins.json", env))
}

// Pin sets a version pin for a key in the given env.
func (p *PinManager) Pin(env, key string, version int, by string) error {
	pins, _ := p.Load(env)
	pins[key] = PinEntry{Key: key, Version: version, PinnedAt: time.Now().UTC(), PinnedBy: by}
	return p.save(env, pins)
}

// Unpin removes a pin for a key.
func (p *PinManager) Unpin(env, key string) error {
	pins, err := p.Load(env)
	if err != nil {
		return err
	}
	delete(pins, key)
	return p.save(env, pins)
}

// Load returns all pins for the given env.
func (p *PinManager) Load(env string) (map[string]PinEntry, error) {
	data, err := os.ReadFile(p.pinFile(env))
	if os.IsNotExist(err) {
		return map[string]PinEntry{}, nil
	}
	if err != nil {
		return nil, err
	}
	var pins map[string]PinEntry
	if err := json.Unmarshal(data, &pins); err != nil {
		return nil, err
	}
	return pins, nil
}

func (p *PinManager) save(env string, pins map[string]PinEntry) error {
	if err := os.MkdirAll(p.dir, 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(pins, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p.pinFile(env), data, 0600)
}
