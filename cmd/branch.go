package cmd

import (
	"fmt"
	"minigit/core"
)

func Branch(name string) {
	if name == "" {
		err := core.DisplayBranches()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
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
