package env

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeWatchFile(t *testing.T, dir, content string) string {
	t.Helper()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeWatchFile: %v", err)
	}
	return p
}

func TestWatcher_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	path := writeWatchFile(t, dir, "KEY=old\n")

	w := NewWatcher(path, 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch, err := w.Watch(ctx)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	// Modify file after a short delay.
	time.AfterFunc(60*time.Millisecond, func() {
		_ = os.WriteFile(path, []byte("KEY=new\n"), 0600)
	})

	select {
	case ev := <-ch:
		if ev.Path != path {
			t.Errorf("expected path %q, got %q", path, ev.Path)
		}
		if ev.OldHash == ev.NewHash {
			t.Error("expected hashes to differ")
		}
		if ev.At.IsZero() {
			t.Error("expected non-zero timestamp")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatcher_NoEventWhenUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := writeWatchFile(t, dir, "KEY=same\n")

	w := NewWatcher(path, 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ch, err := w.Watch(ctx)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	select {
	case ev := <-ch:
		t.Errorf("unexpected event: %+v", ev)
	case <-ctx.Done():
		// expected: no events fired
	}
}

func TestWatcher_InitialHashError(t *testing.T) {
	w := NewWatcher("/nonexistent/.env", 50*time.Millisecond)
	_, err := w.Watch(context.Background())
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestWatcher_ChannelClosedOnCancel(t *testing.T) {
	dir := t.TempDir()
	path := writeWatchFile(t, dir, "A=1\n")

	w := NewWatcher(path, 20*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())

	ch, err := w.Watch(ctx)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}
	cancel()

	// Channel should close shortly after cancel.
	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return // closed as expected
			}
		case <-timer.C:
			t.Fatal("channel was not closed after context cancel")
		}
	}
}
