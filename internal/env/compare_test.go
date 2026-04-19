package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare_Added(t *testing.T) {
	old := map[string]string{"A": "1"}
	new := map[string]string{"A": "1", "B": "2"}
	r := Compare(old, new)
	assert.Equal(t, map[string]string{"B": "2"}, r.Added)
	assert.Empty(t, r.Removed)
	assert.Empty(t, r.Changed)
	assert.True(t, r.HasChanges())
}

func TestCompare_Removed(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	new := map[string]string{"A": "1"}
	r := Compare(old, new)
	assert.Equal(t, map[string]string{"B": "2"}, r.Removed)
	assert.Empty(t, r.Added)
	assert.True(t, r.HasChanges())
}

func TestCompare_Changed(t *testing.T) {
	old := map[string]string{"A": "old"}
	new := map[string]string{"A": "new"}
	r := Compare(old, new)
	assert.Equal(t, [2]string{"old", "new"}, r.Changed["A"])
	assert.True(t, r.HasChanges())
}

func TestCompare_Same(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	new := map[string]string{"A": "1", "B": "2"}
	r := Compare(old, new)
	assert.False(t, r.HasChanges())
	assert.Equal(t, old, r.Same)
}

func TestCompare_SortedKeys(t *testing.T) {
	old := map[string]string{}
	new := map[string]string{"Z": "1", "A": "2", "M": "3"}
	r := Compare(old, new)
	assert.Equal(t, []string{"A", "M", "Z"}, r.SortedAdded())
}

func TestCompare_Mixed(t *testing.T) {
	old := map[string]string{"A": "1", "B": "old", "C": "3"}
	new := map[string]string{"A": "1", "B": "new", "D": "4"}
	r := Compare(old, new)
	assert.Equal(t, map[string]string{"D": "4"}, r.Added)
	assert.Equal(t, map[string]string{"C": "3"}, r.Removed)
	assert.Equal(t, [2]string{"old", "new"}, r.Changed["B"])
	assert.Equal(t, map[string]string{"A": "1"}, r.Same)
}
