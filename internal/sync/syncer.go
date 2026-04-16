package sync

import (
	"fmt"
	"log"

	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/vault"
)

// Result holds the outcome of a sync operation.
type Result struct {
	Path       string
	Keyssynced int
	BackedUp   bool
	OutputFile string
}

// Syncer orchestrates reading secrets from Vault and writing them to a .env file.
type Syncer struct {
	client *vault.Client
	writer *env.Writer
}

// New creates a new Syncer.
func New(client *vault.Client, writer *env.Writer) *Syncer {
	return &Syncer{client: client, writer: writer}
}

// Sync reads secrets from the given Vault path and writes them to outFile.
func (s *Syncer) Sync(mountPath, secretPath, outFile string, backup bool) (*Result, error) {
	normalized := vault.NormalizePath(mountPath, secretPath)

	secrets, err := s.client.ReadSecrets(normalized)
	if err != nil {
		return nil, fmt.Errorf("reading secrets from %q: %w", normalized, err)
	}

	if len(secrets) == 0 {
		log.Printf("warning: no secrets found at path %q", normalized)
	}

	bachedUp, err := s.writer.Write(outFile, secrets, backup)
	if err != nil {
		return nil, fmt.Errorf("writing env file %q: %w", outFile, err)
	}

	return &Result{
		Path:       normalized,
		Keyssynced: len(secrets),
		BackedUp:   bachedUp,
		OutputFile: outFile,
	}, nil
}
