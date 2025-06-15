package utils

import (
    "os"
    "strings"
    "fmt"
)

func GetCurrentBranch() (string, error) {
	data, err := os.ReadFile(".miniGit/HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}

func GetCurrentBranchRefPath() (string, error) {
	data, err := GetCurrentBranch()
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(data, "ref: "), nil
}

func GetCurrentBranchAndParentCommitHash() (string, string, error) {
	headRefPath, err := GetCurrentBranchRefPath()
	if err != nil {
		return "", "", err
	}

	parentHash, err := GetCommitHashFromRef(headRefPath)
	if err != nil {
		return "", "", err
	}

	return headRefPath, parentHash, nil
}

func GetHeadCommit() (string, error) {
	refPath, err := GetCurrentBranchRefPath()
	if err != nil {
		return "", err
	}

	return GetCommitHashFromRef(refPath)
}
