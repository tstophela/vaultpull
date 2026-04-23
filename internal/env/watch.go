package env

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchEvent describes a change detected in a watched env file.
type WatchEvent struct {
	Path    string
	OldHash string
	NewHash string
	At      time.Time
}

// Watcher monitors an env file for changes by polling its content hash.
type Watcher struct {
	path     string
	interval time.Duration
	lastHash string
	now      func() time.Time
}

// NewWatcher creates a Watcher for the given file path and poll interval.
func NewWatcher(path string, interval time.Duration) *Watcher {
	return &Watcher{
		path:     path,
		interval: interval,
		now:      time.Now,
	}
}

// Watch starts polling the file and sends WatchEvents on the returned channel.
// It stops when ctx is cancelled. The channel is closed on exit.
func (w *Watcher) Watch(ctx context.Context) (<-chan WatchEvent, error) {
	h, err := w.hashFile()
	if err != nil {
		return nil, fmt.Errorf("watch: initial hash failed: %w", err)
	}
	w.lastHash = h

	ch := make(chan WatchEvent, 4)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				current, err := w.hashFile()
				if err != nil {
					continue
				}
				if current != w.lastHash {
					ch <- WatchEvent{
						Path:    w.path,
						OldHash: w.lastHash,
						NewHash: current,
						At:      w.now(),
					}
					w.lastHash = current
				}
			}
		}
	}()
	return ch, nil
}

func (w *Watcher) hashFile() (string, error) {
	f, err := os.Open(w.path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
