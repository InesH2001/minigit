package utils

import (
    "crypto/sha1"
    "encoding/hex"
)

func HashContent(content []byte) string {
    hash := sha1.Sum(content)
    return hex.EncodeToString(hash[:])
}