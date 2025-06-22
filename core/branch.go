package core

import (
	"fmt"
	"minigit/utils"
	"os"
	"path/filepath"
	"strings"
)

func ListBranches() ([]string, error) {
	branchesPath := ".miniGit/refs/heads"

	if !utils.FileExists(branchesPath) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(branchesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read branches directory: %w", err)
	}

	var branches []string
	for _, entry := range entries {
		if !entry.IsDir() {
			branches = append(branches, entry.Name())
		}
	}

	return branches, nil
}

func GetCurrentBranchName() (string, error) {
	headData, err := os.ReadFile(".miniGit/HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD: %w", err)
	}

	headRef := strings.TrimSpace(string(headData))
	if !strings.HasPrefix(headRef, "ref: refs/heads/") {
		return "", fmt.Errorf("HEAD is not pointing to a branch")
	}

	return strings.TrimPrefix(headRef, "ref: refs/heads/"), nil
}

func CreateBranch(name string) error {
	if utils.BranchExists(name) {
		return fmt.Errorf("branch '%s' already exists", name)
	}

	currentCommitHash, err := utils.GetHeadCommit()
	if err != nil {
		return fmt.Errorf("failed to get current commit: %w", err)
	}

	branchRefPath := filepath.Join(".miniGit", "refs", "heads", name)
	err = utils.WriteFile(branchRefPath, []byte(currentCommitHash))
	if err != nil {
		return fmt.Errorf("failed to create branch reference: %w", err)
	}

	return nil
}

func SwitchToBranch(branchName string) error {
	if !utils.BranchExists(branchName) {
		return fmt.Errorf("branch '%s' does not exist", branchName)
	}

	currentBranch, err := GetCurrentBranchName()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	if currentBranch == branchName {
		return fmt.Errorf("already on branch '%s'", branchName)
	}

	if utils.HasUncommittedChanges() {
		return fmt.Errorf("you have uncommitted changes. Please commit or stash them before switching branches")
	}

	branchCommitHash, err := utils.GetCommitHashFromRef("refs/heads/" + branchName)
	if err != nil {
		return fmt.Errorf("failed to get branch commit: %w", err)
	}

	err = restoreWorkingDirectory(branchCommitHash)
	if err != nil {
		return fmt.Errorf("failed to restore working directory: %w", err)
	}

	headContent := "ref: refs/heads/" + branchName
	err = utils.WriteFile(".miniGit/HEAD", []byte(headContent))
	if err != nil {
		return fmt.Errorf("failed to update HEAD: %w", err)
	}
	err = utils.WriteFile(".miniGit/index", []byte(""))
	if err != nil {
		return fmt.Errorf("failed to clear index: %w", err)
	}

	return nil
}

func restoreWorkingDirectory(commitHash string) error {
	targetTree := utils.ReadTreeFromCommit(commitHash)
	currentBranch, err := GetCurrentBranchName()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	currentCommitHash, err := utils.GetCommitHashFromRef("refs/heads/" + currentBranch)
	if err != nil {
		return fmt.Errorf("failed to get current commit: %w", err)
	}

	currentTree := utils.ReadTreeFromCommit(currentCommitHash)
	for filePath := range currentTree {
		if _, exists := targetTree[filePath]; !exists {
			err := os.Remove(filePath)
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove file %s: %w", filePath, err)
			}
		}
	}
	currentIndex, err := utils.ReadIndex()
	if err != nil {
		return fmt.Errorf("failed to read current index: %w", err)
	}

	for filePath := range currentIndex {
		if _, exists := targetTree[filePath]; !exists {
			err := os.Remove(filePath)
			if err != nil && !os.IsNotExist(err) {

			}
		}
	}
	for filePath, blobHash := range targetTree {
		content := utils.GetBlobContent(blobHash)
		err := utils.WriteFile(filePath, []byte(content))
		if err != nil {
			return fmt.Errorf("failed to restore file %s: %w", filePath, err)
		}
	}

	return nil
}
