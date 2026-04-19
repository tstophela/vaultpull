package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ImportResult holds the outcome of an import operation.
type ImportResult struct {
	Imported int
	Skipped  int
	Keys     []string
}

// Importer reads an existing .env file and merges secrets into a target map.
type Importer struct {
	strategy Strategy
	filter   *Filter
}

// NewImporter creates an Importer with the given merge strategy and filter.
func NewImporter(strategy Strategy, filter *Filter) *Importer {
	return &Importer{strategy: strategy, filter: filter}
}

// ImportFile reads key=value pairs from srcPath and merges them into existing.
func (im *Importer) ImportFile(srcPath string, existing map[string]string) (map[string]string, ImportResult, error) {
	f, err := os.Open(srcPath)
	if err != nil {
		return nil, ImportResult{}, fmt.Errorf("import: open %s: %w", srcPath, err)
	}
	defer f.Close()

	incoming := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		if im.filter == nil || im.filter.Allow(key) {
			incoming[key] = val
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, ImportResult{}, fmt.Errorf("import: scan %s: %w", srcPath, err)
	}

	result := make(map[string]string)
	for k, v := range existing {
		result[k] = v
	}

	var res ImportResult
	for k, v := range incoming {
		if im.strategy.Apply(k, existing[k], v) {
			result[k] = v
			res.Imported++
			res.Keys = append(res.Keys, k)
		} else {
			res.Skipped++
		}
	}
	return result, res, nil
}
