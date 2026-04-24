package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
)

func init() {
	scopeCmd := &cobra.Command{
		Use:   "scope",
		Short: "Manage named environment scopes",
	}

	indexPath := filepath.Join(".vaultpull", "scopes.json")

	// scope register
	registerCmd := &cobra.Command{
		Use:   "register <name> <path>",
		Short: "Register a named scope pointing to an env file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			m := env.NewScopeManager(indexPath)
			s := env.Scope{Name: args[0], Path: args[1]}
			if err := m.Register(s); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "scope %q registered -> %s\n", args[0], args[1])
			return nil
		},
	}

	// scope list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all registered scopes",
		RunE: func(cmd *cobra.Command, _ []string) error {
			m := env.NewScopeManager(indexPath)
			scopes, err := m.List()
			if err != nil {
				return err
			}
			if len(scopes) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no scopes registered")
				return nil
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tPATH")
			for _, s := range scopes {
				fmt.Fprintf(w, "%s\t%s\n", s.Name, s.Path)
			}
			return w.Flush()
		},
	}

	// scope remove
	removeCmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a registered scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			m := env.NewScopeManager(indexPath)
			if err := m.Remove(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "scope %q removed\n", args[0])
			return nil
		},
	}

	// scope use — print the path for a scope (useful in shell scripts)
	useCmd := &cobra.Command{
		Use:   "use <name>",
		Short: "Print the env file path for a scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			m := env.NewScopeManager(indexPath)
			s, ok, err := m.Get(args[0])
			if err != nil {
				return err
			}
			if !ok {
				fmt.Fprintf(os.Stderr, "scope %q not found\n", args[0])
				os.Exit(1)
			}
			fmt.Fprintln(cmd.OutOrStdout(), s.Path)
			return nil
		},
	}

	scopeCmd.AddCommand(registerCmd, listCmd, removeCmd, useCmd)
	rootCmd.AddCommand(scopeCmd)
}
