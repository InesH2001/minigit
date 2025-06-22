package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func LoadGitignore() ([]string, error) {
	var patterns []string

	if !FileExists(".gitignore") {
		return patterns, nil
	}

	file, err := os.Open(".gitignore")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		patterns = append(patterns, line)
	}

	return patterns, scanner.Err()
}

func IsIgnored(filePath string) bool {
	patterns, err := LoadGitignore()
	if err != nil {
		return false
	}

	for _, pattern := range patterns {
		if matchPattern(filePath, pattern) {
			return true
		}
	}

	return false
}

func matchPattern(filePath, pattern string) bool {
	filePath = filepath.Clean(filePath)
	pattern = filepath.Clean(pattern)

	if filePath == pattern {
		return true
	}

	if strings.Contains(pattern, "*") {
		matched, _ := filepath.Match(pattern, filePath)
		if matched {
			return true
		}

		matched, _ = filepath.Match(pattern, filepath.Base(filePath))
		return matched
	}

	if strings.HasSuffix(pattern, "/") {
		dirPattern := strings.TrimSuffix(pattern, "/")
		return strings.HasPrefix(filePath, dirPattern+"/") || filePath == dirPattern
	}

	if strings.HasPrefix(pattern, "*.") {
		ext := strings.TrimPrefix(pattern, "*")
		return strings.HasSuffix(filePath, ext)
	}

	return strings.Contains(filePath, pattern)
}

func FilterIgnoredFiles(files []string) []string {
	var filtered []string

	for _, file := range files {
		if isInternalPath(file) {
			continue
		}

		if !IsIgnored(file) {
			filtered = append(filtered, file)
		}
	}

	return filtered
}

func isInternalPath(path string) bool {
	return strings.HasPrefix(path, ".git") ||
		strings.HasPrefix(path, ".miniGit") ||
		strings.Contains(path, "/.git/") ||
		strings.Contains(path, "/.miniGit/")
}
