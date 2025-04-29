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
		filepath.Join(miniGitPath, "refs", "heads"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %v", dir, err)
		}
	}

	headPath := filepath.Join(miniGitPath, "HEAD")
	if err := os.WriteFile(headPath, []byte("ref: refs/heads/master"), 0644); err != nil {
		return fmt.Errorf("error creating HEAD file: %v", err)
	}

	fmt.Println("miniGit correctely initizalized")
	return nil
}
