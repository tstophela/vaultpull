package env

import "sort"

// CompareResult holds the result of comparing two env maps.
type CompareResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // key -> [old, new]
	Same    map[string]string
}

// Compare compares two env maps (old vs new) and returns a CompareResult.
func Compare(old, next map[string]string) CompareResult {
	r := CompareResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
		Same:    make(map[string]string),
	}
	for k, v := range next {
		ov, ok := old[k]
		if !ok {
			r.Added[k] = v
		} else if ov != v {
			r.Changed[k] = [2]string{ov, v}
		} else {
			r.Same[k] = v
		}
	}
	for k, v := range old {
		if _, ok := next[k]; !ok {
			r.Removed[k] = v
		}
	}
	return r
}

// HasChanges returns true if there are any added, removed, or changed keys.
func (r CompareResult) HasChanges() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Changed) > 0
}

// SortedAdded returns added keys in sorted order.
func (r CompareResult) SortedAdded() []string {
	return sortedKeys(r.Added)
}

// SortedRemoved returns removed keys in sorted order.
func (r CompareResult) SortedRemoved() []string {
	return sortedKeys(r.Removed)
}

// SortedChanged returns changed keys in sorted order.
func (r CompareResult) SortedChanged() []string {
	keys := make([]string, 0, len(r.Changed))
	for k := range r.Changed {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
