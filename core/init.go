package core

import (
	"fmt"
	"os"
	"path/filepath"
)

func InitRepo() error {
	miniGitPath := ".miniGit"

	if _, err := os.Stat(miniGitPath); !os.IsNotExist(err) {
		return fmt.Errorf("miniGit already initialized")
	}

	dirs := []string{
		filepath.Join(miniGitPath, "objects"),
		filepath.Join(miniGitPath, "objects/blobs"),
		filepath.Join(miniGitPath, "objects/commits"),
		filepath.Join(miniGitPath, "objects/trees"),
		filepath.Join(miniGitPath, "refs", "heads"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %v", dir, err)
		}
	}

	headPath := filepath.Join(miniGitPath, "HEAD")
	if err := os.WriteFile(headPath, []byte("ref: refs/heads/main"), 0644); err != nil {
		return fmt.Errorf("error creating HEAD file: %v", err)
	}

	indexPath := filepath.Join(miniGitPath, "index")
	if err := os.WriteFile(indexPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("error creating index file: %v", err)
	}

	if err := StoreUsername(); err != nil {
		return fmt.Errorf("failed to store username: %v", err)
	}

	return nil
}
