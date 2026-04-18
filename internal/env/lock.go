package env

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LockManager prevents concurrent writes to the same .env file.
type LockManager struct {
	dir string
}

// NewLockManager creates a LockManager that stores lock files in dir.
func NewLockManager(dir string) *LockManager {
	return &LockManager{dir: dir}
}

func (lm *LockManager) lockPath(target string) string {
	base := filepath.Base(target)
	return filepath.Join(lm.dir, base+".lock")
}

// Acquire creates a lock file for target. Returns error if already locked.
func (lm *LockManager) Acquire(target string) error {
	if err := os.MkdirAll(lm.dir, 0700); err != nil {
		return fmt.Errorf("lock dir: %w", err)
	}
	lp := lm.lockPath(target)
	f, err := os.OpenFile(lp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("file %q is locked by another process", target)
		}
		return fmt.Errorf("acquire lock: %w", err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%d\n%s\n", os.Getpid(), time.Now().Format(time.RFC3339))
	return err
}

// Release removes the lock file for target.
func (lm *LockManager) Release(target string) error {
	lp := lm.lockPath(target)
	if err := os.Remove(lp); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("release lock: %w", err)
	}
	return nil
}

// IsLocked reports whether target is currently locked.
func (lm *LockManager) IsLocked(target string) bool {
	_, err := os.Stat(lm.lockPath(target))
	return err == nil
}
