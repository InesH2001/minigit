package cmd

import (
	"fmt"
	"minigit/core"
)

func RunLog() {
	err := core.Log()
	if err != nil {
		fmt.Println("Error running log:", err)
	}
}
