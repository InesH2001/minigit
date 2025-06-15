package utils

import (
	"bufio"
	"os"
	"strings"
	"fmt"
)

func GetCommitHashFromRef(ref string) (string, error) {
	data, err := os.ReadFile(".miniGit/" + ref)
	if err != nil {
		return "", fmt.Errorf("failed to read ref %s: %w", ref, err)
	}
	return strings.TrimSpace(string(data)), nil
}

func GetParentsFromCommit(commitHash string) []string {
    data, _ := os.ReadFile("./miniGit/objects/commits/" + commitHash)
    lines := strings.Split(string(data), "\n")
    var parents []string

    for _, line := range lines {
        if strings.HasPrefix(line, "parent: ") {
            parentHash := strings.TrimSpace(strings.TrimPrefix(line, "parent: "))
            parents = append(parents, parentHash)
        }
    }

    return parents
}

func ReadLastCommit() map[string]string {
	result := make(map[string]string)

	headData, err := os.ReadFile(".miniGit/HEAD")
	if err != nil {
		return result
	}

	headContent := strings.TrimSpace(string(headData))

	var commitHash string
	if strings.HasPrefix(headContent, "ref: ") {
		ref := strings.TrimPrefix(headContent, "ref: ")
		refPath := ".miniGit/" + ref
		hashBytes, err := os.ReadFile(refPath)
		if err != nil {
			return result
		}
		commitHash = strings.TrimSpace(string(hashBytes))
	} else {
		commitHash = headContent
	}

	commitPath := ".miniGit/objects/commits/" + commitHash

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
