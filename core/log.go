package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Log() error {
	headData, err := os.ReadFile(".miniGit/HEAD")
	if err != nil {
		return fmt.Errorf("failed to read HEAD: %v", err)
	}

	headRef := strings.TrimSpace(string(headData))

	var commitHash string
	if strings.HasPrefix(headRef, "ref: ") {
		refPath := ".miniGit/" + strings.TrimPrefix(headRef, "ref: ")
		hashBytes, err := os.ReadFile(refPath)
		if err != nil {
			return fmt.Errorf("failed to read ref file: %v", err)
		}
		commitHash = strings.TrimSpace(string(hashBytes))
	} else {
		commitHash = headRef
	}

	for commitHash != "" {
		commitPath := ".miniGit/objects/commits/" + commitHash
		file, err := os.Open(commitPath)
		if err != nil {
			break
		}

		fmt.Printf("commit %s\n", commitHash)

		scanner := bufio.NewScanner(file)

		var nextHash string
		var message string
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "parent: ") {
				nextHash = strings.TrimPrefix(line, "parent: ")
			} else if strings.HasPrefix(line, "message: ") {
				message = strings.TrimPrefix(line, "message: ")
			}
		}

		if message != "" {
			fmt.Printf("    %s\n\n", message)
		} else {
			fmt.Println()
		}

		commitHash = nextHash
		file.Close()
	}

	return nil
}
