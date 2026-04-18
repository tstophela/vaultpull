package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLockManager_AcquireAndRelease(t *testing.T) {
	dir := t.TempDir()
	lm := NewLockManager(dir)
	target := filepath.Join(dir, ".env")

	if err := lm.Acquire(target); err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	if !lm.IsLocked(target) {
		t.Fatal("expected locked")
	}
	if err := lm.Release(target); err != nil {
		t.Fatalf("Release: %v", err)
	}
	if lm.IsLocked(target) {
		t.Fatal("expected unlocked after release")
	}
}

func TestLockManager_DoubleAcquireFails(t *testing.T) {
	dir := t.TempDir()
	lm := NewLockManager(dir)
	target := filepath.Join(dir, ".env")

	if err := lm.Acquire(target); err != nil {
		t.Fatalf("first Acquire: %v", err)
	}
	defer lm.Release(target) //nolint

	if err := lm.Acquire(target); err == nil {
		t.Fatal("expected error on double acquire")
	}
}

func TestLockManager_ReleaseNonExistent(t *testing.T) {
	dir := t.TempDir()
	lm := NewLockManager(dir)
	target := filepath.Join(dir, ".env")

	if err := lm.Release(target); err != nil {
		t.Fatalf("Release non-existent should not error: %v", err)
	}
}

func TestLockManager_CreatesLockDir(t *testing.T) {
	base := t.TempDir()
	lockDir := filepath.Join(base, "locks", "nested")
	lm := NewLockManager(lockDir)
	target := filepath.Join(base, ".env")

	if err := lm.Acquire(target); err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	if _, err := os.Stat(lockDir); err != nil {
		t.Fatalf("lock dir not created: %v", err)
	}
	lm.Release(target) //nolint
}
