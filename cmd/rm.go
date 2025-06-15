package cmd

import (
	"fmt"
	"minigit/core"
)

func RmCommand(args []string) {
	cachedOnly := false
	force := false
	files := []string{}

	for _, arg := range args {
		switch arg {
		case "--cached":
			cachedOnly = true
		case "-f":
			force = true
		default:
			files = append(files, arg)
		}
	}

	if len(files) == 0 {
		fmt.Println("Please specify a file to remove.")
		return
	}

	for _, file := range files {
		if err := core.RemoveFile(file, force, cachedOnly); err != nil {
			fmt.Println(err)
		}
	}
}
