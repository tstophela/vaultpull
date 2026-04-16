package vault_test

import (
	"testing"

	"github.com/yourorg/vaultpull/internal/vault"
)

func TestNormalizePath_KVv1(t *testing.T) {
	got := vault.NormalizePath("secret", "secret/myapp", vault.KVv1)
	if got != "secret/myapp" {
		t.Errorf("KVv1: expected %q, got %q", "secret/myapp", got)
	}
}

func TestNormalizePath_KVv2_AlreadyNormalized(t *testing.T) {
	got := vault.NormalizePath("secret", "secret/data/myapp", vault.KVv2)
	if got != "secret/data/myapp" {
		t.Errorf("already normalized: expected %q, got %q", "secret/data/myapp", got)
	}
}

func TestNormalizePath_KVv2_InjectsData(t *testing.T) {
	got := vault.NormalizePath("secret", "secret/myapp", vault.KVv2)
	if got != "secret/data/myapp" {
		t.Errorf("inject data: expected %q, got %q", "secret/data/myapp", got)
	}
}

func TestSplitMountPath(t *testing.T) {
	tests := []struct {
		input         string
		wantMount     string
		wantSubpath   string
	}{
		{"secret/myapp/prod", "secret", "myapp/prod"},
		{"secret/myapp", "secret", "myapp"},
		{"secret", "secret", ""},
		{"/secret/myapp", "secret", "myapp"},
	}

	for _, tt := range tests {
		m, s := vault.SplitMountPath(tt.input)
		if m != tt.wantMount || s != tt.wantSubpath {
			t.Errorf("SplitMountPath(%q) = (%q, %q), want (%q, %q)",
				tt.input, m, s, tt.wantMount, tt.wantSubpath)
		}
	}
}
