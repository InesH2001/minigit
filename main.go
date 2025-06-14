package main

import (
	"fmt"
	"minigit/cmd"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command. See 'minigit help'")
		return
	}

	switch os.Args[1] {
	case "init":
		cmd.InitCommand()
	
	case "set-user":
    	cmd.SetUserCommand(os.Args[2:])	
	
	case "get-user":
    	cmd.GetUserCommand()

	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: minigit add <file1> <file2> ... or '.' for all")
			return
		}
		cmd.Add(os.Args[2:])

	case "diff":
		cmd.DiffCommand()

	case "commit":
		cmd.CommitCommand(os.Args[2:])

	case "status":
		cmd.RunStatus()

	case "log":
		cmd.RunLog()

	case "reset":
		cmd.RunReset()

	case "help":
		cmd.Help()

	case "branch":
		branchName := ""
		if len(os.Args) >= 3 {
			branchName = os.Args[2]
		}
		cmd.Branch(branchName)

	case "checkout":
		if len(os.Args) < 3 {
			fmt.Println("Usage: minigit checkout <branchname>")
			return
		}
		cmd.Checkout(os.Args[2])

	case "merge":
		if len(os.Args) < 3 {
			fmt.Println("Usage: minigit merge <branch_name>")
			return
		}
		cmd.Merge(os.Args[2])

	case "revert":
		if len(os.Args) < 3 {
			fmt.Println("Usage: minigit revert <commit_hash>")
			return
		}
		cmd.Revert(os.Args[2:])

	default:
		fmt.Println("Unknown command. See 'minigit help'")
	}
}
