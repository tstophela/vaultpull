package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
)

func init() {
	var fileA, fileB string

	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare two .env files and show differences",
		RunE: func(cmd *cobra.Command, args []string) error {
			readerA := env.NewReader(fileA)
			oldMap, err := readerA.Read()
			if err != nil {
				return fmt.Errorf("reading %s: %w", fileA, err)
			}

			readerB := env.NewReader(fileB)
			newMap, err := readerB.Read()
			if err != nil {
				return fmt.Errorf("reading %s: %w", fileB, err)
			}

			result := env.Compare(oldMap, newMap)

			if !result.HasChanges() {
				fmt.Fprintln(cmd.OutOrStdout(), "No differences found.")
				return nil
			}

			w := cmd.OutOrStdout()
			for _, k := range result.SortedAdded() {
				fmt.Fprintf(w, "+ %s=%s\n", k, result.Added[k])
			}
			for _, k := range result.SortedRemoved() {
				fmt.Fprintf(w, "- %s=%s\n", k, result.Removed[k])
			}
			for _, k := range result.SortedChanged() {
				pair := result.Changed[k]
				fmt.Fprintf(w, "~ %s: %s -> %s\n", k, pair[0], pair[1])
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&fileA, "from", "f", ".env", "Base env file")
	cmd.Flags().StringVarP(&fileB, "to", "t", ".env.new", "Target env file")

	if err := cmd.MarkFlagRequired("from"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	rootCmd.AddCommand(cmd)
}
