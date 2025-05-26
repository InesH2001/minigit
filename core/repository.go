package core

import (
    "os"
    "strings"
)

func getCurrentBranchRefPath() string {
    data, _ := os.ReadFile(".miniGit/HEAD")
    path := strings.TrimSpace(string(data))

    if strings.HasPrefix(path, "ref: ") {
        return strings.TrimPrefix(path, "ref: ")
    }

    return path
}

func getCommitHashFromRef(ref string) string {
    data, _ := os.ReadFile(".miniGit/" + ref)
    return strings.TrimSpace(string(data))
}

func readTreeFromCommit(commitID string) map[string]string {
    data, _ := os.ReadFile(".miniGit/objects/commits/" + commitID)
    lines := strings.Split(string(data), "\n")

    var treeHash string
    for _, line := range lines {
        if strings.HasPrefix(line, "tree: ") {
            treeHash = strings.TrimSpace(strings.TrimPrefix(line, "tree: "))
            break
        }
    }

    files := make(map[string]string)
    if treeHash != "" {
        treeData, err := os.ReadFile(".miniGit/objects/trees/" + treeHash)
        if err != nil {
            return files
        }

        treeLines := strings.Split(string(treeData), "\n")
        for _, line := range treeLines {
            if strings.Contains(line, " ") {
                parts := strings.SplitN(line, " ", 2)
                if len(parts) == 2 {
                    files[parts[0]] = parts[1]
                }
            }
        }
    }

    return files
}

func getBlobContent(blobHash string) string {
    data, _ := os.ReadFile(".miniGit/objects/blob/" + blobHash)
    return string(data)
}

func getParentsFromCommit(commitHash string) []string {
    data, _ := os.ReadFile("./miniGit/objects/commits/" + commitHash)
    lines := strings.Split(string(data), "\n")
    var parents []string

    for _, line := range lines {
        if strings.HasPrefix(line, "parent: ") {
            parentHash := strings.TrimSpace(strings.TrimPrefix(line, "parent: "))
            parents = append(parents, parentHash)
        }
    }

    return parents
}

func getUniqueUnionKeys(maps ...map[string]string) []string {
    keys := make(map[string]struct{})

    for _, m := range maps {
        for k := range m {
            keys[k] = struct{}{}
        }
    }

    var result []string
    for k := range keys {
        result = append(result, k)
    }

    return result
}
