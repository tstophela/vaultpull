package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpull/internal/env"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show sync history for an env file",
	RunE:  runHistory,
}

var (
	historyEnvFile string
	historyDir     string
	historyLimit   int
)

func init() {
	historyCmd.Flags().StringVarP(&historyEnvFile, "file", "f", ".env", "env file to show history for")
	historyCmd.Flags().StringVar(&historyDir, "history-dir", ".vaultpull/history", "directory where history is stored")
	historyCmd.Flags().IntVarP(&historyLimit, "limit", "n", 10, "max entries to display (0 = all)")
	rootCmd.AddCommand(historyCmd)
}

func runHistory(cmd *cobra.Command, _ []string) error {
	hm := env.NewHistoryManager(historyDir)
	entries, err := hm.Load(historyEnvFile)
	if err != nil {
		return fmt.Errorf("loading history: %w", err)
	}
	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No history found.")
		return nil
	}
	if historyLimit > 0 && len(entries) > historyLimit {
		entries = entries[len(entries)-historyLimit:]
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tADDED\tUPDATED\tREMOVED")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			len(e.Added),
			len(e.Updated),
			len(e.Removed),
		)
	}
	return w.Flush()
}
