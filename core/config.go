package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func StoreUsername(usernameArg ...string) error {
	var username string

	if len(usernameArg) > 0 && strings.TrimSpace(usernameArg[0]) != "" {
		username = usernameArg[0]
	} else {
		fmt.Print("Enter your name (username): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		username = strings.TrimSpace(input)
	}

	return os.WriteFile(".miniGit/config", []byte("username="+username+"\n"), 0644)
}

func GetUsername() string {
	data, err := os.ReadFile(".miniGit/config")
	if err != nil {
		return "Unknown"
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "username=") {
			return strings.TrimPrefix(line, "username=")
		}
	}
	return "Unknown"
}
