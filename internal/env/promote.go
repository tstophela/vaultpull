package env

import (
	"fmt"
	"sort"
)

// PromoteOptions configures a promotion run.
type PromoteOptions struct {
	SourceEnv  string
	TargetEnv  string
	Keys       []string // empty = all keys
	DryRun     bool
	Overwrite  bool
}

// PromoteResult holds the outcome of a promotion.
type PromoteResult struct {
	Promoted []string
	Skipped  []string
	DryRun   bool
}

func (r PromoteResult) Summary() string {
	prefix := ""
	if r.DryRun {
		prefix = "[dry-run] "
	}
	return fmt.Sprintf("%spromoted %d key(s), skipped %d key(s)", prefix, len(r.Promoted), len(r.Skipped))
}

// Promoter copies secrets from one env snapshot to another.
type Promoter struct {
	snapshots *SnapshotManager
}

func NewPromoter(dir string) *Promoter {
	return &Promoter{snapshots: NewSnapshotManager(dir)}
}

func (p *Promoter) Promote(opts PromoteOptions) (PromoteResult, error) {
	src, err := p.snapshots.Load(opts.SourceEnv)
	if err != nil {
		return PromoteResult{}, fmt.Errorf("load source %q: %w", opts.SourceEnv, err)
	}

	dst, err := p.snapshots.Load(opts.TargetEnv)
	if err != nil {
		dst = map[string]string{}
	}

	keys := opts.Keys
	if len(keys) == 0 {
		for k := range src {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}

	result := PromoteResult{DryRun: opts.DryRun}
	for _, k := range keys {
		v, ok := src[k]
		if !ok {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		if _, exists := dst[k]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		dst[k] = v
		result.Promoted = append(result.Promoted, k)
	}

	if !opts.DryRun && len(result.Promoted) > 0 {
		if err := p.snapshots.Save(opts.TargetEnv, dst); err != nil {
			return result, fmt.Errorf("save target %q: %w", opts.TargetEnv, err)
		}
	}
	return result, nil
}
