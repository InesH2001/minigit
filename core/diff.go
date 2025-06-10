package core

import (
	"fmt"
	"os"
	"strings"

	"minigit/utils"
)

func Diff() error {
	index := utils.ReadIndex()

	for filePath, hash := range index {
		currentContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Reading impossible for %s : %v\n", filePath, err)
			continue
		}

		blobPath := ".miniGit/objects/blobs/" + hash
		indexedContent, err := os.ReadFile(blobPath)
		if err != nil {
			fmt.Printf("Blob not find for %s\n", filePath)
			continue
		}

		if string(currentContent) != string(indexedContent) {
			fmt.Printf("Diff for %s :\n", filePath)
			showLineDiff(string(indexedContent), string(currentContent))
		}
	}

	return nil
}

func showLineDiff(oldText, newText string) {
	oldLines := strings.Split(oldText, "\n")
	newLines := strings.Split(newText, "\n")
	max := maxLen(oldLines, newLines)

	for i := 0; i < max; i++ {
		var oldLine, newLine string
		if i < len(oldLines) {
			oldLine = oldLines[i]
		}
		if i < len(newLines) {
			newLine = newLines[i]
		}
		if oldLine != newLine {
			if oldLine != "" {
				fmt.Println("- " + oldLine)
			}
			if newLine != "" {
				fmt.Println("+ " + newLine)
			}
		}
	}
}

func maxLen(a, b []string) int {
	if len(a) > len(b) {
		return len(a)
	}
	return len(b)
}
