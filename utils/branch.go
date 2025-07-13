package utils

import (
	"os"
	"strings"
)

func BranchExists(branchName string) bool {
	branchPath := ".miniGit/refs/heads/" + branchName
	return FileExists(branchPath)
}

func HasUncommittedChanges() bool {
	index, err := ReadIndex()
	if err != nil {
		return false
	}

	if len(index) == 0 {
		return false
	}

	lastCommit := ReadLastCommit()
	for filePath, hash := range index {
		if lastHash, ok := lastCommit[filePath]; !ok || lastHash != hash {
			return true
		}
	}

	return false
}

func GetCurrentBranchName() string {
	data, err := os.ReadFile(".miniGit/HEAD")
	if err != nil {
		return ""
	}

	ref := strings.TrimSpace(string(data))
	if strings.HasPrefix(ref, "ref: refs/heads/") {
		return strings.TrimPrefix(ref, "ref: refs/heads/")
	}

	return ""
}
