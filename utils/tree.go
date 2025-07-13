package utils

import (
	"os"
	"strings"
)

func GetTreeHashFromCommit(commitID string) string {
	data, _ := os.ReadFile(".miniGit/objects/commits/" + commitID)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "tree: ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "tree: "))
		}
	}
	return ""
}

func ReadTreeFromCommit(commitID string) map[string]string {
	treeHash := GetTreeHashFromCommit(commitID)

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

func AreTreesEqual(commitA, commitB string) bool {
	treeHashA := GetTreeHashFromCommit(commitA)
	treeHashB := GetTreeHashFromCommit(commitB)

	return treeHashA != "" && treeHashA == treeHashB
}

func WriteTree(treeHash, content string) error {
	return WriteFile(".miniGit/objects/trees/"+treeHash, []byte(content))
}

func GetBlobContent(blobHash string) string {
	data, _ := os.ReadFile(".miniGit/objects/blobs/" + blobHash)
	return string(data)
}