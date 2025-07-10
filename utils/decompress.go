package utils

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func ReadAndDecompressBlob(hash string) ([]byte, error) {
	blobPath := filepath.Join(".miniGit", "objects", "blobs", hash)
	compressedData, err := os.ReadFile(blobPath)
	if err != nil {
		return nil, err
	}
	return Decompress(compressedData)
}

func Decompress(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}
