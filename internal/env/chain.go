package env

import (
	"fmt"
	"sort"
)

// ChainEntry represents a single source in the resolution chain.
type ChainEntry struct {
	Name   string
	Values map[string]string
}

// ChainResult holds the resolved value and its origin.
type ChainResult struct {
	Key    string
	Value  string
	Source string
}

// Resolver resolves environment variables from an ordered chain of sources.
// Earlier sources take precedence over later ones.
type Resolver struct {
	chain []ChainEntry
}

// NewResolver creates a Resolver with the given ordered chain entries.
func NewResolver(chain []ChainEntry) *Resolver {
	return &Resolver{chain: chain}
}

// Resolve returns ChainResult for each key found across all sources.
// The first source that defines a key wins.
func (r *Resolver) Resolve() []ChainResult {
	seen := make(map[string]ChainResult)

	for _, entry := range r.chain {
		for k, v := range entry.Values {
			if _, exists := seen[k]; !exists {
				seen[k] = ChainResult{Key: k, Value: v, Source: entry.Name}
			}
		}
	}

	results := make([]ChainResult, 0, len(seen))
	for _, cr := range seen {
		results = append(results, cr)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})
	return results
}

// ResolveKey returns the ChainResult for a single key, or an error if not found.
func (r *Resolver) ResolveKey(key string) (ChainResult, error) {
	for _, entry := range r.chain {
		if v, ok := entry.Values[key]; ok {
			return ChainResult{Key: key, Value: v, Source: entry.Name}, nil
		}
	}
	return ChainResult{}, fmt.Errorf("key %q not found in any source", key)
}

// Flatten returns a merged map of all resolved key/value pairs.
func (r *Resolver) Flatten() map[string]string {
	results := r.Resolve()
	out := make(map[string]string, len(results))
	for _, cr := range results {
		out[cr.Key] = cr.Value
	}
	return out
}
