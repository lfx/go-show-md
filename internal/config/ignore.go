package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

var DefaultIgnoredPatterns = []string{
	".git", ".venv", "venv", "node_modules", "vendor", ".env",
}

func GetGlobalIgnorePatterns() []string {
	var patterns []string
	home, err := os.UserHomeDir()
	if err == nil {
		globalPath := filepath.Join(home, ".go-show-md-ignore")
		patterns = append(patterns, readLines(globalPath)...)
	}
	return patterns
}

func GetLocalIgnorePatterns(dir string) []string {
	localPath := filepath.Join(dir, ".go-show-md-ignore")
	return readLines(localPath)
}

func readLines(path string) []string {
	var lines []string
	f, err := os.Open(path)
	if err != nil {
		return lines
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}
	return lines
}

func MergeIgnorePatterns(defaults, global, local, configPatterns []string) []string {
	patternMap := make(map[string]bool)
	var merged []string

	addPatterns := func(patterns []string) {
		for _, p := range patterns {
			if !patternMap[p] {
				patternMap[p] = true
				merged = append(merged, p)
			}
		}
	}

	addPatterns(defaults)
	addPatterns(global)
	addPatterns(local)
	addPatterns(configPatterns)

	return merged
}

func IsIgnored(path string, info os.FileInfo, patterns []string) bool {
	name := info.Name()

	// First match against base name
	for _, pattern := range patterns {
		if matched, err := filepath.Match(pattern, name); err == nil && matched {
			return true
		}
	}

	return false
}
