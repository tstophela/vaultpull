package env

import "fmt"

// Strategy controls how incoming secrets are merged with existing env values.
type Strategy int

const (
	// StrategyOverwrite replaces all existing values with incoming secrets.
	StrategyOverwrite Strategy = iota
	// StrategyPreserve keeps existing values; only adds new keys.
	StrategyPreserve
	// StrategyInteractive is reserved for future prompt-based merging.
	StrategyInteractive
)

// ParseStrategy converts a string flag value into a Strategy.
func ParseStrategy(s string) (Strategy, error) {
	switch s {
	case "overwrite", "":
		return StrategyOverwrite, nil
	case "preserve":
		return StrategyPreserve, nil
	case "interactive":
		return StrategyInteractive, nil
	default:
		return StrategyOverwrite, fmt.Errorf("unknown merge strategy %q: choose overwrite, preserve, or interactive", s)
	}
}

// Apply merges incoming secrets into existing values according to the strategy.
// It returns the resulting map that should be written to the .env file.
func (st Strategy) Apply(existing, incoming map[string]string) map[string]string {
	result := make(map[string]string, len(incoming))

	switch st {
	case StrategyPreserve:
		// Start with incoming, then overwrite with anything already set locally.
		for k, v := range incoming {
			result[k] = v
		}
		for k, v := range existing {
			result[k] = v
		}
	default: // StrategyOverwrite
		for k, v := range existing {
			result[k] = v
		}
		for k, v := range incoming {
			result[k] = v
		}
	}

	return result
}
