package utils

import (
	"os"
)

func MergeInProgress() bool {
	_, err := os.Stat(".miniGit/MERGE_HEAD")
	return err == nil
}