// cmd/start.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [OPTIONS] CONTAINER",
	Short: "Start one or more stopped containers",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, identifier := range args {
			if ContainerMgr.UpdateContainerStatus(identifier, "running") {
				fmt.Printf("Started container '%s'\n", identifier)
			} else {
				fmt.Printf("No such container: '%s'\n", identifier)
			}
		}
	},
}
