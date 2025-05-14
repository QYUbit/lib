package envfile

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// Loads an .env file and sets enviroment variables accordingly.
func Load(path string) error {
	path = ensurePath(path)

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	lines, err := getLines(file)
	if err != nil {
		return err
	}

	varMap := parseLines(lines)
	insertEnvVars(varMap)

	return nil
}

func ensurePath(path string) string {
	if path == "" {
		return ".env"
	}
	return path
}

func getLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func parseLines(lines []string) map[string]string {
	varMap := make(map[string]string)

	for _, line := range lines {
		line = strings.TrimSpace(ignoreComments(line))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		before, after, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		varMap[strings.TrimSpace(before)] = strings.TrimSpace(after)
	}

	return varMap
}

func ignoreComments(line string) string {
	if idx := strings.Index(line, "#"); idx != -1 {
		return line[:idx]
	}
	return line
}

func insertEnvVars(varMap map[string]string) {
	for k, v := range varMap {
		if err := os.Setenv(k, v); err != nil {
			// Optionale Fehlerbehandlung oder Logging
		}
	}
}
