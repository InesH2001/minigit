package cmd

import (
	"fmt"
	"minigit/core"
)

func Branch(name string) {
	if name == "" {
		branches, err := core.ListBranches()
		if err != nil {
			fmt.Printf("Error listing branches: %v\n", err)
			return
		}

		currentBranch, err := core.GetCurrentBranchName()
		if err != nil {
			fmt.Printf("Error getting current branch: %v\n", err)
			return
		}

		for _, branch := range branches {
			if branch == currentBranch {
				fmt.Printf("* %s\n", branch)
			} else {
				fmt.Printf("  %s\n", branch)
			}
		}
	} else {
		err := core.CreateBranch(name)
		if err != nil {
			fmt.Printf("Error creating branch: %v\n", err)
			return
		}
		fmt.Printf("Branch '%s' created\n", name)
	}
}
