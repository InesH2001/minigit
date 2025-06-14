package utils

import (
	"os"
	"strings"
)

func GetCurrentBranch() string {
	data, err := os.ReadFile(".miniGit/HEAD")
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(data))
}
