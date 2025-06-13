package core

import (
	"fmt"
	"minigit/utils"
	"os"
)

func Status() error {
	index := utils.ReadIndex()
	lastCommit := utils.ReadLastCommit()
	workingFiles := utils.ListFiles()

	currentBranch := utils.GetCurrentBranch()
	fmt.Printf("On branch %s\n\n", currentBranch)

	fmt.Println("Changes to be committed:")
	for filePath, hash := range index {
		if lastHash, ok := lastCommit[filePath]; !ok || lastHash != hash {
			fmt.Printf("  modified:   %s\n", filePath)
		}
	}
	fmt.Println()

	fmt.Println("Changes not staged for commit:")
	for _, filePath := range workingFiles {
		if indexHash, ok := index[filePath]; ok {
			content, _ := os.ReadFile(filePath)
			if utils.HashContent(content) != indexHash {
				fmt.Printf("  modified:   %s\n", filePath)
			}
		}
	}
	fmt.Println()

	fmt.Println("Untracked files:")
	for _, filePath := range workingFiles {
		if _, tracked := index[filePath]; !tracked {
			fmt.Printf("  %s\n", filePath)
		}
	}

	return nil
}
