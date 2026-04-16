package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Writer handles writing secrets to .env files.
type Writer struct {
	path   string
	backup bool
}

// NewWriter creates a new Writer for the given file path.
// If backup is true, an existing file will be backed up before overwriting.
func NewWriter(path string, backup bool) *Writer {
	return &Writer{path: path, backup: backup}
}

// Write serializes the provided secrets map into a .env file.
// Keys are sorted for deterministic output.
func (w *Writer) Write(secrets map[string]string) error {
	if w.backup {
		if err := w.backupExisting(); err != nil {
			return fmt.Errorf("backup failed: %w", err)
		}
	}

	var sb strings.Builder
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := secrets[k]
		if strings.ContainsAny(v, " \t\n#") {
			v = fmt.Sprintf("%q", v)
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}

	return os.WriteFile(w.path, []byte(sb.String()), 0600)
}

func (w *Writer) backupExisting() error {
	if _, err := os.Stat(w.path); os.IsNotExist(err) {
		return nil
	}
	data, err := os.ReadFile(w.path)
	if err != nil {
		return err
	}
	return os.WriteFile(w.path+".bak", data, 0600)
}
