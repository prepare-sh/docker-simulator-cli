// cmd/images.go
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "List images",
	Run: func(cmd *cobra.Command, args []string) {
		images := ImageMgr.ListImages()
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "IMAGE ID\tREPOSITORY\tTAG")
		for _, img := range images {
			fmt.Fprintf(w, "%s\t%s\t%s\n", img.ID, img.Name, img.Tag)
		}
		w.Flush()
	},
}
