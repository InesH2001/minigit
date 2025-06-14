package cmd

import (
	"fmt"
	"minigit/core"
	"os"
)

func RunReset() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go reset <commit-hash>")
		return
	}
	commitHash := os.Args[2]
	if err := core.Reset(commitHash); err != nil {
		fmt.Println("Error during reset:", err)
	}
}
