package env

// DiffResult holds the categorised changes between existing and incoming secrets.
type DiffResult struct {
	Added   map[string]string
	Updated map[string]string
	Unchanged map[string]string
}

// Diff compares existing key/value pairs against incoming ones and returns
// a DiffResult describing what is new, changed, or unchanged.
// Keys present only in existing are treated as unchanged (not overwritten).
func Diff(existing, incoming map[string]string) DiffResult {
	res := DiffResult{
		Added:     make(map[string]string),
		Updated:   make(map[string]string),
		Unchanged: make(map[string]string),
	}
	for k, v := range incoming {
		old, exists := existing[k]
		if !exists {
			res.Added[k] = v
		} else if old != v {
			res.Updated[k] = v
		} else {
			res.Unchanged[k] = v
		}
	}
	return res
}

// Merge combines existing and incoming maps, with incoming values taking
// precedence for overlapping keys.
func Merge(existing, incoming map[string]string) map[string]string {
	out := make(map[string]string, len(existing)+len(incoming))
	for k, v := range existing {
		out[k] = v
	}
	for k, v := range incoming {
		out[k] = v
	}
	return out
}
