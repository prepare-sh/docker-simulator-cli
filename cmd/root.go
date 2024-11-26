// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"prepare.sh/dockermock/data"

	"github.com/spf13/cobra"
)

var (
	ContainerMgr *data.ContainerManager
	ImageMgr     *data.ImageManager
)

var rootCmd = &cobra.Command{
	Use:   "docker",
	Short: "A mock Docker CLI for demonstration purposes",
	Long:  `This is a Prepare.sh Docker (mock) CLI application that simulates Docker for lab environments.`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Initialize managers
	ContainerMgr = data.NewContainerManager()
	ImageMgr = data.NewImageManager()

	// Add subcommands
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(execCmd)
	rootCmd.AddCommand(psCmd)
	rootCmd.AddCommand(imagesCmd)
	rootCmd.AddCommand(pruneCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(runCmd)
}
