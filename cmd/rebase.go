package cmd

import (
	"fmt"
	"minigit/core"
)

func Rebase(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: minigit rebase <branch_name> or minigit rebase --abort")
		return
	}

	if args[0] == "--abort" {
		err := core.RebaseAbort()
		if err != nil {
			fmt.Printf("Rebase abort failed: %v\n", err)
			return
		}
		return
	}

	branchName := args[0]
	err := core.Rebase(branchName)
	if err != nil {
		fmt.Printf("Rebase failed: %v\n", err)
		return
	}
}
