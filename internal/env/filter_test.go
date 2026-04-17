package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter_Allow_NoPrefixes(t *testing.T) {
	f := NewFilter(nil, nil)
	assert.True(t, f.Allow("ANY_KEY"))
	assert.True(t, f.Allow("another"))
}

func TestFilter_Allow_WithPrefix(t *testing.T) {
	f := NewFilter([]string{"APP_"}, nil)
	assert.True(t, f.Allow("APP_SECRET"))
	assert.True(t, f.Allow("app_token")) // case-insensitive
	assert.False(t, f.Allow("DB_PASSWORD"))
}

func TestFilter_Allow_Excluded(t *testing.T) {
	f := NewFilter(nil, []string{"VAULT_TOKEN"})
	assert.False(t, f.Allow("VAULT_TOKEN"))
	assert.False(t, f.Allow("vault_token")) // case-insensitive
	assert.True(t, f.Allow("OTHER_KEY"))
}

func TestFilter_Allow_ExcludeTakesPrecedence(t *testing.T) {
	f := NewFilter([]string{"APP_"}, []string{"APP_SKIP"})
	assert.True(t, f.Allow("APP_KEEP"))
	assert.False(t, f.Allow("APP_SKIP"))
}

func TestFilter_Apply(t *testing.T) {
	f := NewFilter([]string{"DB_"}, []string{"DB_LEGACY"})
	secrets := map[string]string{
		"DB_HOST":   "localhost",
		"DB_PASS":   "secret",
		"DB_LEGACY": "old",
		"APP_KEY":   "value",
	}

	result := f.Apply(secrets)

	assert.Equal(t, map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "secret",
	}, result)
}

func TestFilter_Apply_EmptySecrets(t *testing.T) {
	f := NewFilter(nil, nil)
	result := f.Apply(map[string]string{})
	assert.Empty(t, result)
}
