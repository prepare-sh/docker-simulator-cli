// cmd/login.go
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
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
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Username for '%s': ", server)
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		fmt.Printf("Password: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		// For demonstration, we'll just print a success message
		fmt.Printf("Login Succeeded for user '%s' on '%s'\n", username, server)
	},
}
