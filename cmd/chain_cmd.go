package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
)

func init() {
	var files []string
	var showSource bool
	var key string

	chainCmd := &cobra.Command{
		Use:   "chain",
		Short: "Resolve env vars from an ordered chain of .env files",
		Long: `Resolves environment variables from multiple .env files in priority order.
The first file that defines a key wins. Use --source to see which file each key came from.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(files) == 0 {
				return fmt.Errorf("at least one --file is required")
			}

			chain := make([]env.ChainEntry, 0, len(files))
			for _, f := range files {
				r := env.NewReader(f)
				vals, err := r.Read()
				if err != nil {
					fmt.Fprintf(os.Stderr, "warning: skipping %s: %v\n", f, err)
					continue
				}
				chain = append(chain, env.ChainEntry{Name: f, Values: vals})
			}

			resolver := env.NewResolver(chain)

			if key != "" {
				cr, err := resolver.ResolveKey(key)
				if err != nil {
					return err
				}
				if showSource {
					fmt.Printf("%s=%s  # from %s\n", cr.Key, cr.Value, cr.Source)
				} else {
					fmt.Printf("%s=%s\n", cr.Key, cr.Value)
				}
				return nil
			}

			for _, cr := range resolver.Resolve() {
				if showSource {
					fmt.Printf("%s=%s  # from %s\n", cr.Key, cr.Value, cr.Source)
				} else {
					fmt.Printf("%s=%s\n", cr.Key, cr.Value)
				}
			}
			return nil
		},
	}

	chainCmd.Flags().StringArrayVarP(&files, "file", "f", nil, "env files in priority order (first wins)")
	chainCmd.Flags().BoolVar(&showSource, "source", false, "annotate each key with its source file")
	chainCmd.Flags().StringVarP(&key, "key", "k", "", "resolve a single key")

	rootCmd.AddCommand(chainCmd)
}
