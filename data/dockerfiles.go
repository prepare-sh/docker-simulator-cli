// data/dockerfile.go
package data

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// DockerfileCommand represents a command in a Dockerfile
type DockerfileCommand struct {
	Instruction string
	Arguments   string
}

// ParseDockerfile reads and parses a Dockerfile
func ParseDockerfile(path string) ([]DockerfileCommand, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open Dockerfile: %v", err)
	}
	defer file.Close()

	var commands []DockerfileCommand
	scanner := bufio.NewScanner(file)

	// For handling multi-line commands
	var currentCommand string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle line continuations
		if strings.HasSuffix(line, "\\") {
			currentCommand += line[:len(line)-1] + " "
			continue
		}

		// Complete the multi-line command or process single line
		if currentCommand != "" {
			line = currentCommand + line
			currentCommand = ""
		}

		// Parse the command
		parts := strings.SplitN(line, " ", 2)
		instruction := strings.ToUpper(parts[0])
		args := ""
		if len(parts) > 1 {
			args = parts[1]
		}

		commands = append(commands, DockerfileCommand{
			Instruction: instruction,
			Arguments:   args,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading Dockerfile: %v", err)
	}

	return commands, nil
}
