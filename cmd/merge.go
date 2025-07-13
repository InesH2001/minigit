package cmd

import (
	"fmt"
	"minigit/core"
)

func Merge(branchName string) {
	err := core.Merge(branchName)
	if err != nil {
		fmt.Println("Merge failed:", err)
		return
	}
}

func MergeAbort() {
    err := core.MergeAbort()
    if err != nil {
        fmt.Println("Merge abort failed:", err)
    }
}