package cmd

import (
	"fmt"
	"os"
	"strings"
)

func Checkout(name string) {
	branchPath := ".miniGit/refs/heads/" + name

	commitHashBytes, err := os.ReadFile(branchPath)
	if err != nil {
		panic(err)
	}
	commitHash := strings.TrimSpace(string(commitHashBytes))

	headContent := "ref: refs/heads/" + name
	err = os.WriteFile(".miniGit/HEAD", []byte(headContent), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Switched to branch '%s' (commit %s)\n", name, commitHash)
}
