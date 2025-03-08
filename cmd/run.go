// cmd/run.go
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var (
	// Command flags
	detached      bool
	containerName string
	portMappings  []string
	envVars       []string
	namespace     = "docker" // Default namespace
)

var runCmd = &cobra.Command{
	Use:   "run [OPTIONS] IMAGE [COMMAND]",
	Short: "Run a command in a new container",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		image := args[0]
		name, tag := parseImage(image)
		imageFull := fmt.Sprintf("%s:%s", name, tag)

		// Check if image exists locally or in repository
		imageExists := false
		for _, img := range ImageMgr.ListImages() {
			if (img.Name == name && img.Tag == tag) ||
				(img.ID == name && (tag == "latest" || img.Tag == tag)) {
				imageExists = true
				imageFull = fmt.Sprintf("%s:%s", img.Name, img.Tag)
				break
			}
		}
		if !imageExists {
			fmt.Printf("Image '%s' not found. Please pull it first.\n", imageFull)
			return
		}

		// Ensure the namespace exists
		ensureNamespace(namespace)

		// Generate a container name if not provided, ensuring Kubernetes compatibility
		podName := containerName
		if podName == "" {
			podName = generateK8sCompatibleName()
		} else {
			// Make sure provided name is k8s compatible
			podName = makeK8sCompatible(podName)
		}

		// Prepare command argument
		command := []string{}
		if len(args) > 1 {
			command = args[1:]
		}

		// Run the container on Kubernetes
		err := runContainerOnK8s(podName, imageFull, command, portMappings, envVars)
		if err != nil {
			fmt.Printf("Failed to run container: %v\n", err)
			return
		}

		// Create a record in our local database
		container := ContainerMgr.CreateContainer(podName, imageFull)
		fmt.Printf("Created and started container '%s' (ID: %s) from image '%s'\n", container.Name, container.ID, container.Image)

		// If not detached, follow logs until Ctrl+C
		if !detached {
			// Set up signal handling for Ctrl+C
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				<-sigCh
				fmt.Println("\nStopping container...")
				deleteContainerFromK8s(podName)
				ContainerMgr.RemoveContainer(container.ID)
				cancel()
				os.Exit(0)
			}()

			// Follow logs
			fmt.Println("Attaching to container. Press Ctrl+C to stop.")
			followPodLogs(ctx, podName)
		}
	},
}

// Generate a Kubernetes-compatible container name
func generateK8sCompatibleName() string {
	adjectives := []string{"bold", "brave", "calm", "eager", "fierce", "gentle", "happy", "jolly", "kind", "lively"}
	nouns := []string{"ant", "bear", "cat", "dog", "eagle", "fox", "giraffe", "horse", "iguana", "jaguar"}

	rand.Seed(time.Now().UnixNano())
	adj := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]

	// Ensure name is k8s compatible
	return fmt.Sprintf("%s-%s-%d", strings.ToLower(adj), strings.ToLower(noun), rand.Intn(1000))
}

// Make a string Kubernetes-compatible (lowercase, alphanumeric with hyphens)
func makeK8sCompatible(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace underscores with hyphens
	name = strings.ReplaceAll(name, "_", "-")

	// Remove any character that's not alphanumeric or hyphen
	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	// Ensure it starts and ends with alphanumeric
	s := result.String()
	if len(s) == 0 {
		return "container-" + fmt.Sprintf("%d", rand.Intn(1000))
	}

	// If first character is not alphanumeric, prepend 'c'
	if !(s[0] >= 'a' && s[0] <= 'z') && !(s[0] >= '0' && s[0] <= '9') {
		s = "c" + s
	}

	// If last character is not alphanumeric, append random number
	lastIndex := len(s) - 1
	if !(s[lastIndex] >= 'a' && s[lastIndex] <= 'z') && !(s[lastIndex] >= '0' && s[lastIndex] <= '9') {
		s = s + fmt.Sprintf("%d", rand.Intn(10))
	}

	return s
}

// Make sure the namespace exists on the K8s cluster
func ensureNamespace(ns string) {
	cmd := exec.Command("kubectl", "get", "namespace", ns)
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Creating namespace '%s'\n", ns)
		createCmd := exec.Command("kubectl", "create", "namespace", ns)
		if err := createCmd.Run(); err != nil {
			fmt.Printf("Warning: Failed to create namespace: %v\n", err)
		}
	}
}

// Run a container on Kubernetes
func runContainerOnK8s(name, image string, command, ports, env []string) error {
	args := []string{"run", name, "--image", image, "-n", namespace}

	// Add port mappings
	for _, port := range ports {
		args = append(args, "--port", port)
	}

	// Add environment variables
	for _, envVar := range env {
		args = append(args, "--env", envVar)
	}

	// Add command if specified
	if len(command) > 0 {
		args = append(args, "--command", "--")
		args = append(args, command...)
	}

	// If detached, add --restart=Always
	if detached {
		args = append(args, "--restart=Always")
	} else {
		// For interactive sessions
		args = append(args, "--restart=Never")
	}

	fmt.Printf("Starting container with: kubectl %s\n", strings.Join(args, " "))
	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Follow pod logs in real-time
func followPodLogs(ctx context.Context, podName string) {
	// Wait a moment for the pod to start
	time.Sleep(2 * time.Second)

	cmd := exec.CommandContext(ctx, "kubectl", "logs", "-f", podName, "-n", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run() // Ignore errors as the context might be canceled
}

// Delete a container from K8s
func deleteContainerFromK8s(name string) {
	cmd := exec.Command("kubectl", "delete", "pod", name, "-n", namespace)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: Failed to delete pod: %v\n%s\n", err, stderr.String())
	} else {
		fmt.Printf("Container '%s' deleted\n", name)
	}
}

func init() {
	// Add flags
	runCmd.Flags().BoolVarP(&detached, "detach", "d", false, "Run container in background")
	runCmd.Flags().StringVar(&containerName, "name", "", "Assign a name to the container")
	runCmd.Flags().StringArrayVarP(&portMappings, "publish", "p", []string{}, "Publish a container's port(s) to the host")
	runCmd.Flags().StringArrayVarP(&envVars, "env", "e", []string{}, "Set environment variables")
}
