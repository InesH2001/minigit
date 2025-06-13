package core

import (
	"bufio"
	"fmt"
	"minigit/utils"
	"os"
	"strings"
)

func Reset(commitHash string) error {
	fmt.Println("→ Start Reset with commitHash:", commitHash)

	commitPath := ".minigit/objects/commits/" + commitHash
	fmt.Println("→ Reading commit file at:", commitPath)
	data, err := os.ReadFile(commitPath)
	if err != nil {
		return fmt.Errorf("❌ commit %s not found (%v)", commitHash, err)
	}

	var treeHash string
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("→ Scanning commit line:", line)
		if strings.HasPrefix(line, "tree: ") {
			treeHash = strings.TrimPrefix(line, "tree: ")
			fmt.Println("→ Found treeHash:", treeHash)
			break
		}
	}

	if treeHash == "" {
		return fmt.Errorf("❌ tree not found in commit %s", commitHash)
	}

	treePath := ".minigit/objects/trees/" + treeHash
	fmt.Println("→ Reading tree file at:", treePath)
	treeData, err := os.ReadFile(treePath)
	if err != nil {
		return fmt.Errorf("❌ tree %s not found (%v)", treeHash, err)
	}

	lines := strings.Split(string(treeData), "\n")
	newIndex := make(map[string]string)

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			fmt.Println("→ Invalid tree line (skipped):", line)
			continue
		}
		filePath := parts[0]
		hash := parts[1]
		fmt.Printf("→ Restoring file %s from blob %s\n", filePath, hash)

		blobPath := ".minigit/objects/blobs/" + hash
		fmt.Println("→ Reading blob at:", blobPath)
		content, err := os.ReadFile(blobPath)
		if err != nil {
			return fmt.Errorf("❌ missing blob: %s (%v)", hash, err)
		}

		fmt.Println("→ Writing content to:", filePath)
		if err := utils.WriteFile(filePath, content); err != nil {
			return fmt.Errorf("❌ failed to write file %s (%v)", filePath, err)
		}

		newIndex[filePath] = hash
	}

	fmt.Println("→ Updating index...")
	if err := utils.WriteIndex(newIndex); err != nil {
		return fmt.Errorf("❌ failed to update index: %v", err)
	}

	fmt.Println("→ Reading HEAD...")
	headData, err := os.ReadFile(".minigit/HEAD")
	if err != nil {
		return fmt.Errorf("❌ failed to read HEAD (%v)", err)
	}
	ref := strings.TrimPrefix(strings.TrimSpace(string(headData)), "ref: ")
	fmt.Println("→ Writing new HEAD ref to .minigit/" + ref)
	if err := utils.WriteFile(".minigit/"+ref, []byte(commitHash)); err != nil {
		return fmt.Errorf("❌ failed to update HEAD ref: %v", err)
	}

	fmt.Println("✅ Reset successful to", commitHash)
	return nil
}
