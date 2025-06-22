package core

import (
	"fmt"
	"minigit/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type StashEntry struct {
	ID      int
	Message string
	Branch  string
	Time    time.Time
	Index   map[string]string
	WorkDir map[string]string
}

func Stash(message string) error {
	if !utils.HasUncommittedChanges() && !hasUnstagedChanges() {
		return fmt.Errorf("no changes to stash")
	}

	currentBranch := utils.GetCurrentBranchName()
	if currentBranch == "" {
		return fmt.Errorf("failed to get current branch")
	}

	if message == "" {
		message = fmt.Sprintf("WIP on %s", currentBranch)
	}

	stashEntry := StashEntry{
		ID:      getNextStashID(),
		Message: message,
		Branch:  currentBranch,
		Time:    time.Now(),
	}

	index, err := utils.ReadIndex()
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}
	stashEntry.Index = index

	workDir, err := captureWorkingDirectory()
	if err != nil {
		return fmt.Errorf("failed to capture working directory: %w", err)
	}
	stashEntry.WorkDir = workDir

	err = saveStash(stashEntry)
	if err != nil {
		return fmt.Errorf("failed to save stash: %w", err)
	}

	err = cleanWorkingDirectory()
	if err != nil {
		return fmt.Errorf("failed to clean working directory: %w", err)
	}

	fmt.Printf("Saved working directory and index state on %s: %s\n", currentBranch, message)
	return nil
}

func StashPop() error {
	stashes, err := listStashes()
	if err != nil {
		return fmt.Errorf("failed to list stashes: %w", err)
	}

	if len(stashes) == 0 {
		return fmt.Errorf("no stash entries found")
	}

	latestStash := stashes[0]

	err = restoreStash(latestStash)
	if err != nil {
		return fmt.Errorf("failed to restore stash: %w", err)
	}

	err = removeStash(latestStash.ID)
	if err != nil {
		return fmt.Errorf("failed to remove stash: %w", err)
	}

	fmt.Printf("Dropped stash@{0} (%s)\n", latestStash.Message)
	return nil
}

func StashList() error {
	stashes, err := listStashes()
	if err != nil {
		return fmt.Errorf("failed to list stashes: %w", err)
	}

	if len(stashes) == 0 {
		fmt.Println("No stash entries found.")
		return nil
	}

	for i, stash := range stashes {
		fmt.Printf("stash@{%d}: On %s: %s\n", i, stash.Branch, stash.Message)
	}

	return nil
}

func hasUnstagedChanges() bool {
	workingFiles := utils.ListFiles()
	index, _ := utils.ReadIndex()

	for _, filePath := range workingFiles {
		if utils.IsIgnored(filePath) {
			continue
		}

		if indexHash, ok := index[filePath]; ok {
			content, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}
			if utils.HashContent(content) != indexHash {
				return true
			}
		} else {
			if len(index) > 0 {
				return true
			}
		}
	}

	return false
}

func getNextStashID() int {
	stashes, _ := listStashes()
	return len(stashes)
}

func captureWorkingDirectory() (map[string]string, error) {
	workDir := make(map[string]string)
	workingFiles := utils.ListFiles()

	for _, filePath := range workingFiles {
		if utils.IsIgnored(filePath) {
			continue
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		hash := utils.HashContent(content)
		workDir[filePath] = hash

		blobPath := filepath.Join(".miniGit", "objects", "blobs", hash)
		if !utils.FileExists(blobPath) {
			utils.WriteFile(blobPath, content)
		}
	}

	return workDir, nil
}

func saveStash(stash StashEntry) error {
	stashDir := ".miniGit/stash"
	if err := os.MkdirAll(stashDir, 0755); err != nil {
		return err
	}

	stashFile := filepath.Join(stashDir, fmt.Sprintf("stash_%d", stash.ID))

	var content strings.Builder
	content.WriteString(fmt.Sprintf("message: %s\n", stash.Message))
	content.WriteString(fmt.Sprintf("branch: %s\n", stash.Branch))
	content.WriteString(fmt.Sprintf("time: %s\n", stash.Time.Format(time.RFC3339)))
	content.WriteString("index:\n")

	for path, hash := range stash.Index {
		content.WriteString(fmt.Sprintf("  %s %s\n", path, hash))
	}

	content.WriteString("workdir:\n")
	for path, hash := range stash.WorkDir {
		content.WriteString(fmt.Sprintf("  %s %s\n", path, hash))
	}

	return utils.WriteFile(stashFile, []byte(content.String()))
}

func listStashes() ([]StashEntry, error) {
	stashDir := ".miniGit/stash"
	if !utils.FileExists(stashDir) {
		return []StashEntry{}, nil
	}

	files, err := os.ReadDir(stashDir)
	if err != nil {
		return nil, err
	}

	var stashes []StashEntry
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "stash_") {
			stash, err := loadStash(filepath.Join(stashDir, file.Name()))
			if err == nil {
				stashes = append(stashes, stash)
			}
		}
	}

	return stashes, nil
}

func loadStash(filePath string) (StashEntry, error) {
	var stash StashEntry
	stash.Index = make(map[string]string)
	stash.WorkDir = make(map[string]string)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return stash, err
	}

	lines := strings.Split(string(content), "\n")
	section := ""

	for _, line := range lines {
		originalLine := line
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "message: ") {
			stash.Message = strings.TrimPrefix(line, "message: ")
		} else if strings.HasPrefix(line, "branch: ") {
			stash.Branch = strings.TrimPrefix(line, "branch: ")
		} else if line == "index:" {
			section = "index"
		} else if line == "workdir:" {
			section = "workdir"
		} else if strings.HasPrefix(originalLine, "  ") {
			parts := strings.Fields(strings.TrimSpace(originalLine))
			if len(parts) == 2 {
				if section == "index" {
					stash.Index[parts[0]] = parts[1]
				} else if section == "workdir" {
					stash.WorkDir[parts[0]] = parts[1]
				}
			}
		}
	}

	fileName := filepath.Base(filePath)
	if strings.HasPrefix(fileName, "stash_") {
		idStr := strings.TrimPrefix(fileName, "stash_")
		if id, err := strconv.Atoi(idStr); err == nil {
			stash.ID = id
		}
	}

	return stash, nil
}

func restoreStash(stash StashEntry) error {
	err := utils.WriteIndex(stash.Index)
	if err != nil {
		return fmt.Errorf("failed to restore index: %w", err)
	}

	for filePath, hash := range stash.WorkDir {
		content := utils.GetBlobContent(hash)
		err := utils.WriteFile(filePath, []byte(content))
		if err != nil {
			return fmt.Errorf("failed to restore file %s: %w", filePath, err)
		}
	}

	return nil
}

func cleanWorkingDirectory() error {
	currentBranch := utils.GetCurrentBranchName()
	if currentBranch == "" {
		return fmt.Errorf("failed to get current branch")
	}

	commitHash, err := utils.GetCommitHashFromRef("refs/heads/" + currentBranch)
	if err != nil {
		return err
	}

	index, err := utils.ReadIndex()
	if err != nil {
		return err
	}

	tree := utils.ReadTreeFromCommit(commitHash)

	for filePath := range index {
		if _, exists := tree[filePath]; !exists {
			err := os.Remove(filePath)
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove file %s: %w", filePath, err)
			}
		}
	}

	for filePath, blobHash := range tree {
		content := utils.GetBlobContent(blobHash)
		err := utils.WriteFile(filePath, []byte(content))
		if err != nil {
			return err
		}
	}

	return utils.WriteFile(".miniGit/index", []byte(""))
}

func removeStash(stashID int) error {
	stashFile := filepath.Join(".miniGit/stash", fmt.Sprintf("stash_%d", stashID))
	return os.Remove(stashFile)
}
