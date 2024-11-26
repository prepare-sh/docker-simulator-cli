// cmd/run.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [OPTIONS] IMAGE [COMMAND]",
	Short: "Run a command in a new container",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		image := args[0]
		name, tag := parseImage(image)
		imageFull := fmt.Sprintf("%s:%s", name, tag)

		// Check if image exists
		imageExists := false
		for _, img := range ImageMgr.ListImages() {
			if img.Name == name && img.Tag == tag {
				imageExists = true
				break
			}
		}
		if !imageExists {
			fmt.Printf("Image '%s' not found. Please pull it first.\n", imageFull)
			return
		}

		// Create container
		containerName := fmt.Sprintf("container_%d", ContainerMgr.Counter)
		container := ContainerMgr.CreateContainer(containerName, imageFull)
		fmt.Printf("Created and started container '%s' (ID: %s) from image '%s'\n", container.Name, container.ID, container.Image)

		// If a command is provided, simulate exec
		if len(args) > 1 {
			command := args[1:]
			fmt.Printf("Executing command '%v' in container '%s'\n", command, container.ID)
		}
	},
}
