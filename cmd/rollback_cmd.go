package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
)

func init() {
	rollbackCmd := &cobra.Command{
		Use:   "rollback",
		Short: "Manage rollback points for an env file",
	}

	saveCmd := &cobra.Command{
		Use:   "save",
		Short: "Save current env file as a rollback point",
		RunE: func(cmd *cobra.Command, args []string) error {
			envFile, _ := cmd.Flags().GetString("file")
			label, _ := cmd.Flags().GetString("label")
			r := env.NewReader(envFile)
			secrets, err := r.Read()
			if err != nil {
				return fmt.Errorf("read env: %w", err)
			}
			rm := env.NewRollbackManager(rollbackDir(envFile), envFile)
			rp, err := rm.Save(label, secrets)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Saved rollback point %s (label: %s)\n", rp.ID, rp.Label)
			return nil
		},
	}
	saveCmd.Flags().String("file", ".env", "env file to snapshot")
	saveCmd.Flags().String("label", "", "optional label for this rollback point")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available rollback points",
		RunE: func(cmd *cobra.Command, args []string) error {
			envFile, _ := cmd.Flags().GetString("file")
			rm := env.NewRollbackManager(rollbackDir(envFile), envFile)
			points, err := rm.List()
			if err != nil {
				return err
			}
			if len(points) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No rollback points found.")
				return nil
			}
			for _, p := range points {
				fmt.Fprintf(cmd.OutOrStdout(), "%s  %s  %s\n", p.ID, p.CreatedAt.Format("2006-01-02 15:04:05"), p.Label)
			}
			return nil
		},
	}
	listCmd.Flags().String("file", ".env", "env file")

	restoreCmd := &cobra.Command{
		Use:   "restore <id>",
		Short: "Restore env file from a rollback point",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			envFile, _ := cmd.Flags().GetString("file")
			rm := env.NewRollbackManager(rollbackDir(envFile), envFile)
			rp, err := rm.Restore(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Restored rollback point %s (label: %s) to %s\n", rp.ID, rp.Label, envFile)
			return nil
		},
	}
	restoreCmd.Flags().String("file", ".env", "env file to restore")

	rollbackCmd.AddCommand(saveCmd, listCmd, restoreCmd)
	rootCmd.AddCommand(rollbackCmd)
}

func rollbackDir(envFile string) string {
	base := filepath.Dir(envFile)
	return filepath.Join(base, ".vaultpull", "rollbacks", sanitizeRollbackName(filepath.Base(envFile)))
}

func sanitizeRollbackName(name string) string {
	out := make([]byte, len(name))
	for i := range name {
		if os.PathSeparator == name[i] || name[i] == '.' {
			out[i] = '_'
		} else {
			out[i] = name[i]
		}
	}
	return string(out)
}
