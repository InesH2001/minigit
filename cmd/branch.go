package cmd

import (
	"fmt"
	"minigit/core"
	"minigit/utils"
)

func Branch(args []string) {
	if len(args) == 0 {
		err := core.DisplayBranches()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}

	var action, branchName, newName string
	force := false

	for i, arg := range args {
		switch arg {
		case "-d", "--delete":
			action = "delete"
			if i+1 < len(args) {
				branchName = args[i+1]
			}
		case "-D":
			action = "delete"
			force = true
			if i+1 < len(args) {
				branchName = args[i+1]
			}
		case "-m", "--move":
			action = "rename"
			if i+1 < len(args) {
				branchName = args[i+1]
			}
			if i+2 < len(args) {
				newName = args[i+2]
			}
		default:
			if action == "" && branchName == "" {
				action = "create"
				branchName = arg
			}
		}
	}

	switch action {
	case "create":
		if branchName == "" {
			fmt.Println("Usage: minigit branch <branch_name>")
			return
		}
		err := core.CreateBranch(branchName)
		if err != nil {
			fmt.Printf("Error creating branch: %v\n", err)
			return
		}
		fmt.Printf("Branch '%s' created\n", branchName)

	case "delete":
		if branchName == "" {
			fmt.Println("Usage: minigit branch -d <branch_name> or minigit branch -D <branch_name>")
			return
		}
		err := utils.DeleteBranch(branchName, force)
		if err != nil {
			fmt.Printf("Error deleting branch: %v\n", err)
			return
		}
		fmt.Printf("Branch '%s' deleted\n", branchName)

	case "rename":
		if branchName == "" || newName == "" {
			fmt.Println("Usage: minigit branch -m <old_name> <new_name>")
			return
		}
		err := utils.RenameBranch(branchName, newName)
		if err != nil {
			fmt.Printf("Error renaming branch: %v\n", err)
			return
		}
		fmt.Printf("Branch '%s' renamed to '%s'\n", branchName, newName)

	default:
		fmt.Println("Usage:")
		fmt.Println("  minigit branch                    # List branches")
		fmt.Println("  minigit branch <name>             # Create branch")
		fmt.Println("  minigit branch -d <name>          # Delete branch (safe)")
		fmt.Println("  minigit branch -D <name>          # Delete branch (force)")
		fmt.Println("  minigit branch -m <old> <new>     # Rename branch")
	}
}
