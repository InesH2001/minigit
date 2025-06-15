package core

import (
	"fmt"
	"os"
	"minigit/utils"
)

func RemoveFile(file string, force bool, cachedOnly bool) error {

	index, err := utils.ReadIndex()
	if err != nil {
		return err
	}

	_, isStaged := index[file]

	if isStaged {
		if cachedOnly || force {
			delete(index, file)
		} else {
			return fmt.Errorf("error: '%s' has changes staged. Use --cached or -f to force removal", file)
		}
	}

	if !cachedOnly {
		headCommit, err := utils.GetHeadCommit()
		if err != nil {
			return err
		}

		tree := utils.ReadTreeFromCommit(headCommit)
		if _, inTree := tree[file]; inTree {
			if err := os.Remove(file); err != nil {
				return fmt.Errorf("failed to remove '%s' from disk: %w", file, err)
			}
			index[file] = ""
		} else if !isStaged {
			return fmt.Errorf("fatal: pathspec '%s' did not match any files", file)
		}
	}
	if err := utils.WriteIndex(index); err != nil {
		return fmt.Errorf("failed to update index: %w", err)
	}

	fmt.Printf("File '%s' removed\n", file)
	return nil
}
