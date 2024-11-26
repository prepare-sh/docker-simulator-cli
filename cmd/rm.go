// cmd/rm.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm [OPTIONS] CONTAINER",
	Short: "Remove one or more containers",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, identifier := range args {
			if ContainerMgr.RemoveContainer(identifier) {
				fmt.Printf("Removed container '%s'\n", identifier)
			} else {
				fmt.Printf("No such container: '%s'\n", identifier)
			}
		}
	},
}
