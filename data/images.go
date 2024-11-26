// data/images.go
package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Image represents a mock Docker image
type Image struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

// ImageManager manages mock images
type ImageManager struct {
	images  map[string]*Image
	counter int
	mu      sync.Mutex
}

// NewImageManager initializes an ImageManager with persisted data
func NewImageManager() *ImageManager {
	im := &ImageManager{
		images:  make(map[string]*Image),
		counter: 1,
	}
	if err := im.Load(); err != nil {
		fmt.Println("Warning: Unable to load images data:", err)
	}
	return im
}

// Load reads images data from the JSON file
func (im *ImageManager) Load() error {
	im.mu.Lock()
	defer im.mu.Unlock()

	path := GetImagesFilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// No data to load
		return nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var images []*Image
	if err := json.Unmarshal(data, &images); err != nil {
		return err
	}

	for _, img := range images {
		im.images[img.ID] = img
		// Update counter based on existing IDs
		var idNum int
		_, err := fmt.Sscanf(img.ID, "i%d", &idNum)
		if err == nil && idNum >= im.counter {
			im.counter = idNum + 1
		}
	}

	return nil
}

// Save writes images data to the JSON file
func (im *ImageManager) Save() error {
	im.mu.Lock()
	defer im.mu.Unlock()

	// Create slice of images
	images := make([]*Image, 0, len(im.images))
	for _, img := range im.images {
		images = append(images, img)
	}

	// Marshal the data
	data, err := json.MarshalIndent(images, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal images: %v", err)
	}

	// Get the file path
	filePath := GetImagesFilePath()

	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write the file
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	fmt.Printf("Successfully saved %d images to %s\n", len(images), filePath)
	return nil
}

// PullImage simulates pulling an image
func (im *ImageManager) PullImage(name, tag string) *Image {
	im.mu.Lock()
	defer im.mu.Unlock()

	// Check if image already exists
	for _, img := range im.images {
		if img.Name == name && img.Tag == tag {
			fmt.Printf("Image %s:%s already exists\n", name, tag)
			return img
		}
	}

	// Simulate download progress
	fmt.Printf("Pulling image %s:%s\n", name, tag)
	for i := 0; i <= 100; i += 10 {
		fmt.Printf("Download progress: %d%%\n", i)
		time.Sleep(100 * time.Millisecond) // Simulate download time
	}

	id := fmt.Sprintf("i%03d", im.counter)
	im.counter++
	image := &Image{
		ID:   id,
		Name: name,
		Tag:  tag,
	}
	im.images[id] = image

	if err := im.Save(); err != nil {
		fmt.Printf("Warning: Failed to save image data: %v\n", err)
	}

	fmt.Printf("Successfully pulled %s:%s\n", name, tag)
	return image
}

// PushImage simulates pushing an image
func (im *ImageManager) PushImage(name, tag string) bool {
	im.mu.Lock()
	defer im.mu.Unlock()

	for _, img := range im.images {
		if img.Name == name && img.Tag == tag {
			return true
		}
	}
	return false
}

// ListImages lists all images
func (im *ImageManager) ListImages() []*Image {
	im.mu.Lock()
	defer im.mu.Unlock()

	list := []*Image{}
	for _, img := range im.images {
		list = append(list, img)
	}
	return list
}

// RemoveImage removes an image
func (im *ImageManager) RemoveImage(identifier string) bool {
	im.mu.Lock()
	defer im.mu.Unlock()

	for id, img := range im.images {
		if img.ID == identifier || img.Name == identifier {
			delete(im.images, id)
			im.Save()
			return true
		}
	}
	return false
}