package cmd

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourorg/vaultpull/internal/env"
)

func TestPromoteCmd_BasicPromote(t *testing.T) {
	dir := t.TempDir()
	sm := env.NewSnapshotManager(dir)
	require.NoError(t, sm.Save("staging", map[string]string{"KEY": "val"}))

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"promote", "staging", "prod", "--snapshots-dir", dir})
	err := rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "promoted 1")
}

func TestPromoteCmd_DryRun(t *testing.T) {
	dir := t.TempDir()
	sm := env.NewSnapshotManager(dir)
	require.NoError(t, sm.Save("staging", map[string]string{"KEY": "val"}))

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"promote", "staging", "prod", "--dry-run", "--snapshots-dir", dir})
	err := rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "[dry-run]")

	_, err = sm.Load("prod")
	assert.Error(t, err, "prod snapshot should not exist after dry-run")
}

func TestPromoteCmd_SelectedKeys(t *testing.T) {
	dir := t.TempDir()
	sm := env.NewSnapshotManager(dir)
	require.NoError(t, sm.Save("staging", map[string]string{"A": "1", "B": "2"}))

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"promote", "staging", "prod", "--keys", "A", "--snapshots-dir", dir})
	err := rootCmd.Execute()
	require.NoError(t, err)

	dst, err := sm.Load("prod")
	require.NoError(t, err)
	assert.Equal(t, "1", dst["A"])
	_, ok := dst["B"]
	assert.False(t, ok)

	_ = filepath.Join(dir, "prod") // just ensure dir used
}
