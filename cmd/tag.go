// cmd/tag.go
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag SOURCE_IMAGE[:TAG] TARGET_IMAGE[:TAG]",
	Short: "Create a tag TARGET_IMAGE that refers to SOURCE_IMAGE",
	Long: `Create a tag TARGET_IMAGE that refers to SOURCE_IMAGE.
If no tag is specified for SOURCE_IMAGE, 'latest' is assumed.
If no tag is specified for TARGET_IMAGE, 'latest' is assumed.`,
	Example: `  docker tag myimage:1.0 myrepo/myimage:2.0
  docker tag myimage myrepo/myimage:latest
  docker tag myimage:1.0 myrepo/myimage`,
	Args: cobra.ExactArgs(2),
	Run:  runTag,
}

func init() {
	rootCmd.AddCommand(tagCmd)
}

func runTag(cmd *cobra.Command, args []string) {
	sourceRef := args[0]
	targetRef := args[1]

	// Parse source image reference
	sourceName, sourceTag := parseImageReference(sourceRef)
	if sourceTag == "" {
		sourceTag = "latest"
	}

	// Parse target image reference
	targetName, targetTag := parseImageReference(targetRef)
	if targetTag == "" {
		targetTag = "latest"
	}

	// Create the new tag
	if success := ImageMgr.TagImage(sourceName, sourceTag, targetName, targetTag); success {
		fmt.Printf("Successfully tagged %s:%s as %s:%s\n", sourceName, sourceTag, targetName, targetTag)
	} else {
		fmt.Printf("Error: No such image: %s:%s\n", sourceName, sourceTag)
	}
}

func parseImageReference(ref string) (name, tag string) {
	parts := strings.Split(ref, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return ref, ""
}
