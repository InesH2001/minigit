package core

import (
	"fmt"
	"minigit/utils"
	"os"
	"strings"
)

func Merge(branchName string) error {
	hasConflict := false
	modifiedFiles := []string{}

	currentBranch, err := utils.GetCurrentBranchRefPath()
	if err != nil {
		return err
	}

	headCommitHash, err := utils.GetCommitHashFromRef(currentBranch)
	if err != nil {
		return err
	}

	if branchName == utils.GetCurrentBranchName() {
		return fmt.Errorf("cannot merge a branch with itself")
	}

	branchCommitHash, err := utils.GetCommitHashFromRef("refs/heads/" + branchName)
	if err != nil {
		return err
	}

	if utils.AreTreesEqual(headCommitHash, branchCommitHash) {
		fmt.Println("Branches are identical. Nothing to merge.")
		return nil
	}

	err = os.WriteFile(".miniGit/MERGE_HEAD", []byte(headCommitHash), 0644)
	if err != nil {
		return fmt.Errorf("failed to save merge state: %w", err)
	}

	commonAncestorHash := findCommonCommitAncestorHash(headCommitHash, branchCommitHash)

	headBlobs := utils.ReadTreeFromCommit(headCommitHash)
	branchBlobs := utils.ReadTreeFromCommit(branchCommitHash)
	baseBlobs := utils.ReadTreeFromCommit(commonAncestorHash)

	allFiles := utils.GetUniqueUnionKeys(headBlobs, branchBlobs, baseBlobs)

	for _, file := range allFiles {
		base := utils.GetBlobContent(baseBlobs[file])
		head := utils.GetBlobContent(headBlobs[file])
		branch := utils.GetBlobContent(branchBlobs[file])

		merged := mergeThreeVersions(base, head, branch, branchName)

		err := os.WriteFile(file, []byte(merged), 0644)
		if err != nil {
			return err
		}

		if strings.Contains(merged, "<<<<<<<") {
			hasConflict = true
			fmt.Printf("⚠️ Conflict detected in %s. Please resolve manually. ⚠️\n", file)
		} else if merged != base {
			fmt.Printf("Merge completed without conflict for %s.\n", file)
			modifiedFiles = append(modifiedFiles, file)
		}

	}

	if !hasConflict {
		if len(modifiedFiles) == 0 {
			fmt.Println("No changes to merge. The branch is already up to date with", branchName)
			return nil
		}

		for _, file := range modifiedFiles {
			if err := Add(file); err != nil {
				return fmt.Errorf("failed to add file %s: %w", file, err)
			}
		}

		err := Commit(CommitParams{
			Message: "Merge branch '" + branchName + "'",
			Author:  "MiniGit",
		})
		if err != nil {
			return fmt.Errorf("failed to commit merge: %w", err)
		}

		if err := os.Remove(".miniGit/MERGE_HEAD"); err != nil {
			return fmt.Errorf("failed to remove MERGE_HEAD: %w", err)
		}
	}
	return nil
}

func mergeThreeVersions(base, head, branch, branchName string) string {
	baseLines := strings.Split(base, "\n")
	headLines := strings.Split(head, "\n")
	branchLines := strings.Split(branch, "\n")

	maxLen := utils.Max(len(baseLines), len(headLines), len(branchLines))
	var result strings.Builder

	for i := 0; i < maxLen; i++ {
		baseLine := utils.GetLine(baseLines, i)
		headLine := utils.GetLine(headLines, i)
		branchLine := utils.GetLine(branchLines, i)

		switch {
		case headLine == branchLine:
			result.WriteString(headLine + "\n")
		case headLine != baseLine && branchLine == baseLine:
			result.WriteString(headLine + "\n")
		case branchLine != baseLine && headLine == baseLine:
			result.WriteString(branchLine + "\n")
		case headLine == "" && branchLine != "":
			result.WriteString(branchLine + "\n")
		case branchLine == "" && headLine != "":
			result.WriteString(headLine + "\n")
		default:
			result.WriteString("<<<<<<< HEAD\n")
			result.WriteString(headLine + "\n")
			result.WriteString("=======\n")
			result.WriteString(branchLine + "\n")
			result.WriteString(">>>>>>> " + branchName + "\n")
		}
	}

	return result.String()
}

func findCommonCommitAncestorHash(commit1, commit2 string) string {
	visited := map[string]bool{}
	queue := []string{commit1}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		visited[current] = true

		parents := utils.GetParentsFromCommit(current)
		queue = append(queue, parents...)
	}

	queue = []string{commit2}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			return current
		}

		parents := utils.GetParentsFromCommit(current)
		queue = append(queue, parents...)
	}

	return ""
}

func MergeAbort() error {
	data, err := os.ReadFile(".miniGit/MERGE_HEAD")
	if err != nil {
		return fmt.Errorf("No merge in progress or unable to read MERGE_HEAD: %w", err)
	}

	commitHash := strings.TrimSpace(string(data))

	tree := utils.ReadTreeFromCommit(commitHash)
	for file, blobHash := range tree {
		content := utils.GetBlobContent(blobHash)
		if err := os.WriteFile(file, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to restore file %s: %w", file, err)
		}
	}

	if err := os.Remove(".miniGit/MERGE_HEAD"); err != nil {
		return fmt.Errorf("failed to remove MERGE_HEAD: %w", err)
	}

	if err := utils.WriteFile(".miniGit/index", []byte("")); err != nil {
		return fmt.Errorf("commit succeeded, but failed to clear index after merge: %w", err)
	}

	fmt.Println("Merge aborted. Working directory restored to previous commit.")
	return nil
}
