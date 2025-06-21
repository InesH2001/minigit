package main

import (
	"fmt"
	"minigit/cmd"
	"os"
	"minigit/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command. See 'minigit help'")
		return
	}

	command := os.Args[1]
	
	blockedDuringMerge := map[string]bool{
		"merge":    true,
		"checkout": true,
		"reset":    true,
		"rebase":   true,
		"revert":   true,
	}

	if utils.MergeInProgress() && blockedDuringMerge[command] {
		if !(command == "merge" && len(os.Args) >= 3 && os.Args[2] == "--abort") {
			fmt.Printf("Cannot run '%s': a merge is in progress. Use 'minigit merge --abort' or 'commit' to finish.\n", command)
			return
		}
	}

	switch command {
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
		if len(os.Args) == 3 && os.Args[2] == "--abort" {
			cmd.MergeAbort()
		} else if len(os.Args) < 3 {
			fmt.Println("Usage: minigit merge <branch_name> or minigit merge --abort")
		} else {
			cmd.Merge(os.Args[2])
		}

	case "revert":
		if len(os.Args) < 3 {
			fmt.Println("Usage: minigit revert <commit_hash>")
			return
		}
		cmd.Revert(os.Args[2:])
	
	case "rm":
		if len(os.Args) < 1 {
			fmt.Println("Usage: minigit rm [--cached] [-f] <file>")
			return
		}
		cmd.RmCommand(os.Args[2:])

	default:
		fmt.Println("Unknown command. See 'minigit help'")
	}
}
