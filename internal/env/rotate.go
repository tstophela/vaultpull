package env

import (
	"fmt"
	"time"
)

// RotationRecord tracks a single secret rotation event.
type RotationRecord struct {
	Key       string
	OldValue  string
	NewValue  string
	RotatedAt time.Time
}

// Rotator detects and records secret rotations between two env maps.
type Rotator struct {
	redactor *Redactor
	now      func() time.Time
}

// NewRotator creates a Rotator with optional custom time source.
func NewRotator(r *Redactor) *Rotator {
	if r == nil {
		r = NewRedactor()
	}
	return &Rotator{redactor: r, now: time.Now}
}

// Detect compares existing and incoming env maps and returns records for
// keys whose values have changed.
func (r *Rotator) Detect(existing, incoming map[string]string) []RotationRecord {
	var records []RotationRecord
	for key, newVal := range incoming {
		oldVal, exists := existing[key]
		if !exists || oldVal == newVal {
			continue
		}
		records = append(records, RotationRecord{
			Key:       key,
			OldValue:  r.redactor.Redact(key, oldVal),
			NewValue:  r.redactor.Redact(key, newVal),
			RotatedAt: r.now(),
		})
	}
	return records
}

// Summary returns a human-readable summary of rotation records.
func Summary(records []RotationRecord) string {
	if len(records) == 0 {
		return "no secrets rotated"
	}
	s := fmt.Sprintf("%d secret(s) rotated:\n", len(records))
	for _, rec := range records {
		s += fmt.Sprintf("  %s: %s -> %s (at %s)\n",
			rec.Key, rec.OldValue, rec.NewValue,
			rec.RotatedAt.Format(time.RFC3339))
	}
	return s
}
