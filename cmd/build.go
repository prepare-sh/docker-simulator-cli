// cmd/build.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"prepare.sh/dockermock/data"
)

var (
	buildTag     string
	buildFile    string
	buildNoCache bool
)

var buildCmd = &cobra.Command{
	Use:   "build [OPTIONS] PATH",
	Short: "Build an image from a Dockerfile",
	Long: `Build an image from a Dockerfile and a 'context'.
The PATH specifies the location of the Dockerfile and build context.`,
	Example: `  docker build .
  docker build -t myapp:latest .
  docker build -f custom.Dockerfile .
  docker build --no-cache .`,
	Args: cobra.ExactArgs(1),
	Run:  runBuild,
}

func init() {
	buildCmd.Flags().StringVarP(&buildTag, "tag", "t", "", "Name and optionally a tag in the 'name:tag' format")
	buildCmd.Flags().StringVarP(&buildFile, "file", "f", "Dockerfile", "Name of the Dockerfile")
	buildCmd.Flags().BoolVar(&buildNoCache, "no-cache", false, "Do not use cache when building the image")
	rootCmd.AddCommand(buildCmd)
}

func runBuild(cmd *cobra.Command, args []string) {
	contextPath := args[0]
	dockerfilePath := filepath.Join(contextPath, buildFile)

	// Check if context path exists
	if _, err := os.Stat(contextPath); os.IsNotExist(err) {
		fmt.Printf("Error: context path does not exist: %s\n", contextPath)
		return
	}

	// Check if Dockerfile exists
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		fmt.Printf("Error: Dockerfile not found in %s\n", dockerfilePath)
		return
	}

	// Parse tag if provided
	var name, tag string
	if buildTag != "" {
		parts := strings.Split(buildTag, ":")
		name = parts[0]
		if len(parts) > 1 {
			tag = parts[1]
		} else {
			tag = "latest"
		}
	} else {
		// Generate default name from directory
		name = filepath.Base(contextPath)
		tag = "latest"
	}

	// Parse Dockerfile
	commands, err := data.ParseDockerfile(dockerfilePath)
	if err != nil {
		fmt.Printf("Error parsing Dockerfile: %v\n", err)
		return
	}

	// Simulate build process
	fmt.Printf("Building image from context: %s\n", contextPath)
	fmt.Printf("Using Dockerfile: %s\n", dockerfilePath)

	// Process Dockerfile commands
	if !processDockerfileCommands(commands) {
		fmt.Println("\nBuild failed")
		return
	}

	// Create the image
	image := ImageMgr.BuildImage(name, tag)
	if image != nil {
		fmt.Printf("\nSuccessfully built image '%s:%s' (ID: %s)\n", image.Name, image.Tag, image.ID)
	} else {
		fmt.Println("\nBuild failed")
	}
}

func processDockerfileCommands(commands []data.DockerfileCommand) bool {
	fmt.Println("\nSending build context to Docker daemon...")
	time.Sleep(500 * time.Millisecond)

	fmt.Println("Step 1/6 : Analyzing build context")
	time.Sleep(300 * time.Millisecond)

	for i, cmd := range commands {
		fmt.Printf("Step %d/%d : %s %s\n", i+1, len(commands), cmd.Instruction, cmd.Arguments)

		// Simulate command execution
		fmt.Printf(" ---> Running in temporary container\n")
		time.Sleep(500 * time.Millisecond)

		// Generate a pseudo-random layer ID
		layerID := fmt.Sprintf("%x", time.Now().UnixNano()%0xFFFFFF)
		fmt.Printf(" ---> %s\n", layerID)

		if i < len(commands)-1 {
			fmt.Println(" ---> Removing intermediate container")
		}

		// Simulate specific command behaviors
		switch cmd.Instruction {
		case "FROM":
			if !handleFromCommand(cmd.Arguments) {
				return false
			}
		case "COPY", "ADD":
			fmt.Printf(" ---> Copying files...\n")
		case "RUN":
			fmt.Printf(" ---> Running command: %s\n", cmd.Arguments)
		}
	}

	return true
}

func handleFromCommand(baseImage string) bool {
	// Check if base image exists (you could integrate this with your image management system)
	fmt.Printf(" ---> Using base image %s\n", baseImage)
	time.Sleep(300 * time.Millisecond)
	return true
}
