package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Branch(name string) {
	if name == "" {
		branches, err := filepath.Glob(".minigit/refs/heads/*")
		if err != nil {
			panic(err)
		}

		headData, err := os.ReadFile(".minigit/HEAD")
		if err != nil {
			panic(err)
		}
		headRef := strings.TrimSpace(string(headData))
		currentBranch := strings.TrimPrefix(headRef, "ref: refs/heads/")

		for _, branchPath := range branches {
			branchName := filepath.Base(branchPath)

			if branchName == currentBranch {
				fmt.Print("* ")
				fmt.Println(branchName)
			} else {
				fmt.Println("  " + branchName)
			}
		}
		return
	}

	headData, err := os.ReadFile(".minigit/HEAD")
	if err != nil {
		panic(err)
	}
	headRef := strings.TrimSpace(string(headData))
	if !strings.HasPrefix(headRef, "ref: ") {
		panic("Invalid HEAD format")
	}
	refPath := strings.TrimPrefix(headRef, "ref: ")

	commitHashBytes, err := os.ReadFile(".minigit/" + refPath)
	if err != nil {
		panic(err)
	}
	commitHash := strings.TrimSpace(string(commitHashBytes))

	err = os.WriteFile(".minigit/refs/heads/"+name, []byte(commitHash), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Branch '%s' created at commit %s\n", name, commitHash)
}
