package utils

import (
	"bufio"
	"os"
	"strings"
)

func ReadLastCommit() map[string]string {
	result := make(map[string]string)

	headData, err := os.ReadFile(".minigit/HEAD")
	if err != nil {
		return result
	}

	branch := strings.TrimSpace(string(headData))
	refPath := ".minigit/refs/" + branch

	commitHashBytes, err := os.ReadFile(refPath)
	if err != nil {
		return result
	}

	commitHash := strings.TrimSpace(string(commitHashBytes))
	commitPath := ".minigit/objects/commits/" + commitHash

	f, err := os.Open(commitPath)
	if err != nil {
		return result
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}

	return result
}
