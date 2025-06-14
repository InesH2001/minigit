package cmd

import (
    "fmt"
    "errors"  
    "minigit/core"
)

func CommitCommand(args []string) {
	params, err := ParseCommitArgs(args)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Usage: ./minigit commit -m <message> [-a <author>]")
		return
	}

    if err := core.Commit(params); err != nil {
        fmt.Println("Error:", err)
    }
}

func ParseCommitArgs(args []string) (core.CommitParams, error) {
	params := core.CommitParams{
		Author: core.GetUsername(),
        Message: "",
	}

	for i := 0; i < len(args); i++ {
		if args[i] == "-m" || args[i] == "--message" {
			if i+1 < len(args) {
				params.Message = args[i+1]
				i++
			}
		}
	}

	if params.Message == "" {
		return params, errors.New("missing commit message")
	}

	return params, nil
}
