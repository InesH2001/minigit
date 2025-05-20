package cmd

import (
	"fmt"

	"minigit/core"
	"minigit/utils"
)

func Add(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage : mgit add <fichier|.>")
		return
	}

	if args[0] == "." {
		files, err := utils.ListFilesRecursive(".")
		if err != nil {
			fmt.Println("Erreur lors de la lecture des fichiers :", err)
			return
		}
		for _, f := range files {
			core.Add(f)
		}
	} else {
		for _, f := range args {
			core.Add(f)
		}
	}
}
