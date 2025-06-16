package utils

import (
	"fmt"
	"os"
	"strings"
)

func ReadIndex() (map[string]string, error) {
	content, err := os.ReadFile(".miniGit/index")
	if err != nil {
		return nil, fmt.Errorf("failed to read index: %w", err)
	}
	lines := strings.Split(string(content), "\n")
	index := make(map[string]string)
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		index[parts[0]] = parts[1]
	}
	return index, nil
}

func WriteIndex(index map[string]string) error {
	f, err := os.Create(".miniGit/index")
	if err != nil {
		return err
	}
	defer f.Close()

	for path, hash := range index {
		if _, err := fmt.Fprintf(f, "%s %s\n", path, hash); err != nil {
			return err
		}
	}
	return nil
}
