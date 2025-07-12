package core

import (
	"fmt"
	"minigit/utils"
	"os"
	"strings"
)

func Revert(commitHash string, author string) error {
	treeHash, parentCommitHash, message, err := loadCommitMeta(commitHash)
	if err != nil {
		return err
	}

	targetTree, err := readTree(treeHash)
	if err != nil {
		return err
	}

	parentTreeHash, _, _, err := loadCommitMeta(parentCommitHash)
	if err != nil {
		return err
	}

	parentTree, err := readTree(parentTreeHash)
	if err != nil {
		return err
	}

	diff := computeInverseDiff(parentTree, targetTree)
	index, err := applyDiff(diff)
	if err != nil {
		return err
	}

	if err := utils.WriteIndex(index); err != nil {
		return fmt.Errorf("failed to update index: %w", err)
	}

	return Commit(CommitParams{
		Author:  author,
		Message: fmt.Sprintf("Revert \"%s\"", strings.ReplaceAll(message, "\"", "'")),
	})
}

func loadCommitMeta(commitHash string) (treeHash string, parentHash string, message string, err error) {
	commitPath := ".miniGit/objects/commits/" + commitHash
	data, err := os.ReadFile(commitPath)
	if err != nil {
		return "", "", "", fmt.Errorf("commit not found: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "tree: ") {
			treeHash = strings.TrimPrefix(line, "tree: ")
		}
		if strings.HasPrefix(line, "parent: ") {
            parentHash = strings.TrimPrefix(line, "parent: ")
        }
		if strings.HasPrefix(line, "message: ") {
			message = strings.TrimPrefix(line, "message: ")
		}
	}
	return treeHash, parentHash, message, nil
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

func computeInverseDiff(parentTree, targetTree map[string]string) map[string]*string {
	diff := make(map[string]*string)

	for path, targetHash := range targetTree {
		parentHash, exists := parentTree[path]
		if !exists {
			diff[path] = nil
		} else if parentHash != targetHash {
			diff[path] = &parentHash
		}
	}

	for path, parentHash := range parentTree {
		if _, exists := targetTree[path]; !exists {
			diff[path] = &parentHash
		}
	}

	return diff
}

func applyDiff(diff map[string]*string) (map[string]string, error) {
	index := make(map[string]string)

	for path, blobHash := range diff {
		if blobHash == nil {
			err := os.Remove(path)
			if err != nil && !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to delete %s: %w", path, err)
			}
		} else {
			content, err := utils.ReadAndDecompressBlob(*blobHash)
			if err != nil {
				return nil, fmt.Errorf("failed to read blob %s: %w", *blobHash, err)
			}
			err = os.WriteFile(path, content, 0644)
			if err != nil {
				return nil, fmt.Errorf("failed to restore file %s: %w", path, err)
			}
			index[path] = *blobHash
		}
	}

	return index, nil
}
