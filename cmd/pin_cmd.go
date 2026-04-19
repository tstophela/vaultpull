package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"vaultpull/internal/env"
)

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin",
		Short: "Manage pinned secret versions",
	}

	pinDir := pinCmd.PersistentFlags().String("pin-dir", ".vaultpull/pins", "Directory to store pin files")
	envName := pinCmd.PersistentFlags().String("env", "default", "Environment name")

	setCmd := &cobra.Command{
		Use:   "set <key> <version>",
		Short: "Pin a key to a specific secret version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("version must be an integer: %w", err)
			}
			by := os.Getenv("USER")
			if by == "" {
				by = "unknown"
			}
			pm := env.NewPinManager(*pinDir)
			if err := pm.Pin(*envName, args[0], version, by); err != nil {
				return err
			}
			fmt.Printf("Pinned %s to version %d\n", args[0], version)
			return nil
		},
	}

	unsetCmd := &cobra.Command{
		Use:   "unset <key>",
		Short: "Remove a version pin for a key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pm := env.NewPinManager(*pinDir)
			if err := pm.Unpin(*envName, args[0]); err != nil {
				return err
			}
			fmt.Printf("Unpinned %s\n", args[0])
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all pinned keys for an environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			pm := env.NewPinManager(*pinDir)
			pins, err := pm.Load(*envName)
			if err != nil {
				return err
			}
			if len(pins) == 0 {
				fmt.Println("No pins set.")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "KEY\tVERSION\tPINNED BY\tPINNED AT")
			for _, e := range pins {
				fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", e.Key, e.Version, e.PinnedBy, e.PinnedAt.Format("2006-01-02 15:04:05"))
			}
			return w.Flush()
		},
	}

	pinCmd.AddCommand(setCmd, unsetCmd, listCmd)
	rootCmd.AddCommand(pinCmd)
}
