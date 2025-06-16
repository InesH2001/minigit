package cmd

import (
	"fmt"
	"minigit/core"
)

func InitCommand() {
	err := core.InitRepo()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("miniGit correctly initialized")
		fmt.Println("You can now start using miniGit for version control.")
	}
}
