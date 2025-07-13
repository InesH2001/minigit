package core

import (
	"fmt"
	"minigit/utils"
	"os"
	"path/filepath"
	"strings"
)

func Add(filePath string) error {

	if isInternalPath(filePath) {
		fmt.Printf("Ignore : %s (internal folder)\n", filePath)
		return nil
	}

	if utils.IsIgnored(filePath) {
		fmt.Printf("Ignore : %s (gitignore)\n", filePath)
		return nil
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || isInternalPath(path) {
				return nil
			}
			return Add(path)
		})
	}

	if !utils.FileExists(filePath) {
		fmt.Printf("Ignore : %s (does not exist)\n", filePath)
		return nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	hash := utils.HashContent(content)
	index, err := utils.ReadIndex()
	if err != nil {
		return err
	}

	if currentHash, ok := index[filePath]; ok && currentHash == hash {
		fmt.Printf("Unchanged : %s\n", filePath)
		return nil
	}

	blobPath := filepath.Join(".miniGit", "objects", "blobs", hash)
	if !utils.FileExists(blobPath) {
		compressedContent, err := utils.Compress(content)
		if err != nil {
			return fmt.Errorf("compression failed for %s: %v", filePath, err)
		}

		err = os.WriteFile(blobPath, compressedContent, 0644)
		if err != nil {
			return err
		}
		fmt.Printf("Blob created (compressed): %s\n", hash)
	}

	index[filePath] = hash
	err = utils.WriteIndex(index)
	if err != nil {
		return err
	}
	fmt.Printf("Added : %s\n", filePath)
	return nil
}

func isInternalPath(path string) bool {
	return strings.HasPrefix(path, ".git") ||
		strings.HasPrefix(path, ".miniGit") ||
		strings.Contains(path, "/.git/") ||
		strings.Contains(path, "/.miniGit/")
}
