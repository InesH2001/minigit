package core

import (
    "fmt"
    "os"
    "strings"
    "minigit/utils"
)

func Merge(branchName string) error {
    hasConflict := false
    modifiedFiles := []string{}

    currentBranch := getCurrentBranchRefPath()

    headCommitHash := getCommitHashFromRef(currentBranch)
    branchCommitHash := getCommitHashFromRef("refs/heads/" + branchName)
    commonAncestorHash := findCommonCommitAncestorHash(headCommitHash, branchCommitHash)

    headBlobs := readTreeFromCommit(headCommitHash)
    branchBlobs := readTreeFromCommit(branchCommitHash)
    baseBlobs := readTreeFromCommit(commonAncestorHash)

    allFiles := getUniqueUnionKeys(headBlobs, branchBlobs, baseBlobs)

    for _, file := range allFiles {
        base := getBlobContent(baseBlobs[file])
        head := getBlobContent(headBlobs[file])
        branch := getBlobContent(branchBlobs[file])

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

        parents := getParentsFromCommit(current)
        queue = append(queue, parents...)
    }

    queue = []string{commit2}
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]

        if visited[current] {
            return current
        }

        parents := getParentsFromCommit(current)
        queue = append(queue, parents...)
    }

    return ""
}


