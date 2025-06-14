package cmd

import (
	"fmt"
	"minigit/core"
)

func RunStatus() {
	err := core.Status()
	if err != nil {
		fmt.Println("Error running status:", err)
	}
}
