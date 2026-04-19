package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
)

func init() {
	var strategyFlag string
	var prefixes []string
	var outputFlag string

	importCmd := &cobra.Command{
		Use:   "import <source.env>",
		Short: "Import secrets from an existing .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcPath := args[0]

			strategy, err := env.ParseStrategy(strategyFlag)
			if err != nil {
				return fmt.Errorf("invalid strategy: %w", err)
			}

			var filter *env.Filter
			if len(prefixes) > 0 {
				filter = env.NewFilter(prefixes, nil)
			}

			existing := map[string]string{}
			if outputFlag != "" {
				r := env.NewReader(outputFlag)
				if loaded, err := r.Read(); err == nil {
					existing = loaded
				}
			}

			im := env.NewImporter(strategy, filter)
			merged, res, err := im.ImportFile(srcPath, existing)
			if err != nil {
				return err
			}

			if outputFlag != "" {
				w := env.NewWriter(outputFlag)
				if err := w.Write(merged); err != nil {
					return fmt.Errorf("write output: %w", err)
				}
				fmt.Fprintf(os.Stdout, "Imported %d key(s), skipped %d\n", res.Imported, res.Skipped)
			} else {
				keys := make([]string, 0, len(merged))
				for k := range merged {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					fmt.Fprintf(os.Stdout, "%s=%s\n", k, merged[k])
				}
			}
			return nil
		},
	}

	importCmd.Flags().StringVarP(&strategyFlag, "strategy", "s", "overwrite", "Merge strategy: overwrite or preserve")
	importCmd.Flags().StringArrayVarP(&prefixes, "prefix", "p", nil, "Only import keys with these prefixes")
	importCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Target .env file to merge into")

	rootCmd.AddCommand(importCmd)
}
