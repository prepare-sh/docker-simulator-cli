// data/containers.go
package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

// Container represents a mock container
type Container struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"` // e.g., running, stopped
}

// ContainerManager manages mock containers
type ContainerManager struct {
	containers map[string]*Container
	Counter    int // Renamed from 'counter' to 'Counter' to export it
	mu         sync.Mutex
}

// NewContainerManager initializes a ContainerManager with persisted data
func NewContainerManager() *ContainerManager {
	cm := &ContainerManager{
		containers: make(map[string]*Container),
		Counter:    1, // Initialize Counter
	}
	if err := cm.Load(); err != nil {
		fmt.Println("Warning: Unable to load containers data:", err)
	}
	return cm
}

// Load reads containers data from the JSON file
func (cm *ContainerManager) Load() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	path := GetContainersFilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// No data to load
		return nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var containers []*Container
	if err := json.Unmarshal(data, &containers); err != nil {
		return err
	}

	for _, c := range containers {
		cm.containers[c.ID] = c
		// Update counter based on existing IDs
		var idNum int
		_, err := fmt.Sscanf(c.ID, "c%d", &idNum)
		if err == nil && idNum >= cm.Counter {
			cm.Counter = idNum + 1
		}
	}

	return nil
}

// Save writes containers data to the JSON file
func (cm *ContainerManager) Save() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	containers := []*Container{}
	for _, c := range cm.containers {
		containers = append(containers, c)
	}

	data, err := json.MarshalIndent(containers, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(GetContainersFilePath(), data, 0644)
}

// CreateContainer simulates creating a container
func (cm *ContainerManager) CreateContainer(name, image string) *Container {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	id := fmt.Sprintf("c%03d", cm.Counter)
	cm.Counter++ // Increment Counter
	container := &Container{
		ID:     id,
		Name:   name,
		Image:  image,
		Status: "running",
	}
	cm.containers[id] = container
	cm.Save()
	return container
}

// GetContainer retrieves a container by ID or name
func (cm *ContainerManager) GetContainer(identifier string) (*Container, bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for _, c := range cm.containers {
		if c.ID == identifier || c.Name == identifier {
			return c, true
		}
	}
	return nil, false
}

// ListContainers lists all containers
func (cm *ContainerManager) ListContainers() []*Container {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	list := []*Container{}
	for _, c := range cm.containers {
		list = append(list, c)
	}
	return list
}

// RemoveContainer removes a container
func (cm *ContainerManager) RemoveContainer(identifier string) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for id, c := range cm.containers {
		if c.ID == identifier || c.Name == identifier {
			delete(cm.containers, id)
			cm.Save()
			return true
		}
	}
	return false
}

// UpdateContainerStatus updates the status of a container
func (cm *ContainerManager) UpdateContainerStatus(identifier, status string) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if c, exists := cm.GetContainer(identifier); exists {
		c.Status = status
		cm.Save()
		return true
	}
	return false
}
