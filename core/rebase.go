package core

import (
	"fmt"
	"minigit/utils"
	"os"
	"strings"
)

func getParentsFromCommitLocal(commitHash string) []string {
	data, err := os.ReadFile(".miniGit/objects/commits/" + commitHash)
	if err != nil {
		return []string{}
	}

	lines := strings.Split(string(data), "\n")
	var parents []string

	for _, line := range lines {
		if strings.HasPrefix(line, "parent: ") {
			parentHash := strings.TrimSpace(strings.TrimPrefix(line, "parent: "))
			if parentHash != "" {
				parents = append(parents, parentHash)
			}
		}
	}

	return parents
}

func Rebase(targetBranch string) error {
	if !utils.BranchExists(targetBranch) {
		return fmt.Errorf("branch '%s' does not exist", targetBranch)
	}

	currentBranch := utils.GetCurrentBranchName()
	if currentBranch == "" {
		return fmt.Errorf("failed to get current branch")
	}

	if currentBranch == targetBranch {
		return fmt.Errorf("cannot rebase a branch onto itself")
	}

	if utils.HasUncommittedChanges() {
		return fmt.Errorf("you have uncommitted changes. Please commit or stash them before rebasing")
	}

	currentCommitHash, err := utils.GetCommitHashFromRef("refs/heads/" + currentBranch)
	if err != nil {
		return fmt.Errorf("failed to get current commit: %w", err)
	}

	targetCommitHash, err := utils.GetCommitHashFromRef("refs/heads/" + targetBranch)
	if err != nil {
		return fmt.Errorf("failed to get target commit: %w", err)
	}

	commonAncestor := findCommonAncestor(currentCommitHash, targetCommitHash)
	if commonAncestor == "" {
		return fmt.Errorf("no common ancestor found")
	}

	if commonAncestor == targetCommitHash {
		fmt.Println("Current branch is already up to date with", targetBranch)
		return nil
	}

	if commonAncestor == currentCommitHash {
		fmt.Printf("Fast-forwarding %s to %s\n", currentBranch, targetBranch)
		branchRefPath := ".miniGit/refs/heads/" + currentBranch
		err = utils.WriteFile(branchRefPath, []byte(targetCommitHash))
		if err != nil {
			return fmt.Errorf("failed to update branch reference: %w", err)
		}

		err = restoreWorkingDirectory(targetCommitHash)
		if err != nil {
			return fmt.Errorf("failed to restore working directory: %w", err)
		}

		return nil
	}

	commits := collectCommitsSince(currentCommitHash, commonAncestor)
	if len(commits) == 0 {
		fmt.Println("No commits to rebase")
		return nil
	}

	fmt.Printf("Rebasing %d commits from %s onto %s\n", len(commits), currentBranch, targetBranch)

	err = saveRebaseState(targetBranch, commits)
	if err != nil {
		return fmt.Errorf("failed to save rebase state: %w", err)
	}

	err = restoreWorkingDirectory(targetCommitHash)
	if err != nil {
		return fmt.Errorf("failed to restore working directory: %w", err)
	}

	err = utils.WriteFile(".miniGit/index", []byte(""))
	if err != nil {
		return fmt.Errorf("failed to clear index: %w", err)
	}
	branchRefPath := ".miniGit/refs/heads/" + currentBranch
	err = utils.WriteFile(branchRefPath, []byte(targetCommitHash))
	if err != nil {
		return fmt.Errorf("failed to move current branch to base of rebase: %w", err)
	}
	for i := len(commits) - 1; i >= 0; i-- {
		commit := commits[i]
		fmt.Printf("Applying: %s\n", commit.Message)
		err = applyCommit(commit)
		if err != nil {
			return fmt.Errorf("failed to apply commit %s: %w", commit.Hash, err)
		}
	}

	cleanRebaseState()
	fmt.Printf("Successfully rebased %s onto %s\n", currentBranch, targetBranch)
	return nil
}

type CommitInfo struct {
	Hash    string
	Message string
	Author  string
	Tree    map[string]string
}

func collectCommitsSince(startCommit, endCommit string) []CommitInfo {
	var commits []CommitInfo
	current := startCommit

	for current != endCommit && current != "" {
		commit := loadCommitInfo(current)
		if commit.Hash != "" {
			commits = append(commits, commit)
		}

		parents := getParentsFromCommitLocal(current)
		if len(parents) > 0 {
			current = parents[0]
		} else {
			break
		}
	}

	return commits
}

func loadCommitInfo(commitHash string) CommitInfo {
	commitPath := ".miniGit/objects/commits/" + commitHash
	data, err := os.ReadFile(commitPath)
	if err != nil {
		return CommitInfo{}
	}

	lines := strings.Split(string(data), "\n")
	commit := CommitInfo{Hash: commitHash}

	for _, line := range lines {
		if strings.HasPrefix(line, "message: ") {
			commit.Message = strings.TrimPrefix(line, "message: ")
		} else if strings.HasPrefix(line, "author: ") {
			commit.Author = strings.TrimPrefix(line, "author: ")
		} else if strings.HasPrefix(line, "tree: ") {
			treeHash := strings.TrimPrefix(line, "tree: ")
			commit.Tree = loadTreeFromHash(treeHash)
		}
	}

	return commit
}

func loadTreeFromHash(treeHash string) map[string]string {
	treePath := ".miniGit/objects/trees/" + treeHash
	data, err := os.ReadFile(treePath)
	if err != nil {
		return make(map[string]string)
	}

	tree := make(map[string]string)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			tree[parts[0]] = parts[1]
		}
	}

	return tree
}

func applyCommit(commit CommitInfo) error {
	err := utils.WriteFile(".miniGit/index", []byte(""))
	if err != nil {
		return fmt.Errorf("failed to clear index: %w", err)
	}
	baseFiles := make(map[string]string)
	workingFiles := utils.ListFiles()
	for _, filePath := range workingFiles {
		content, err := os.ReadFile(filePath)
		if err == nil {
			baseFiles[filePath] = string(content)
		}
	}
	baseOfRebaseFiles := make(map[string]bool)
	for filePath := range baseFiles {
		if !strings.Contains(filePath, "feature") {
			baseOfRebaseFiles[filePath] = true
		}
	}

	hasChanges := false
	for filePath, blobHash := range commit.Tree {
		commitContent := utils.GetBlobContent(blobHash)
		currentContent := baseFiles[filePath]
		if baseOfRebaseFiles[filePath] {
		} else {
			if currentContent != commitContent {
				err := utils.WriteFile(filePath, []byte(commitContent))
				if err != nil {
					return fmt.Errorf("failed to update file %s: %w", filePath, err)
				}
				hasChanges = true
			}
		}
		err = Add(filePath)
		if err != nil {
			return fmt.Errorf("failed to add file %s: %w", filePath, err)
		}
	}

	for _, filePath := range workingFiles {
		if _, exists := commit.Tree[filePath]; !exists {
			err = Add(filePath)
			if err != nil {
			}
		}
	}
	if !hasChanges {
		newIndex, err := utils.ReadIndex()
		if err != nil {
			return fmt.Errorf("failed to read new index: %w", err)
		}

		lastCommit := utils.ReadLastCommit()
		indexChanged := false

		for file, hash := range newIndex {
			if lastHash, ok := lastCommit[file]; !ok || lastHash != hash {
				indexChanged = true
				break
			}
		}

		if !indexChanged {
			for file := range lastCommit {
				if _, exists := newIndex[file]; !exists {
					indexChanged = true
					break
				}
			}
		}

		if !indexChanged {
			return nil
		}
	}
	return Commit(CommitParams{
		Message: commit.Message,
		Author:  commit.Author,
	})
}

func findCommonAncestor(commit1, commit2 string) string {
	if commit1 == commit2 {
		return commit1
	}
	ancestors1 := make(map[string]bool)
	collectAncestors(commit1, ancestors1)
	return findFirstCommonAncestor(commit2, ancestors1)
}

func collectAncestors(commitHash string, ancestors map[string]bool) {
	if commitHash == "" || ancestors[commitHash] {
		return
	}

	ancestors[commitHash] = true
	parents := getParentsFromCommitLocal(commitHash)

	for _, parent := range parents {
		if parent != "" {
			collectAncestors(parent, ancestors)
		}
	}
}

func findFirstCommonAncestor(commitHash string, ancestors map[string]bool) string {
	if commitHash == "" {
		return ""
	}

	if ancestors[commitHash] {
		return commitHash
	}

	parents := getParentsFromCommitLocal(commitHash)

	for _, parent := range parents {
		if parent != "" {
			result := findFirstCommonAncestor(parent, ancestors)
			if result != "" {
				return result
			}
		}
	}

	return ""
}

func saveRebaseState(targetBranch string, commits []CommitInfo) error {
	rebaseDir := ".miniGit/rebase"
	err := os.MkdirAll(rebaseDir, 0755)
	if err != nil {
		return err
	}

	err = utils.WriteFile(".miniGit/rebase/target", []byte(targetBranch))
	if err != nil {
		return err
	}

	var content strings.Builder
	for _, commit := range commits {
		content.WriteString(fmt.Sprintf("%s %s\n", commit.Hash, commit.Message))
	}

	return utils.WriteFile(".miniGit/rebase/commits", []byte(content.String()))
}

func cleanRebaseState() {
	os.RemoveAll(".miniGit/rebase")
}

func RebaseAbort() error {
	if !utils.FileExists(".miniGit/rebase") {
		return fmt.Errorf("no rebase in progress")
	}

	currentBranch := utils.GetCurrentBranchName()
	if currentBranch == "" {
		return fmt.Errorf("failed to get current branch")
	}

	commitHash, err := utils.GetCommitHashFromRef("refs/heads/" + currentBranch)
	if err != nil {
		return fmt.Errorf("failed to get original commit: %w", err)
	}

	err = restoreWorkingDirectory(commitHash)
	if err != nil {
		return fmt.Errorf("failed to restore working directory: %w", err)
	}

	cleanRebaseState()
	fmt.Println("Rebase aborted. Restored to original state.")
	return nil
}
