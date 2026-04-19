package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
)

var tagManager = env.NewTagManager()

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags on secret keys",
}

var tagSetCmd = &cobra.Command{
	Use:   "set <secret-key> <tag-key>=<tag-value>",
	Short: "Set a tag on a secret key",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		secretKey := args[0]
		parts := strings.SplitN(args[1], "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("tag must be in key=value format")
		}
		tagManager.Set(secretKey, parts[0], parts[1])
		fmt.Fprintf(cmd.OutOrStdout(), "tagged %s with %s=%s\n", secretKey, parts[0], parts[1])
		return nil
	},
}

var tagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintln(cmd.OutOrStdout(), tagManager.Summary())
		return nil
	},
}

var tagFilterCmd = &cobra.Command{
	Use:   "filter <tag-key>=<tag-value>",
	Short: "Filter secret keys by tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		parts := strings.SplitN(args[0], "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("filter must be in key=value format")
		}
		results := tagManager.FilterByTag(parts[0], parts[1])
		if len(results) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "no matching keys")
			return nil
		}
		for _, k := range results {
			fmt.Fprintln(cmd.OutOrStdout(), k)
		}
		return nil
	},
}

func init() {
	tagCmd.AddCommand(tagSetCmd)
	tagCmd.AddCommand(tagListCmd)
	tagCmd.AddCommand(tagFilterCmd)
	rootCmd.AddCommand(tagCmd)
}
