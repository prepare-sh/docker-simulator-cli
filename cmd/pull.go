// cmd/pull.go
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull [OPTIONS] IMAGE",
	Short: "Pull an image from a registry",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		image := args[0]
		name, tag := parseImage(image)
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
