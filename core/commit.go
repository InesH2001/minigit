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
	if utils.FileExists(".miniGit/MERGE_HEAD") {
    	_ = os.Remove(".miniGit/MERGE_HEAD")
	}

	index, err := utils.ReadIndex()
	if err != nil {
		return err
	}

	if len(index) == 0 {
		return fmt.Errorf("no files staged for commit")
	}
	
	headRef, parentHash, err := utils.GetCurrentBranchAndParentCommitHash()
	if err != nil {
		return err
	}

	treeContent := buildTree(index, parentHash)
	treeHash := utils.HashContent([]byte(treeContent))

	oldTreeHash := utils.GetTreeHashFromCommit(parentHash)
	if treeHash == oldTreeHash {
		return fmt.Errorf("no changes since last commit")
	}

	if err := utils.WriteTree(treeHash, treeContent); err != nil {
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

func buildTree(index map[string]string, parentHash string) string {
	previousTree := utils.ReadTreeFromCommit(parentHash)

	for file, hash := range index {
		if hash == "" {
			delete(previousTree, file)
		} else {
			previousTree[file] = hash
		}
	}

	var builder strings.Builder
	for path, hash := range previousTree {
		builder.WriteString(fmt.Sprintf("%s %s\n", path, hash))
	}
	return builder.String()
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
