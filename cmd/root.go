package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/env"
	"github.com/user/vaultpull/internal/sync"
	"github.com/user/vaultpull/internal/vault"
)

var (
	outputFile string
	backup     bool
	dryRun     bool
)

// vaultClientAdapter wraps vault.Client to match sync.VaultReader interface
type vaultClientAdapter struct {
	client *vault.Client
}

func (a *vaultClientAdapter) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	// Note: vault.Client.ReadSecrets doesn't use context, but we accept it for interface compatibility
	return a.client.ReadSecrets(path)
}

var rootCmd = &cobra.Command{
	Use:   "vaultpull [mount] [path]",
	Short: "Sync HashiCorp Vault secrets into a local .env file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}

		vaultCfg := vault.Config{
			Address: cfg.VaultAddr,
			Token:   cfg.VaultToken,
			Timeout: 0, // use default
		}
		client, err := vault.NewClient(vaultCfg)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		writer := env.NewWriter(outputFile, backup)
		adapter := &vaultClientAdapter{client: client}
		s := sync.New(adapter, writer, nil, nil, nil)

		ctx := cmd.Context()
		return s.Sync(ctx, args[0], args[1], outputFile, os.Stdout, dryRun)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", ".env", "Output .env file path")
	rootCmd.Flags().BoolVarP(&backup, "backup", "b", true, "Backup existing .env before overwriting")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing to disk")
}
