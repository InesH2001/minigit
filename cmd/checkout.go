package cmd

import (
	"fmt"
	"minigit/core"
)

func Checkout(branchName string) {
	err := core.SwitchToBranch(branchName)
	if err != nil {
		fmt.Printf("Error switching to branch: %v\n", err)
		return
	}
	fmt.Printf("Switched to branch '%s'\n", branchName)
}
