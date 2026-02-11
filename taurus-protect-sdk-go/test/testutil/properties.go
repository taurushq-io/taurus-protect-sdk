package testutil

import (
	"bufio"
	"os"
	"strings"
)

// ParseProperties parses a .properties file (key=value format).
// Supports # comments, blank lines, and literal \n escape in values (for PEM keys).
func ParseProperties(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	props := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, "=")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		// Replace literal \n with newline (for PEM keys stored in properties)
		value = strings.ReplaceAll(value, `\n`, "\n")
		props[key] = value
	}
	return props, scanner.Err()
}
