package utils

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ListFilesRecursive(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if strings.Contains(path, "/.miniGit") {
			return nil
		}

		if strings.HasPrefix(filepath.Base(path), ".") && filepath.Base(path) != ".gitignore" {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return FilterIgnoredFiles(files), err
}

func WriteFile(path string, content []byte) error {
	dir := getDir(path)
	if dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(path, content, 0644)
}

func getDir(path string) string {
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		return ""
	}
	return path[:lastSlash]
}

func ListFiles() []string {
	files, err := ListFilesRecursive(".")
	if err != nil {
		return []string{}
	}
	return files
}
