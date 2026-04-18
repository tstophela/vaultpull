package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
)

var (
	encryptPassphrase string
	encryptDecrypt    bool
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt [value]",
	Short: "Encrypt or decrypt a secret value using AES-GCM",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if encryptPassphrase == "" {
			encryptPassphrase = os.Getenv("VAULTPULL_ENCRYPT_KEY")
		}
		if encryptPassphrase == "" {
			return fmt.Errorf("passphrase required: use --passphrase or VAULTPULL_ENCRYPT_KEY")
		}

		e := env.NewEncryptor(encryptPassphrase)
		value := args[0]

		if encryptDecrypt {
			plain, err := e.Decrypt(value)
			if err != nil {
				return fmt.Errorf("decrypt: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), plain)
			return nil
		}

		encoded, err := e.Encrypt(value)
		if err != nil {
			return fmt.Errorf("encrypt: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), encoded)
		return nil
	},
}

func init() {
	encryptCmd.Flags().StringVar(&encryptPassphrase, "passphrase", "", "Encryption passphrase (or set VAULTPULL_ENCRYPT_KEY)")
	encryptCmd.Flags().BoolVar(&encryptDecrypt, "decrypt", false, "Decrypt the value instead of encrypting")
	rootCmd.AddCommand(encryptCmd)
}
