// cmd/login.go
package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"prepare.sh/dockermock/data"
)

var loginCmd = &cobra.Command{
	Use:   "login [OPTIONS] [SERVER]",
	Short: "Log in to a Docker registry",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		server := "https://index.docker.io/v1/"
		if len(args) > 0 {
			server = args[0]
		}

		// For GitHub Container Registry (ghcr.io), use GitHub CLI
		if strings.Contains(server, "ghcr.io") {
			fmt.Println("Detected GitHub Container Registry. Using GitHub CLI for authentication...")

			// Check if gh CLI is available
			_, err := exec.LookPath("gh")
			if err != nil {
				fmt.Println("Error: GitHub CLI (gh) not found. Please install it to authenticate with ghcr.io")
				os.Exit(1)
			}

			// Check if already authenticated with gh
			authStatusCmd := exec.Command("gh", "auth", "status")
			output, err := authStatusCmd.CombinedOutput()
			if err != nil || !strings.Contains(string(output), "Logged in") {
				fmt.Println("Not logged in to GitHub. Please authenticate with GitHub CLI...")

				// Run gh auth login
				loginCmd := exec.Command("gh", "auth", "login")
				loginCmd.Stdin = os.Stdin
				loginCmd.Stdout = os.Stdout
				loginCmd.Stderr = os.Stderr

				if err := loginCmd.Run(); err != nil {
					fmt.Println("Failed to authenticate with GitHub:", err)
					os.Exit(1)
				}
			}

			// Get GitHub token
			tokenCmd := exec.Command("gh", "auth", "token")
			token, err := tokenCmd.Output()
			if err != nil {
				fmt.Println("Failed to get GitHub token:", err)
				os.Exit(1)
			}

			// Store the token in our mock docker config
			username := "USERNAME" // Placeholder - with GitHub token, actual username isn't needed
			storeToken := strings.TrimSpace(string(token))

			if err := storeCredentials(server, username, storeToken); err != nil {
				fmt.Println("Failed to store credentials:", err)
				os.Exit(1)
			}

			fmt.Printf("Login Succeeded to %s\n", server)
		} else {
			// Standard username/password login for other registries
			// This is a simplified version - in a real implementation you'd use proper credential management
			fmt.Println("Standard login flow not implemented for non-GitHub registries in this version.")
		}
	},
}

func storeCredentials(server, username, token string) error {
	configDir := filepath.Join(data.StorageDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Create a simplified docker config structure
	configContent := fmt.Sprintf(`{
  "auths": {
    "%s": {
      "auth": "%s"
    }
  }
}`, server, base64Encode(username+":"+token))

	return os.WriteFile(filepath.Join(configDir, "config.json"), []byte(configContent), 0600)
}

func base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
