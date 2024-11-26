# Mock Docker CLI

**Docker CLI Simulator** is a Go-based command-line application that simulates essential Docker commands for demonstration and educational purposes. It mimics Docker's functionality without performing real container operations, using persistent storage to maintain state across sessions.

## 📦 Features

- **Simulated Commands:** `pull`, `push`, `rm`, `start`, `stop`, `exec`, `ps`, `images`, `prune`, `login`, `run`, `build`
- **Persistent Storage:** Stores container and image data in `/tmp/mock-docker`
- **Easy to Use:** Familiar Docker-like CLI experience

## 🚀 Installation

```bash
git clone https://github.com/prepare-sh/docker-simulator-cli.git
cd dockermock
go build -o dockermock
./dockermock --help