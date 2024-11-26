// cmd/ps.go
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List containers",
	Run: func(cmd *cobra.Command, args []string) {
		containers := ContainerMgr.ListContainers()
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "CONTAINER ID\tNAME\tIMAGE\tSTATUS")
		for _, c := range containers {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", c.ID, c.Name, c.Image, c.Status)
		}
		w.Flush()
	},
}
