package env

import (
	"bufio"
	"os"
	"strings"
)

// Reader reads an existing .env file into a key-value map.
type Reader struct {
	path string
}

// NewReader creates a new Reader for the given file path.
func NewReader(path string) *Reader {
	return &Reader{path: path}
}

// Read parses the .env file and returns a map of key-value pairs.
// Lines beginning with '#' and empty lines are ignored.
// Values wrapped in double quotes are unquoted.
func (r *Reader) Read() (map[string]string, error) {
	f, err := os.Open(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	defer f.Close()

	result := make(map[string]string)
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
		val := strings.TrimSpace(parts[1])

		// Strip surrounding double quotes if present.
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
		}

		result[key] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
