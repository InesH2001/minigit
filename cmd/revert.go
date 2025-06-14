package cmd

import (
    "fmt"
    "minigit/core"
)

func Revert(args []string) {
    if len(args) < 1 {
        fmt.Println("Usage: minigit revert <commit_hash>")
        return
    }

    err := core.Revert(args[0], "John Doe")
    if err != nil {
        fmt.Println("Revert failed:", err)
        return
    }

    fmt.Println("Revert completed successfully.")
}
