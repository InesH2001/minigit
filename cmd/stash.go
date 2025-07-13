package cmd

import (
	"fmt"
	"minigit/core"
)

func StashCommand(args []string) {
	if len(args) == 0 {
		err := core.Stash("")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}

	switch args[0] {
	case "pop":
		err := core.StashPop()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

	case "list":
		err := core.StashList()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

	default:
		message := args[0]
		err := core.Stash(message)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
