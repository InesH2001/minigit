package cmd

import (
    "fmt"
    "minigit/core"
)

func SetUserCommand(args []string) {
	err := core.StoreUsername(args...)
	if err != nil {
		fmt.Println("Error saving username:", err)
		return
	}

	fmt.Println("Username updated successfully.")
}

func GetUserCommand() {
	username := core.GetUsername()
	fmt.Println("Current user (username):", username)
}