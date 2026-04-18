package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
)

func init() {
	var lockDir string

	lockCmd := &cobra.Command{
		Use:   "lock",
		Short: "Manage .env file locks",
	}

	acquireCmd := &cobra.Command{
		Use:   "acquire <file>",
		Short: "Acquire a lock on a .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, _ := filepath.Abs(args[0])
			lm := env.NewLockManager(lockDir)
			if err := lm.Acquire(target); err != nil {
				return fmt.Errorf("lock acquire: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Locked: %s\n", target)
			return nil
		},
	}

	releaseCmd := &cobra.Command{
		Use:   "release <file>",
		Short: "Release a lock on a .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, _ := filepath.Abs(args[0])
			lm := env.NewLockManager(lockDir)
			if err := lm.Release(target); err != nil {
				return fmt.Errorf("lock release: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Released: %s\n", target)
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status <file>",
		Short: "Check lock status of a .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, _ := filepath.Abs(args[0])
			lm := env.NewLockManager(lockDir)
			if lm.IsLocked(target) {
				fmt.Fprintf(cmd.OutOrStdout(), "LOCKED: %s\n", target)
				os.Exit(1)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "UNLOCKED: %s\n", target)
			return nil
		},
	}

	defaultLockDir := filepath.Join(os.TempDir(), "vaultpull", "locks")
	lockCmd.PersistentFlags().StringVar(&lockDir, "lock-dir", defaultLockDir, "directory to store lock files")
	lockCmd.AddCommand(acquireCmd, releaseCmd, statusCmd)
	rootCmd.AddCommand(lockCmd)
}
