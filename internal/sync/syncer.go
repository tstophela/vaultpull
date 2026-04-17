package sync

import (
	"context"
	"fmt"
	"io"
	"os"

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
// If backup is true, an existing file is backed up before overwriting.
func (s *Syncer) Sync(ctx context.Context, path, outPath string, backup bool, auditOut io.Writer) error {
	secrets, err := s.vault.ReadSecrets(ctx, path)
	if err != nil {
		return fmt.Errorf("reading secrets: %w", err)
	}

	if s.filter != nil {
		secrets = s.filter.Apply(secrets)
	}

	var existing map[string]string
	if reader, err2 := env.NewReader(outPath); err2 == nil {
		existing, _ = reader.Read()
	}
	if existing == nil {
		existing = map[string]string{}
	}

	merged := env.Merge(existing, secrets)

	if err := s.writer.Write(outPath, merged, backup); err != nil {
		return fmt.Errorf("writing env file: %w", err)
	}

	if s.auditor != nil && auditOut != nil {
		diff := env.Diff(existing, secrets)
		display := secrets
		if s.redactor != nil {
			display = s.redactor.Redact(secrets)
		}
		_ = display
		s.auditor.Log(auditOut, diff)
	}

	fmt.Fprintf(os.Stdout, "synced %d secrets to %s\n", len(merged), outPath)
	return nil
}
