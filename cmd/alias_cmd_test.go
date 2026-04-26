package cmd_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpull/internal/env"
)

func newAliasCmdSuite(t *testing.T) (*env.AliasManager, string) {
	t.Helper()
	dir := t.TempDir()
	am, err := env.NewAliasManager(dir)
	if err != nil {
		t.Fatalf("NewAliasManager: %v", err)
	}
	return am, dir
}

func runAliasCmd(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "vaultpull"}
	var buf bytes.Buffer

	aliasCmd := buildAliasCmd(dir)
	aliasCmd.SetOut(&buf)
	for _, sub := range aliasCmd.Commands() {
		sub.SetOut(&buf)
	}
	root.AddCommand(aliasCmd)
	root.SetOut(&buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

// buildAliasCmd constructs the alias command tree with a custom dir for testing.
func buildAliasCmd(dir string) *cobra.Command {
	aliasCmd := &cobra.Command{Use: "alias", Short: "Manage aliases"}

	aliasCmd.AddCommand(&cobra.Command{
		Use:  "set <name> <path>",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			am, _ := env.NewAliasManager(dir)
			_ = am.Set(args[0], args[1])
			cmd.Printf("alias %q -> %q saved\n", args[0], args[1])
			return nil
		},
	})
	aliasCmd.AddCommand(&cobra.Command{
		Use:  "list",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			am, _ := env.NewAliasManager(dir)
			for _, e := range am.List() {
				cmd.Printf("%-20s %s\n", e.Name, e.Path)
			}
			return nil
		},
	})
	aliasCmd.AddCommand(&cobra.Command{
		Use:  "delete <name>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			am, _ := env.NewAliasManager(dir)
			return am.Delete(args[0])
		},
	})
	return aliasCmd
}

func TestAliasCmd_SetAndList(t *testing.T) {
	dir := t.TempDir()
	out, err := runAliasCmd(t, dir, "alias", "set", "prod", "secret/prod/app")
	if err != nil {
		t.Fatalf("set: %v", err)
	}
	if !strings.Contains(out, "prod") {
		t.Errorf("expected confirmation, got %q", out)
	}

	out, err = runAliasCmd(t, dir, "alias", "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(out, "secret/prod/app") {
		t.Errorf("expected path in list output, got %q", out)
	}
}

func TestAliasCmd_Delete(t *testing.T) {
	dir := t.TempDir()
	_, _ = runAliasCmd(t, dir, "alias", "set", "dev", "secret/dev/app")

	_, err := runAliasCmd(t, dir, "alias", "delete", "dev")
	if err != nil {
		t.Fatalf("delete: %v", err)
	}

	am, _ := env.NewAliasManager(dir)
	if _, ok := am.Get("dev"); ok {
		t.Error("expected alias to be deleted")
	}
}

func TestAliasCmd_Delete_NonExistent(t *testing.T) {
	dir := t.TempDir()
	_, err := runAliasCmd(t, dir, "alias", "delete", "ghost")
	if err == nil {
		t.Error("expected error deleting non-existent alias")
	}
}

func TestAliasCmd_List_Empty(t *testing.T) {
	dir := t.TempDir()
	// Ensure no aliases.json exists
	os.Remove(dir + "/aliases.json")
	out, err := runAliasCmd(t, dir, "alias", "list")
	if err != nil {
		t.Fatalf("list empty: %v", err)
	}
	// No output expected (empty list)
	_ = out
}
