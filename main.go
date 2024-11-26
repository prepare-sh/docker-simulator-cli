// main.go
package main

import (
	"fmt"
	"os"

	"prepare.sh/dockermock/cmd"
	"prepare.sh/dockermock/data"
)

func main() {
	// Ensure storage directory exists
	if err := data.EnsureStorageDir(); err != nil {
		fmt.Println("Error creating storage directory:", err)
		os.Exit(1)
	}

	cmd.Execute()
}
