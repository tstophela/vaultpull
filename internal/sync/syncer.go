package sync

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/user/vaultpull/internal/env"
)

// VaultReader reads secrets from Vault.
type VaultReader interface {
	ReadSecrets(ctx context.Context, path string) (map[string]string, error)
}

// Syncer orchestrates pulling secrets from Vault and writing them to a file.
type Syncer struct {
	vault  VaultReader
	writer *env.Writer
	filter *env.Filter
	audit  *env.AuditLogger
}

// New creates a Syncer with the provided dependencies.
func New(v VaultReader, w *env.Writer, f *env.Filter, audit *env.AuditLogger) *Syncer {
	return &Syncer{vault: v, writer: w, filter: f, audit: audit}
}

// Sync pulls secrets from the given Vault path and writes them to outPath.
// If backup is true, an existing file at outPath is preserved as outPath.bak.
func (s *Syncer) Sync(ctx context.Context, vaultPath, outPath string, backup bool) error {
	secrets, err := s.vault.ReadSecrets(ctx, vaultPath)
	if err != nil {
		return fmt.Errorf("vault read: %w", err)
	}

	if s.filter != nil {
		secrets = s.filter.Apply(secrets)
	}

	var existing map[string]string
	if r, err := env.NewReader(outPath); err == nil {
		existing, _ = r.Read()
	}

	diff := env.Diff(existing, secrets)
	merged := env.Merge(existing, secrets)

	if err := s.writer.Write(outPath, merged, backup); err != nil {
		return fmt.Errorf("write env: %w", err)
	}

	if s.audit != nil {
		var w io.Writer = os.Stdout
		s.audit.Log(w, diff)
	}

	return nil
}
