// cmd/pull.go
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"prepare.sh/dockermock/data"
)

var pullCmd = &cobra.Command{
	Use:   "pull [OPTIONS] IMAGE",
	Short: "Pull an image from a registry",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		image := args[0]
		name, tag := parseImage(image)

		// Check if this is a ghcr.io image
		if strings.HasPrefix(name, "ghcr.io/") {
			// Verify authentication
			if !isAuthenticatedForRegistry("ghcr.io") {
				fmt.Println("Error: Not authenticated to ghcr.io. Please run 'docker login ghcr.io' first")
				os.Exit(1)
			}
			fmt.Printf("Using GitHub authentication for %s\n", name)
		}

		// Simulate network activity for pulling
		simulatePull(name, tag)

		// Store image in our local database
		img := ImageMgr.PullImage(name, tag)
		fmt.Printf("Successfully pulled image '%s:%s' (ID: %s)\n", img.Name, img.Tag, img.ID)
	},
}

func parseImage(image string) (name, tag string) {
	parts := strings.Split(image, ":")
	if len(parts) == 2 {
		name = parts[0]
		tag = parts[1]
	} else {
		name = image
		tag = "latest"
	}
	return
}

// Check if we have valid authentication for a registry
func isAuthenticatedForRegistry(registry string) bool {
	configPath := filepath.Join(data.StorageDir, "config", "config.json")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return false
	}

	// Read config file
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Warning: Failed to read config file: %v\n", err)
		return false
	}

	// Parse JSON
	var config struct {
		Auths map[string]interface{} `json:"auths"`
	}

	if err := json.Unmarshal(configData, &config); err != nil {
		fmt.Printf("Warning: Failed to parse config: %v\n", err)
		return false
	}

	// Check if we have auth for this registry
	for reg := range config.Auths {
		if strings.Contains(reg, registry) {
			return true
		}
	}

	return false
}

// Simulate a real Docker pull with progress indicators
func simulatePull(name, tag string) {
	fmt.Printf("Pulling from %s\n", name)

	// Simulate downloading different layers
	layers := []string{
		"Pulling fs layer",
		"Downloading",
		"Download complete",
		"Extracting",
		"Pull complete",
	}

	// For each simulated layer
	for i := 0; i < 3; i++ {
		layerId := fmt.Sprintf("%x", i*12345) // fake layer ID

		for _, action := range layers {
			fmt.Printf("%s: %s\n", layerId, action)
			time.Sleep(100 * time.Millisecond)
		}
	}

	fmt.Printf("Digest: sha256:%x\n", time.Now().UnixNano())
	fmt.Printf("Status: Downloaded newer image for %s:%s\n", name, tag)
}
