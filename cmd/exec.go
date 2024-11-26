// cmd/exec.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec [OPTIONS] CONTAINER COMMAND [ARG...]",
	Short: "Execute a command in a running container",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		containerID := args[0]
		command := args[1:]
		container, exists := ContainerMgr.GetContainer(containerID)
		if !exists {
			fmt.Printf("No such container: '%s'\n", containerID)
			return
		}
		if container.Status != "running" {
			fmt.Printf("Cannot exec in container '%s' as it is %s\n", containerID, container.Status)
			return
		}
		fmt.Printf("Executing command '%v' in container '%s'\n", command, containerID)
	},
}
