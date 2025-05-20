package core

import (
	"fmt"
	"os"
	"path/filepath"

	"minigit/utils"
)

func Add(filePath string) error {
	if !utils.FileExists(filePath) {
		fmt.Printf("⚠️  Ignoré : %s (n'existe pas)\n", filePath)
		return nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	hash := utils.HashContent(content)
	index := utils.ReadIndex()

	if currentHash, ok := index[filePath]; ok && currentHash == hash {
		fmt.Printf("Inchangé : %s\n", filePath)
		return nil
	}

	blobPath := filepath.Join(".mgit", "objects", "blob", hash)
	if !utils.FileExists(blobPath) {
		err := os.WriteFile(blobPath, content, 0644)
		if err != nil {
			return err
		}
		fmt.Printf("🆕 Blob créé : %s\n", hash)
	}

	index[filePath] = hash
	err = utils.WriteIndex(index)
	if err != nil {
		return err
	}
	fmt.Printf("Ajouté : %s\n", filePath)
	return nil
}
