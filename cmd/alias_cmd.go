package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpull/internal/env"
)

func init() {
	aliasDir := os.ExpandEnv("${HOME}/.vaultpull")

	aliasCmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage short aliases for Vault secret paths",
	}

	// alias set <name> <path>
	setCmd := &cobra.Command{
		Use:   "set <name> <path>",
		Short: "Register or update an alias",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			am, err := env.NewAliasManager(aliasDir)
			if err != nil {
				return err
			}
			if err := am.Set(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "alias %q -> %q saved\n", args[0], args[1])
			return nil
		},
	}

	// alias get <name>
	getCmd := &cobra.Command{
		Use:   "get <name>",
		Short: "Resolve an alias to its Vault path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			am, err := env.NewAliasManager(aliasDir)
			if err != nil {
				return err
			}
			path, ok := am.Get(args[0])
			if !ok {
				return fmt.Errorf("alias %q not found", args[0])
			}
			fmt.Fprintln(cmd.OutOrStdout(), path)
			return nil
		},
	}

	// alias delete <name>
	delCmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Remove an alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			am, err := env.NewAliasManager(aliasDir)
			if err != nil {
				return err
			}
			if err := am.Delete(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "alias %q removed\n", args[0])
			return nil
		},
	}

	// alias list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all registered aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			am, err := env.NewAliasManager(aliasDir)
			if err != nil {
				return err
			}
			entries := am.List()
			if len(entries) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no aliases defined")
				return nil
			}
			for _, e := range entries {
				fmt.Fprintf(cmd.OutOrStdout(), "%-20s %s\n", e.Name, e.Path)
			}
			return nil
		},
	}

	aliasCmd.AddCommand(setCmd, getCmd, delCmd, listCmd)
	rootCmd.AddCommand(aliasCmd)
}
