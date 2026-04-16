package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/sync"
	"github.com/your-org/vaultpull/internal/vault"
)

var (
	outputFile string
	backup     bool
)

var rootCmd = &cobra.Command{
	Use:   "vaultpull [mount] [path]",
	Short: "Sync HashiCorp Vault secrets into a local .env file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}

		client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		writer := env.NewWriter()
		s := sync.New(client, writer)

		result, err := s.Sync(args[0], args[1], outputFile, backup)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Synced %d keys from %q → %s\n", result.Keyssynced, result.Path, result.OutputFile)
		if result.BackedUp {
			fmt.Printf("  Backup created for previous %s\n", result.OutputFile)
		}
		return nil
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
}
