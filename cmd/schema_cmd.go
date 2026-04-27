package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultpull/internal/env"
)

func init() {
	var schemaFile string
	var envFile string
	var applyDefaults bool

	schemaCmd := &cobra.Command{
		Use:   "schema",
		Short: "Validate a .env file against a JSON schema",
		Long: `Loads a JSON schema describing required keys, patterns, and defaults,
then validates the given .env file and reports any violations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sm, err := env.NewSchemaManager(schemaFile)
			if err != nil {
				return fmt.Errorf("loading schema: %w", err)
			}

			reader := env.NewReader(envFile)
			kvs, err := reader.Read()
			if err != nil {
				return fmt.Errorf("reading env file: %w", err)
			}

			if applyDefaults {
				kvs = sm.ApplyDefaults(kvs)
				w := env.NewWriter(envFile)
				if err := w.Write(kvs); err != nil {
					return fmt.Errorf("writing defaults: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), "defaults applied")
			}

			violations := sm.Validate(kvs)
			if len(violations) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "schema validation passed")
				return nil
			}

			for _, v := range violations {
				fmt.Fprintf(cmd.ErrOrStderr(), "  ✗ %s\n", v.Error())
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "%d violation(s) found\n", len(violations))
			os.Exit(1)
			return nil
		},
	}

	schemaCmd.Flags().StringVarP(&schemaFile, "schema", "s", ".env.schema.json", "path to JSON schema file")
	schemaCmd.Flags().StringVarP(&envFile, "file", "f", ".env", "path to .env file to validate")
	schemaCmd.Flags().BoolVar(&applyDefaults, "apply-defaults", false, "write missing default values into the env file before validating")

	rootCmd.AddCommand(schemaCmd)
}
