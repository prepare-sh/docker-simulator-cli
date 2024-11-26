// cmd/stop.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [OPTIONS] CONTAINER",
	Short: "Stop one or more running containers",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, identifier := range args {
			if ContainerMgr.UpdateContainerStatus(identifier, "stopped") {
				fmt.Printf("Stopped container '%s'\n", identifier)
			} else {
				fmt.Printf("No such container: '%s'\n", identifier)
			}
		}
	},
}
