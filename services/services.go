package services

import (
	"fmt"
	"reflect"
	"sync"
)

// Container is a service container for dependency injection
type Container struct {
	services map[string]interface{}
	mu       sync.RWMutex
}

// NewContainer creates a new service container
func NewContainer() *Container {
	return &Container{
		services: make(map[string]interface{}),
	}
}

// Register registers a service with the container
// The service can be registered by name or by type
func (c *Container) Register(name string, service interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = service
}

// RegisterByType registers a service by its type name
func (c *Container) RegisterByType(service interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	serviceType := reflect.TypeOf(service)
	if serviceType.Kind() == reflect.Ptr {
		serviceType = serviceType.Elem()
	}
	
	c.services[serviceType.Name()] = service
}

// Get retrieves a service by name
func (c *Container) Get(name string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	service, exists := c.services[name]
	if !exists {
		return nil, fmt.Errorf("service '%s' not found", name)
	}
	
	return service, nil
}

// GetByType retrieves a service by type
func (c *Container) GetByType(serviceType interface{}) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	t := reflect.TypeOf(serviceType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	
	service, exists := c.services[t.Name()]
	if !exists {
		return nil, fmt.Errorf("service of type '%s' not found", t.Name())
	}
	
	return service, nil
}

// MustGet retrieves a service by name, panics if not found
func (c *Container) MustGet(name string) interface{} {
	service, err := c.Get(name)
	if err != nil {
		panic(err)
	}
	return service
}

// MustGetByType retrieves a service by type, panics if not found
func (c *Container) MustGetByType(serviceType interface{}) interface{} {
	service, err := c.GetByType(serviceType)
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
func (c *Container) All() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	result := make(map[string]interface{})
	for k, v := range c.services {
		result[k] = v
	}
	return result
}

// Clear removes all services
func (c *Container) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services = make(map[string]interface{})
}

