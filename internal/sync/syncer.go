package sync

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/vaultpull/internal/env"
)

// VaultReader abstracts reading secrets from Vault.
type VaultReader interface {
	ReadSecrets(ctx context.Context, path string) (map[string]string, error)
}

// Syncer orchestrates pulling secrets from Vault and writing them to a file.
type Syncer struct {
	vault    VaultReader
	writer   *env.Writer
	auditor  *env.AuditLogger
	redactor *env.Redactor
	filter   *env.Filter
}

// New creates a Syncer with the provided dependencies.
func New(v VaultReader, w *env.Writer, a *env.AuditLogger, r *env.Redactor, f *env.Filter) *Syncer {
	return &Syncer{vault: v, writer: w, auditor: a, redactor: r, filter: f}
}

// Sync pulls secrets from path and writes them to outPath.
// If dryRun is true, no changes are written to disk.
func (s *Syncer) Sync(ctx context.Context, mount, path, outPath string, auditOut io.Writer, dryRun bool) error {
	secrets, err := s.vault.ReadSecrets(ctx, path)
	if err != nil {
		return fmt.Errorf("reading secrets: %w", err)
	}

	if s.filter != nil {
		secrets = s.filter.Apply(secrets)
	}

	var existing map[string]string
	reader := env.NewReader(outPath)
	var err2 error
	existing, err2 = reader.Read()
	if err2 != nil {
		return fmt.Errorf("reading existing env: %w", err2)
	}
	if existing == nil {
		existing = map[string]string{}
	}

	merged := env.Merge(existing, secrets)
	diff := env.Diff(existing, secrets)

	if !dryRun {
		if err := s.writer.Write(merged); err != nil {
			return fmt.Errorf("writing env file: %w", err)
		}
	}

	if s.auditor != nil {
		added := make([]string, 0, len(diff.Added))
		for k := range diff.Added {
			added = append(added, k)
		}
		updated := make([]string, 0, len(diff.Updated))
		for k := range diff.Updated {
			updated = append(updated, k)
		}
		unchanged := make([]string, 0, len(diff.Unchanged))
		for k := range diff.Unchanged {
			unchanged = append(unchanged, k)
		}
		entry := env.AuditEntry{
			Timestamp: time.Now(),
			Path:      path,
			Added:     added,
			Updated:   updated,
			Unchanged: unchanged,
		}
		s.auditor.Log(entry)
	}

	if dryRun {
		out := auditOut
		if out == nil {
			out = os.Stdout
		}
		fmt.Fprintf(out, "[DRY-RUN] Would sync %d secrets to %s\n", len(merged), outPath)
		s.printDiff(out, diff)
	} else {
		fmt.Fprintf(os.Stdout, "synced %d secrets to %s\n", len(merged), outPath)
	}
	return nil
}

// printDiff displays the diff result in a human-readable format.
func (s *Syncer) printDiff(w io.Writer, diff env.DiffResult) {
	if len(diff.Added) > 0 {
		fmt.Fprintf(w, "\n  Added (%d):\n", len(diff.Added))
		for k := range diff.Added {
			fmt.Fprintf(w, "    + %s\n", k)
		}
	}
	if len(diff.Updated) > 0 {
		fmt.Fprintf(w, "\n  Updated (%d):\n", len(diff.Updated))
		for k := range diff.Updated {
			fmt.Fprintf(w, "    ~ %s\n", k)
		}
	}
	if len(diff.Unchanged) > 0 {
		fmt.Fprintf(w, "\n  Unchanged (%d):\n", len(diff.Unchanged))
		for k := range diff.Unchanged {
			fmt.Fprintf(w, "    = %s\n", k)
		}
	}
}
