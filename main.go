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

	default:
		fmt.Println("Unknown command. See 'minigit help'")
	}
}
