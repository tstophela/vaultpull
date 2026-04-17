package env

import (
	"testing"
)

func TestFormatter_FormatLine_Plain(t *testing.T) {
	f := NewFormatter(false)
	got := f.FormatLine("FOO", "bar")
	if got != "FOO=bar" {
		t.Errorf("expected FOO=bar, got %s", got)
	}
}

func TestFormatter_FormatLine_QuoteAll(t *testing.T) {
	f := NewFormatter(true)
	got := f.FormatLine("FOO", "bar")
	if got != `FOO="bar"` {
		t.Errorf("expected FOO=\"bar\", got %s", got)
	}
}

func TestFormatter_FormatLine_AutoQuoteSpace(t *testing.T) {
	f := NewFormatter(false)
	got := f.FormatLine("MSG", "hello world")
	if got != `MSG="hello world"` {
		t.Errorf("unexpected: %s", got)
	}
}

func TestFormatter_FormatLine_AutoQuoteHash(t *testing.T) {
	f := NewFormatter(false)
	got := f.FormatLine("VAL", "foo#bar")
	if got != `VAL="foo#bar"` {
		t.Errorf("unexpected: %s", got)
	}
}

func TestFormatter_FormatLine_EscapesInnerQuotes(t *testing.T) {
	f := NewFormatter(false)
	got := f.FormatLine("K", `say "hi"`)
	expected := `K="say \"hi\""`
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestFormatter_FormatMap(t *testing.T) {
	f := NewFormatter(false)
	secrets := map[string]string{
		"ALPHA": "one",
		"BETA":  "two",
		"GAMMA": "three four",
	}
	keys := []string{"ALPHA", "BETA", "GAMMA"}
	lines := f.FormatMap(secrets, keys)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "ALPHA=one" {
		t.Errorf("unexpected: %s", lines[0])
	}
	if lines[2] != `GAMMA="three four"` {
		t.Errorf("unexpected: %s", lines[2])
	}
}

func TestFormatter_FormatMap_SkipsMissingKeys(t *testing.T) {
	f := NewFormatter(false)
	secrets := map[string]string{"A": "1"}
	lines := f.FormatMap(secrets, []string{"A", "B"})
	if len(lines) != 1 {
		t.Errorf("expected 1 line, got %d", len(lines))
	}
}
