package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
)

func init() {
	expireCmd := &cobra.Command{
		Use:   "expire",
		Short: "Manage secret expiry policies",
	}

	// expire set <key> <duration>
	setCmd := &cobra.Command{
		Use:   "set <key> <duration>",
		Short: "Set an expiry duration for a secret key (e.g. 24h, 7d)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			durationStr := args[1]

			d, err := time.ParseDuration(durationStr)
			if err != nil {
				return fmt.Errorf("invalid duration %q: %w", durationStr, err)
			}

			dataDir, _ := cmd.Flags().GetString("data-dir")
			mgr := env.NewExpiryManager(dataDir)
			if err := mgr.Set(key, d); err != nil {
				return fmt.Errorf("failed to set expiry: %w", err)
			}

			fmt.Fprintf(os.Stdout, "Expiry set: %s expires in %s\n", key, d)
			return nil
		},
	}
	setCmd.Flags().String("data-dir", ".vaultpull", "Directory for expiry metadata")

	// expire get <key>
	getCmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Get the expiry status of a secret key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]

			dataDir, _ := cmd.Flags().GetString("data-dir")
			mgr := env.NewExpiryManager(dataDir)

			entry, err := mgr.Get(key)
			if err != nil {
				return fmt.Errorf("failed to get expiry: %w", err)
			}

			expired := mgr.IsExpired(key)
			status := "valid"
			if expired {
				status = "EXPIRED"
			}

			fmt.Fprintf(os.Stdout, "Key:     %s\nExpires: %s\nStatus:  %s\n",
				key, entry.ExpiresAt.Format(time.RFC3339), status)
			return nil
		},
	}
	getCmd.Flags().String("data-dir", ".vaultpull", "Directory for expiry metadata")

	// expire check — print all expired keys
	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "List all expired secret keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			dataDir, _ := cmd.Flags().GetString("data-dir")
			mgr := env.NewExpiryManager(dataDir)

			expired, err := mgr.ListExpired()
			if err != nil {
				return fmt.Errorf("failed to list expired keys: %w", err)
			}

			if len(expired) == 0 {
				fmt.Fprintln(os.Stdout, "No expired keys found.")
				return nil
			}

			fmt.Fprintln(os.Stdout, "Expired keys:")
			for _, k := range expired {
				fmt.Fprintf(os.Stdout, "  - %s\n", k)
			}
			return nil
		},
	}
	checkCmd.Flags().String("data-dir", ".vaultpull", "Directory for expiry metadata")

	expireCmd.AddCommand(setCmd, getCmd, checkCmd)
	rootCmd.AddCommand(expireCmd)
}
