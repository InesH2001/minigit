package utils

import (
	"os"
	"strings"
)

func ReadTreeFromCommit(commitID string) map[string]string {
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

func WriteTree(treeHash, content string) error {
	return WriteFile(".miniGit/objects/trees/"+treeHash, []byte(content))
}

func GetBlobContent(blobHash string) string {
	data, _ := os.ReadFile(".miniGit/objects/blobs/" + blobHash)
	return string(data)
}