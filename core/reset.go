package core

import (
	"bufio"
	"fmt"
	"minigit/utils"
	"os"
	"strings"
)

func Reset(commitHash string) error {
	commitPath := ".minigit/objects/commits/" + commitHash
	data, err := os.ReadFile(commitPath)
	if err != nil {
		return fmt.Errorf("commit %s not found (%v)", commitHash, err)
	}

	var treeHash string
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "tree: ") {
			treeHash = strings.TrimPrefix(line, "tree: ")
			break
		}
	}

	if treeHash == "" {
		return fmt.Errorf("tree not found in commit %s", commitHash)
	}

	treePath := ".minigit/objects/trees/" + treeHash
	treeData, err := os.ReadFile(treePath)
	if err != nil {
		return fmt.Errorf("tree %s not found (%v)", treeHash, err)
	}

	lines := strings.Split(string(treeData), "\n")
	newIndex := make(map[string]string)

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		filePath := parts[0]
		hash := parts[1]

		content, err := os.ReadFile(".minigit/objects/blobs/" + hash)
		if err != nil {
			return fmt.Errorf("missing blob: %s (%v)", hash, err)
		}

		if err := utils.WriteFile(filePath, content); err != nil {
			return fmt.Errorf("failed to write file %s (%v)", filePath, err)
		}

		newIndex[filePath] = hash
	}

	if err := utils.WriteIndex(newIndex); err != nil {
		return fmt.Errorf("failed to update index: %v", err)
	}

	headData, err := os.ReadFile(".minigit/HEAD")
	if err != nil {
		return fmt.Errorf("failed to read HEAD (%v)", err)
	}
	ref := strings.TrimPrefix(strings.TrimSpace(string(headData)), "ref: ")
	if err := utils.WriteFile(".minigit/"+ref, []byte(commitHash)); err != nil {
		return fmt.Errorf("failed to update HEAD ref: %v", err)
	}

	fmt.Println("Reset successful to", commitHash)
	return nil
}
