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