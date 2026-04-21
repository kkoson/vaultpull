package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Writer handles writing secrets to .env files.
type Writer struct {
	filePath string
}

// NewWriter creates a new Writer targeting the given file path.
func NewWriter(filePath string) *Writer {
	return &Writer{filePath: filePath}
}

// Write serialises the provided key/value map into a .env file,
// creating or truncating the file as needed.
func (w *Writer) Write(secrets map[string]string) error {
	f, err := os.Create(w.filePath)
	if err != nil {
		return fmt.Errorf("env: create file %q: %w", w.filePath, err)
	}
	defer f.Close()

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		line := fmt.Sprintf("%s=%s\n", k, quoteValue(secrets[k]))
		if _, err := f.WriteString(line); err != nil {
			return fmt.Errorf("env: write key %q: %w", k, err)
		}
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("env: sync file %q: %w", w.filePath, err)
	}
	return nil
}

// quoteValue wraps a value in double-quotes when it contains
// whitespace, quotes, or the '#' comment character.
func quoteValue(v string) string {
	if strings.ContainsAny(v, " \t\n\r\"'#") {
		v = strings.ReplaceAll(v, `"`, `\"`)
		return `"` + v + `"`
	}
	return v
}
