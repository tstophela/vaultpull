package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpull/internal/env"
)

var validateWarnEmpty bool

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate keys in a local .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		reader := env.NewReader(filePath)
		secrets, err := reader.Read()
		if err != nil {
			return fmt.Errorf("reading %s: %w", filePath, err)
		}

		validator := env.NewValidator(validateWarnEmpty)
		if err := validator.Validate(secrets); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return fmt.Errorf("validation errors found")
		}

		fmt.Printf("✓ %s is valid (%d keys)\n", filePath, len(secrets))
		return nil
	},
}

func init() {
	validateCmd.Flags().BoolVar(&validateWarnEmpty, "warn-empty", false, "warn on keys with empty values")
	rootCmd.AddCommand(validateCmd)
}
