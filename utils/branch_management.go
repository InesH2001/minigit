package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func RenameBranch(oldName, newName string) error {
	if !BranchExists(oldName) {
		return fmt.Errorf("branch '%s' does not exist", oldName)
	}
	if BranchExists(newName) {
		return fmt.Errorf("branch '%s' already exists", newName)
	}
	oldPath := filepath.Join(".miniGit", "refs", "heads", oldName)
	newPath := filepath.Join(".miniGit", "refs", "heads", newName)
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return fmt.Errorf("failed to rename branch reference: %w", err)
	}
	currentBranch := GetCurrentBranchName()
	if currentBranch == oldName {
		headContent := "ref: refs/heads/" + newName
		err = WriteFile(".miniGit/HEAD", []byte(headContent))
		if err != nil {
			os.Rename(newPath, oldPath)
			return fmt.Errorf("failed to update HEAD: %w", err)
		}
	}

	return nil
}

func DeleteBranch(branchName string, force bool) error {
	if !BranchExists(branchName) {
		return fmt.Errorf("branch '%s' does not exist", branchName)
	}
	currentBranch := GetCurrentBranchName()
	if currentBranch == branchName {
		return fmt.Errorf("cannot delete the currently active branch '%s'", branchName)
	}
	if !force {
		isMerged, err := IsBranchMerged(branchName, currentBranch)
		if err != nil {
			return fmt.Errorf("failed to check if branch is merged: %w", err)
		}
		if !isMerged {
			return fmt.Errorf("branch '%s' is not fully merged. Use -f to force deletion", branchName)
		}
	}

	branchPath := filepath.Join(".miniGit", "refs", "heads", branchName)
	err := os.Remove(branchPath)
	if err != nil {
		return fmt.Errorf("failed to delete branch reference: %w", err)
	}

	return nil
}

func IsBranchMerged(branchToCheck, targetBranch string) (bool, error) {
	branchCommit, err := GetCommitHashFromRef("refs/heads/" + branchToCheck)
	if err != nil {
		return false, err
	}

	targetCommit, err := GetCommitHashFromRef("refs/heads/" + targetBranch)
	if err != nil {
		return false, err
	}
	return IsCommitAncestor(branchCommit, targetCommit), nil
}

func IsCommitAncestor(commitA, commitB string) bool {
	if commitA == commitB {
		return true
	}

	visited := make(map[string]bool)
	queue := []string{commitB}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true

		if current == commitA {
			return true
		}

		parents := GetParentsFromCommit(current)
		queue = append(queue, parents...)
	}

	return false
}
