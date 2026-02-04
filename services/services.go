package services

import (
	"fmt"
	"sync"
)

// Container is a service container for dependency injection
type Container struct {
	services map[string]any
	mu       sync.RWMutex
}

// NewContainer creates a new service container
func NewContainer() *Container {
	return &Container{
		services: make(map[string]any),
	}
}

// Register registers a service with the container
// The service can be registered by name or by type
func (c *Container) Register(name string, service any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = service
}

// Get retrieves a service by name
func (c *Container) Get(name string) (any, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	service, exists := c.services[name]
	if !exists {
		return nil, fmt.Errorf("service '%s' not found", name)
	}

	return service, nil
}

// MustGet retrieves a service by name, panics if not found
func (c *Container) MustGet(name string) any {
	service, err := c.Get(name)
	if err != nil {
		panic(err)
	}
	return service
}

// Has checks if a service is registered
func (c *Container) Has(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.services[name]
	return exists
}

// All returns all registered services
func (c *Container) All() map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]any)
	for k, v := range c.services {
		result[k] = v
	}
	return result
}

// Clear removes all services
func (c *Container) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services = make(map[string]any)
}
