package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestPromoter(t *testing.T) *Promoter {
	t.Helper()
	return NewPromoter(t.TempDir())
}

func TestPromoter_BasicPromote(t *testing.T) {
	p := newTestPromoter(t)
	require.NoError(t, p.snapshots.Save("staging", map[string]string{"FOO": "bar", "BAZ": "qux"}))

	res, err := p.Promote(PromoteOptions{SourceEnv: "staging", TargetEnv: "prod"})
	require.NoError(t, err)
	assert.Len(t, res.Promoted, 2)
	assert.Empty(t, res.Skipped)
	assert.False(t, res.DryRun)

	dst, err := p.snapshots.Load("prod")
	require.NoError(t, err)
	assert.Equal(t, "bar", dst["FOO"])
}

func TestPromoter_SkipsExistingWithoutOverwrite(t *testing.T) {
	p := newTestPromoter(t)
	require.NoError(t, p.snapshots.Save("staging", map[string]string{"FOO": "new"}))
	require.NoError(t, p.snapshots.Save("prod", map[string]string{"FOO": "old"}))

	res, err := p.Promote(PromoteOptions{SourceEnv: "staging", TargetEnv: "prod", Overwrite: false})
	require.NoError(t, err)
	assert.Empty(t, res.Promoted)
	assert.Contains(t, res.Skipped, "FOO")

	dst, _ := p.snapshots.Load("prod")
	assert.Equal(t, "old", dst["FOO"])
}

func TestPromoter_OverwriteReplaces(t *testing.T) {
	p := newTestPromoter(t)
	require.NoError(t, p.snapshots.Save("staging", map[string]string{"FOO": "new"}))
	require.NoError(t, p.snapshots.Save("prod", map[string]string{"FOO": "old"}))

	res, err := p.Promote(PromoteOptions{SourceEnv: "staging", TargetEnv: "prod", Overwrite: true})
	require.NoError(t, err)
	assert.Contains(t, res.Promoted, "FOO")

	dst, _ := p.snapshots.Load("prod")
	assert.Equal(t, "new", dst["FOO"])
}

func TestPromoter_DryRunDoesNotSave(t *testing.T) {
	p := newTestPromoter(t)
	require.NoError(t, p.snapshots.Save("staging", map[string]string{"FOO": "bar"}))

	res, err := p.Promote(PromoteOptions{SourceEnv: "staging", TargetEnv: "prod", DryRun: true})
	require.NoError(t, err)
	assert.True(t, res.DryRun)
	assert.Contains(t, res.Promoted, "FOO")

	_, err = p.snapshots.Load("prod")
	assert.Error(t, err)
}

func TestPromoter_SelectedKeys(t *testing.T) {
	p := newTestPromoter(t)
	require.NoError(t, p.snapshots.Save("staging", map[string]string{"A": "1", "B": "2", "C": "3"}))

	res, err := p.Promote(PromoteOptions{SourceEnv: "staging", TargetEnv: "prod", Keys: []string{"A", "C"}})
	require.NoError(t, err)
	assert.Len(t, res.Promoted, 2)

	dst, _ := p.snapshots.Load("prod")
	assert.Equal(t, "1", dst["A"])
	_, ok := dst["B"]
	assert.False(t, ok)
}

func TestPromoteResult_Summary(t *testing.T) {
	r := PromoteResult{Promoted: []string{"A", "B"}, Skipped: []string{"C"}, DryRun: true}
	assert.Contains(t, r.Summary(), "[dry-run]")
	assert.Contains(t, r.Summary(), "2")
}
