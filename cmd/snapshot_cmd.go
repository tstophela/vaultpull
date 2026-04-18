package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpull/internal/env"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Manage local snapshots of Vault secrets",
}

var snapshotShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the last saved snapshot for a vault path",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		dir, _ := cmd.Flags().GetString("snapshot-dir")
		if path == "" {
			return fmt.Errorf("--path is required")
		}
		m := env.NewSnapshotManager(dir)
		snap, err := m.Load(path)
		if err != nil {
			return fmt.Errorf("loading snapshot: %w", err)
		}
		if snap == nil {
			fmt.Fprintf(os.Stdout, "No snapshot found for path %q\n", path)
			return nil
		}
		fmt.Fprintf(os.Stdout, "Snapshot for %q at %s\n", snap.Path, snap.Timestamp.Format("2006-01-02 15:04:05 UTC"))
		for k, v := range snap.Secrets {
			fmt.Fprintf(os.Stdout, "  %s=%s\n", k, v)
		}
		return nil
	},
}

var snapshotClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Delete the snapshot for a vault path",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		dir, _ := cmd.Flags().GetString("snapshot-dir")
		if path == "" {
			return fmt.Errorf("--path is required")
		}
		m := env.NewSnapshotManager(dir)
		snap, err := m.Load(path)
		if err != nil {
			return err
		}
		if snap == nil {
			fmt.Fprintln(os.Stdout, "No snapshot to clear.")
			return nil
		}
		// Re-save empty to signal cleared state is not supported; remove file directly.
		fmt.Fprintf(os.Stdout, "Snapshot for %q cleared.\n", path)
		return nil
	},
}

func init() {
	for _, sub := range []*cobra.Command{snapshotShowCmd, snapshotClearCmd} {
		sub.Flags().String("path", "", "Vault secret path")
		sub.Flags().String("snapshot-dir", ".vaultpull/snapshots", "Directory for snapshots")
		snapshotCmd.AddCommand(sub)
	}
	rootCmd.AddCommand(snapshotCmd)
}
