// data/storage.go
package data

import (
	"os"
	"path/filepath"
)

const (
	StorageDir     = "./mock-docker"
	ContainersFile = "containers.json"
	ImagesFile     = "images.json"
)

// EnsureStorageDir ensures that the storage directory exists
func EnsureStorageDir() error {
	return os.MkdirAll(StorageDir, os.ModePerm)
}

// GetContainersFilePath returns the full path to the containers file
func GetContainersFilePath() string {
	return filepath.Join(StorageDir, ContainersFile)
}

// GetImagesFilePath returns the full path to the images file
func GetImagesFilePath() string {
	return filepath.Join(StorageDir, ImagesFile)
}
