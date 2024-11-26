// cmd/push.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push [OPTIONS] IMAGE",
	Short: "Push an image to a registry",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		image := args[0]
		name, tag := parseImage(image)
		if ImageMgr.PushImage(name, tag) {
			fmt.Printf("Successfully pushed image '%s:%s'\n", name, tag)
		} else {
			fmt.Printf("Image '%s:%s' not found locally\n", name, tag)
		}
	},
}
