package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
)

func init() {
	var rulesFlag []string

	lintCmd := &cobra.Command{
		Use:   "lint [file]",
		Short: "Lint a .env file for common issues",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

			reader := env.NewReader()
			envMap, err := reader.Read(filePath)
			if err != nil {
				return fmt.Errorf("reading env file: %w", err)
			}

			var rules []env.LintRule
			for _, r := range rulesFlag {
				rules = append(rules, env.LintRule(strings.TrimSpace(r)))
			}

			linter := env.NewLinter(rules...)
			issues := linter.Lint(envMap)

			if len(issues) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "✔ No lint issues found.")
				return nil
			}

			for _, issue := range issues {
				fmt.Fprintln(cmd.OutOrStdout(), issue.String())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\n%d issue(s) found.\n", len(issues))
			os.Exit(1)
			return nil
		},
	}

	lintCmd.Flags().StringSliceVar(
		&rulesFlag,
		"rules",
		nil,
		"comma-separated list of rules to enable (no_lowercase,no_spaces,no_duplicates,no_empty_keys)",
	)

	rootCmd.AddCommand(lintCmd)
}
