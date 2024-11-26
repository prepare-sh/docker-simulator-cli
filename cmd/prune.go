// cmd/prune.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pruneCmd = &cobra.Command{
	Use:   "prune [OPTIONS]",
	Short: "Remove unused data",
	Run: func(cmd *cobra.Command, args []string) {
		// For demonstration, we'll clear all stopped containers and dangling images
		removedContainers := 0
		for _, c := range ContainerMgr.ListContainers() {
			if c.Status == "stopped" {
				if ContainerMgr.RemoveContainer(c.ID) {
					removedContainers++
				}
			}
		}

		removedImages := 0
		for _, img := range ImageMgr.ListImages() {
			// Assuming dangling images have no containers associated (for simplicity)
			// Here, we'll consider images not used by any container as dangling
			isUsed := false
			for _, c := range ContainerMgr.ListContainers() {
				if c.Image == fmt.Sprintf("%s:%s", img.Name, img.Tag) {
					isUsed = true
					break
				}
			}
			if !isUsed {
				if ImageMgr.RemoveImage(img.ID) {
					removedImages++
				}
			}
		}

		fmt.Printf("Pruned %d containers and %d images\n", removedContainers, removedImages)
	},
}
