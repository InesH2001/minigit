package core

import (
	"fmt"
	"minigit/utils"
	"os"
	"strings"
	"time"
)

type CommitParams struct {
	Message string
	Author  string
}

func Commit(params CommitParams) error {
	index, err := readIndexFile()
	if err != nil {
		return err
	}

	treeContent := buildTree(index)
	treeHash := utils.HashContent([]byte(treeContent))

	if err := writeTree(treeHash, treeContent); err != nil {
		return err
	}

	headRef, parentHash, err := readHeadAndParent()
	if err != nil {
		return err
	}

	commitContent := buildCommit(params, treeHash, parentHash)
	commitHash := utils.HashContent([]byte(commitContent))

	if err := writeCommit(commitHash, commitContent); err != nil {
		return err
	}

	if err := utils.WriteFile(".miniGit/"+headRef, []byte(commitHash)); err != nil {
		return fmt.Errorf("failed to update branch ref: %w", err)
	}

	if err := utils.WriteFile(".miniGit/index", []byte("")); err != nil {
		return fmt.Errorf("failed to clear index: %w", err)
	}

	fmt.Println("Commit successful:", commitHash)
	return nil
}

func readIndexFile() (map[string]string, error) {
	content, err := os.ReadFile(".miniGit/index")
	if err != nil {
		return nil, fmt.Errorf("failed to read index: %w", err)
	}
	lines := strings.Split(string(content), "\n")
	index := make(map[string]string)
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		index[parts[0]] = parts[1]
	}
	return index, nil
}

func buildTree(index map[string]string) string {
	var builder strings.Builder
	for path, hash := range index {
		builder.WriteString(fmt.Sprintf("%s %s\n", path, hash))
	}
	return builder.String()
}

func writeTree(treeHash, content string) error {
	return utils.WriteFile(".miniGit/objects/trees/"+treeHash, []byte(content))
}

func readHeadAndParent() (string, string, error) {
	headContent, err := os.ReadFile(".miniGit/HEAD")
	if err != nil {
		return "", "", fmt.Errorf("failed to read HEAD: %w", err)
	}
	headRef := strings.TrimPrefix(string(headContent), "ref: ")
	headRef = strings.TrimSpace(headRef)

	parentHash := ""
	parentPath := ".miniGit/" + headRef
	if data, err := os.ReadFile(parentPath); err == nil {
		parentHash = strings.TrimSpace(string(data))
	}
	return headRef, parentHash, nil
}

func buildCommit(params CommitParams, treeHash, parentHash string) string {
	return fmt.Sprintf(
		"tree: %s\nparent: %s\nauthor: %s\ndate: %s\nmessage: %s\n",
		treeHash,
		parentHash,
		params.Author,
		time.Now().Format(time.RFC3339),
		params.Message,
	)
}

func writeCommit(hash, content string) error {
	return utils.WriteFile(".miniGit/objects/commits/"+hash, []byte(content))
}
