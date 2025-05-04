package utils

import (
    "io/ioutil"
    "os"
)

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
    lastSlash := -1
    for i := len(path) - 1; i >= 0; i-- {
        if path[i] == '/' {
            lastSlash = i
            break
        }
    }
    if lastSlash == -1 {
        return ""
    }
    return path[:lastSlash]
}
