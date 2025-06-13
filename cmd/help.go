package cmd

import "fmt"

func Help() {
	fmt.Println("MiniGit - A Minimal Git-Like Version Control System")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  minigit <command> [arguments]")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  init             Initialize a new minigit repository")
	fmt.Println("  add <file...>    Add files to the staging area (use '.' for all)")
	fmt.Println("  commit -m <msg>  Commit staged changes with a message")
	fmt.Println("  branch [name]    Create or list branches")
	fmt.Println("  checkout <name>  Switch to specified branch")
	fmt.Println("  diff             Show changes between working dir and index")
	fmt.Println("  status           Show current status of working directory")
	fmt.Println("  log              Show commit history")
	fmt.Println("  reset <hash>     Reset working directory to a specific commit")
	fmt.Println("  merge <branch>   Merge specified branch into current one")
	fmt.Println("  stash            Stash current changes")
	fmt.Println("  revet <hash>     Revert a specific commit (like undo)")
	fmt.Println("  abort            Abort current merge or rebase operation")
	fmt.Println("  rebase <branch>  Reapply commits from current branch onto another")
	fmt.Println()
	fmt.Println("Use 'minigit <command>' for more information on a command.")
}
