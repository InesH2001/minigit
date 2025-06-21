// utils/branch.go
package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func GetCurrentBranchName() string {
	data, err := os.ReadFile(".miniGit/HEAD")
	if err != nil {
		return ""
	}
	ref := strings.TrimSpace(string(data))
	return filepath.Base(ref)
}
