package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourorg/vaultpull/internal/env"
)

func init() {
	var (
		keys      string
		overwrite bool
		dryRun    bool
		snapsDir  string
	)

	cmd := &cobra.Command{
		Use:   "promote <source-env> <target-env>",
		Short: "Promote secrets from one environment snapshot to another",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, dst := args[0], args[1]

			var keyList []string
			if keys != "" {
				for _, k := range strings.Split(keys, ",") {
					if k = strings.TrimSpace(k); k != "" {
						keyList = append(keyList, k)
					}
				}
			}

			p := env.NewPromoter(snapsDir)
			res, err := p.Promote(env.PromoteOptions{
				SourceEnv: src,
				TargetEnv: dst,
				Keys:      keyList,
				DryRun:    dryRun,
				Overwrite: overwrite,
			})
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), res.Summary())
			if len(res.Promoted) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "  promoted:", strings.Join(res.Promoted, ", "))
			}
			if len(res.Skipped) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "  skipped: ", strings.Join(res.Skipped, ", "))
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&keys, "keys", "", "Comma-separated list of keys to promote (default: all)")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing keys in target")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be promoted without writing")
	cmd.Flags().StringVar(&snapsDir, "snapshots-dir", ".vaultpull/snapshots", "Directory storing snapshots")

	rootCmd.AddCommand(cmd)
}
