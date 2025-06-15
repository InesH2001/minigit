package core

import (
	"fmt"
	"os"
	"strings"
	"minigit/utils"
)

func Revert(commitHash string, author string) error {
	treeHash, message, err := loadCommitMeta(commitHash)
	if err != nil {
		return err
	}

	tree, err := readTree(treeHash)
	if err != nil {
		return err
	}

	index, err := restoreFiles(tree)
	if err != nil {
		return err
	}

	if err := utils.WriteIndex(index); err != nil {
		return fmt.Errorf("failed to update index: %w", err)
	}

	return Commit(CommitParams{
		Author:  author,
		Message: fmt.Sprintf("Revert \"%s\"", message),
	})
}

func loadCommitMeta(commitHash string) (treeHash string, message string, err error) {
	commitPath := ".miniGit/objects/commits/" + commitHash
	data, err := os.ReadFile(commitPath)
	if err != nil {
		return "", "", fmt.Errorf("commit not found: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "tree: ") {
			treeHash = strings.TrimPrefix(line, "tree: ")
		}
		if strings.HasPrefix(line, "message: ") {
			message = strings.TrimPrefix(line, "message: ")
		}
	}
	return treeHash, message, nil
}

func readTree(treeHash string) (map[string]string, error) {
	path := ".miniGit/objects/trees/" + treeHash
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read tree: %w", err)
	}

	tree := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		tree[parts[0]] = parts[1]
	}
	return tree, nil
}

func restoreFiles(tree map[string]string) (map[string]string, error) {
	index := make(map[string]string)
	for path, hash := range tree {
		content := utils.GetBlobContent(hash)
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to restore %s: %w", path, err)
		}
		index[path] = hash
	}
	return index, nil
}
