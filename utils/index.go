package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadIndex() map[string]string {
	index := make(map[string]string)
	file, err := os.Open(".minigit/index")
	if err != nil {
		return index
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) == 2 {
			index[parts[0]] = parts[1]
		}
	}
	return index
}

func WriteIndex(index map[string]string) error {
	f, err := os.Create(".minigit/index")
	if err != nil {
		return err
	}
	defer f.Close()

	for path, hash := range index {
		if _, err := fmt.Fprintf(f, "%s %s\n", path, hash); err != nil {
			return err
		}
	}
	return nil
}
