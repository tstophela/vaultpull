package env

import (
	"fmt"
	"io"
	"os"
	"time"
)

// AuditEntry records a single sync operation result.
type AuditEntry struct {
	Timestamp time.Time
	Path      string
	Added     []string
	Updated   []string
	Unchanged []string
}

// AuditLogger writes audit entries to a writer.
type AuditLogger struct {
	w io.Writer
}

// NewAuditLogger creates an AuditLogger writing to w.
// Pass nil to use os.Stdout.
func NewAuditLogger(w io.Writer) *AuditLogger {
	if w == nil {
		w = os.Stdout
	}
	return &AuditLogger{w: w}
}

// Log writes a human-readable summary of the audit entry.
func (a *AuditLogger) Log(e AuditEntry) {
	fmt.Fprintf(a.w, "[%s] sync: %s\n", e.Timestamp.Format(time.RFC3339), e.Path)
	for _, k := range e.Added {
		fmt.Fprintf(a.w, "  + %s\n", k)
	}
	for _, k := range e.Updated {
		fmt.Fprintf(a.w, "  ~ %s\n", k)
	}
	for _, k := range e.Unchanged {
		fmt.Fprintf(a.w, "  = %s\n", k)
	}
	fmt.Fprintf(a.w, "  summary: %d added, %d updated, %d unchanged\n",
		len(e.Added), len(e.Updated), len(e.Unchanged))
}
