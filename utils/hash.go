package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

func HashContent(data []byte) string {
	hash := sha1.Sum(data)
	return hex.EncodeToString(hash[:])
}
