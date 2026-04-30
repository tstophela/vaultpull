package env

import (
	"testing"
)

func TestResolver_Resolve_FirstSourceWins(t *testing.T) {
	r := NewResolver([]ChainEntry{
		{Name: "vault", Values: map[string]string{"DB_PASS": "secret", "API_KEY": "abc"}},
		{Name: "local", Values: map[string]string{"DB_PASS": "override", "LOG_LEVEL": "info"}},
	})

	results := r.Resolve()
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	byKey := make(map[string]ChainResult)
	for _, cr := range results {
		byKey[cr.Key] = cr
	}

	if byKey["DB_PASS"].Value != "secret" || byKey["DB_PASS"].Source != "vault" {
		t.Errorf("DB_PASS: expected vault/secret, got %s/%s", byKey["DB_PASS"].Source, byKey["DB_PASS"].Value)
	}
	if byKey["LOG_LEVEL"].Source != "local" {
		t.Errorf("LOG_LEVEL should come from local, got %s", byKey["LOG_LEVEL"].Source)
	}
}

func TestResolver_Resolve_SortedOutput(t *testing.T) {
	r := NewResolver([]ChainEntry{
		{Name: "src", Values: map[string]string{"ZEBRA": "z", "ALPHA": "a", "MANGO": "m"}},
	})
	results := r.Resolve()
	keys := []string{results[0].Key, results[1].Key, results[2].Key}
	expected := []string{"ALPHA", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("position %d: expected %s, got %s", i, expected[i], k)
		}
	}
}

func TestResolver_ResolveKey_Found(t *testing.T) {
	r := NewResolver([]ChainEntry{
		{Name: "env", Values: map[string]string{"TOKEN": "tok123"}},
	})
	cr, err := r.ResolveKey("TOKEN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cr.Value != "tok123" || cr.Source != "env" {
		t.Errorf("unexpected result: %+v", cr)
	}
}

func TestResolver_ResolveKey_NotFound(t *testing.T) {
	r := NewResolver([]ChainEntry{
		{Name: "env", Values: map[string]string{"FOO": "bar"}},
	})
	_, err := r.ResolveKey("MISSING")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestResolver_Flatten(t *testing.T) {
	r := NewResolver([]ChainEntry{
		{Name: "a", Values: map[string]string{"X": "1"}},
		{Name: "b", Values: map[string]string{"X": "2", "Y": "3"}},
	})
	flat := r.Flatten()
	if flat["X"] != "1" {
		t.Errorf("X should be 1 from source a, got %s", flat["X"])
	}
	if flat["Y"] != "3" {
		t.Errorf("Y should be 3 from source b, got %s", flat["Y"])
	}
}

func TestResolver_EmptyChain(t *testing.T) {
	r := NewResolver(nil)
	results := r.Resolve()
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
