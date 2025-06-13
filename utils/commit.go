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

	headContent := strings.TrimSpace(string(headData))

	var commitHash string
	if strings.HasPrefix(headContent, "ref: ") {
		ref := strings.TrimPrefix(headContent, "ref: ")
		refPath := ".minigit/" + ref
		hashBytes, err := os.ReadFile(refPath)
		if err != nil {
			return result
		}
		commitHash = strings.TrimSpace(string(hashBytes))
	} else {
		commitHash = headContent
	}

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
