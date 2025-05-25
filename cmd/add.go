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
			if err := core.Add(f); err != nil {
				fmt.Printf("Erreur lors de l'ajout de %s : %v\n", f, err)
			}
		}
	} else {
		for _, f := range args {
			if err := core.Add(f); err != nil {
				fmt.Printf("Erreur lors de l'ajout de %s : %v\n", f, err)
			}
		}
	}
}
