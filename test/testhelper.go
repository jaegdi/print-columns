package main

import (
	"os"
	"path/filepath"
	"runtime"
)

// getTestDataPath returns the absolute path to a test data file.
// It uses runtime.Caller to find the test directory location.
func getTestDataPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(currentFile)
	return filepath.Join(testDir, "data", filename)
}

// readTestDataFile reads a test data file and returns its contents as a slice of strings.
func readTestDataFile(filename string) ([]string, error) {
	path := getTestDataPath(filename)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// Split by newlines and remove empty trailing line
	lines := splitLines(string(content))
	return lines, nil
}

// splitLines splits content by newlines, handling both \n and \r\n
func splitLines(content string) []string {
	var lines []string
	var current []byte
	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			lines = append(lines, string(current))
			current = nil
		} else if content[i] == '\r' {
			// Skip \r, will handle \n next
			continue
		} else {
			current = append(current, content[i])
		}
	}
	// Add last line if not empty
	if len(current) > 0 {
		lines = append(lines, string(current))
	}
	return lines
}
