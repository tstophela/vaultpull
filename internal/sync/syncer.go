package sync

import (
	"context"
	"time"

	"github.com/user/vaultpull/internal/env"
)

// VaultReader is the interface for reading secrets from Vault.
type VaultReader interface {
	ReadSecrets(ctx context.Context, path string) (map[string]string, error)
}

// Syncer orchestrates reading from Vault and writing to a .env file.
type Syncer struct {
	vault  VaultReader
	writer *env.Writer
	reader *env.Reader
	audit  *env.AuditLogger
}

// New creates a Syncer with the provided dependencies.
func New(v VaultReader, w *env.Writer, r *env.Reader, a *env.AuditLogger) *Syncer {
	return &Syncer{vault: v, writer: w, reader: r, audit: a}
}

// Sync reads secrets from path and merges them into the target file.
func (s *Syncer) Sync(ctx context.Context, path, targetFile string) error {
	incoming, err := s.vault.ReadSecrets(ctx, path)
	if err != nil {
		return err
	}

	existing, _ := s.reader.Read(targetFile)

	diff := env.Diff(existing, incoming)
	merged := env.Merge(existing, incoming)

	if err := s.writer.Write(targetFile, merged); err != nil {
		return err
	}

	if s.audit != nil {
		s.audit.Log(env.AuditEntry{
			Timestamp: time.Now(),
			Path:      path,
			Added:     diff.Added,
			Updated:   diff.Updated,
			Unchanged: diff.Unchanged,
		})
	}
	return nil
}
