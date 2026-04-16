package env

import (
	"testing"
)

func TestDiff_Added(t *testing.T) {
	existing := map[string]string{}
	incoming := map[string]string{"NEW_KEY": "val"}
	res := Diff(existing, incoming)
	if res.Added["NEW_KEY"] != "val" {
		t.Errorf("expected NEW_KEY in Added")
	}
	if len(res.Updated) != 0 || len(res.Unchanged) != 0 {
		t.Errorf("unexpected entries: %v %v", res.Updated, res.Unchanged)
	}
}

func TestDiff_Updated(t *testing.T) {
	existing := map[string]string{"KEY": "old"}
	incoming := map[string]string{"KEY": "new"}
	res := Diff(existing, incoming)
	if res.Updated["KEY"] != "new" {
		t.Errorf("expected KEY in Updated")
	}
	if len(res.Added) != 0 || len(res.Unchanged) != 0 {
		t.Errorf("unexpected entries")
	}
}

func TestDiff_Unchanged(t *testing.T) {
	existing := map[string]string{"KEY": "same"}
	incoming := map[string]string{"KEY": "same"}
	res := Diff(existing, incoming)
	if res.Unchanged["KEY"] != "same" {
		t.Errorf("expected KEY in Unchanged")
	}
	if len(res.Added) != 0 || len(res.Updated) != 0 {
		t.Errorf("unexpected entries")
	}
}

func TestMerge_IncomingWins(t *testing.T) {
	existing := map[string]string{"A": "1", "B": "2"}
	incoming := map[string]string{"B": "99", "C": "3"}
	out := Merge(existing, incoming)
	if out["A"] != "1" || out["B"] != "99" || out["C"] != "3" {
		t.Errorf("unexpected merge result: %v", out)
	}
}

func TestMerge_EmptyExisting(t *testing.T) {
	out := Merge(map[string]string{}, map[string]string{"X": "y"})
	if out["X"] != "y" {
		t.Errorf("expected X=y")
	}
}
