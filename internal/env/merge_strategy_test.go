package env

import (
	"testing"
)

func TestParseStrategy_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  Strategy
	}{
		{"overwrite", StrategyOverwrite},
		{"", StrategyOverwrite},
		{"preserve", StrategyPreserve},
		{"interactive", StrategyInteractive},
	}
	for _, tc := range cases {
		got, err := ParseStrategy(tc.input)
		if err != nil {
			t.Errorf("ParseStrategy(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseStrategy(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestParseStrategy_Invalid(t *testing.T) {
	_, err := ParseStrategy("magic")
	if err == nil {
		t.Fatal("expected error for unknown strategy, got nil")
	}
}

func TestStrategy_Apply_Overwrite(t *testing.T) {
	existing := map[string]string{"A": "old", "B": "keep"}
	incoming := map[string]string{"A": "new", "C": "added"}

	result := StrategyOverwrite.Apply(existing, incoming)

	if result["A"] != "new" {
		t.Errorf("expected A=new, got %s", result["A"])
	}
	if result["C"] != "added" {
		t.Errorf("expected C=added, got %s", result["C"])
	}
	if result["B"] != "keep" {
		t.Errorf("expected B=keep, got %s", result["B"])
	}
}

func TestStrategy_Apply_Preserve(t *testing.T) {
	existing := map[string]string{"A": "local", "B": "local"}
	incoming := map[string]string{"A": "vault", "C": "vault"}

	result := StrategyPreserve.Apply(existing, incoming)

	if result["A"] != "local" {
		t.Errorf("preserve: expected A=local, got %s", result["A"])
	}
	if result["C"] != "vault" {
		t.Errorf("preserve: expected C=vault, got %s", result["C"])
	}
}

func TestStrategy_Apply_EmptyExisting(t *testing.T) {
	incoming := map[string]string{"X": "1", "Y": "2"}
	result := StrategyOverwrite.Apply(nil, incoming)
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestStrategy_Apply_EmptyIncoming(t *testing.T) {
	existing := map[string]string{"A": "1", "B": "2"}
	result := StrategyOverwrite.Apply(existing, nil)
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if result["A"] != "1" || result["B"] != "2" {
		t.Errorf("expected existing keys to be preserved, got %v", result)
	}
}
