package servicediscovery

import (
	"fmt"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"sync"
)

// DiscoveryType defines the type of service discovery
type DiscoveryType string

const (
	// Service discovery types
	ConsulDiscovery    DiscoveryType = "consul"
	ZookeeperDiscovery DiscoveryType = "zookeeper"
	EtcdDiscovery      DiscoveryType = "etcd"
)

// Config holds common configuration for service registry
type Config struct {
	// Type of service discovery (e.g., "consul", "zookeeper", "etcd")
	Type DiscoveryType

	// Address of the registry service (e.g., "localhost:8500" for Consul)
	Address string

	// Registry options for authentication, etc.
	Token    string
	Username string
	Password string

	// Additional options specific to the implementation
	Options map[string]interface{}
}

// ServiceInfo represents information about a registered service
type ServiceInfo struct {
	Name          string            // Service name
	ID            string            // Service ID (unique identifier)
	Address       string            // Service address (IP or hostname)
	Port          int               // Service port
	Tags          []string          // Service tags/metadata
	Meta          map[string]string // Additional metadata
	CheckURL      string            // Health check URL (if any)
	CheckTTL      string            // Health check TTL (e.g., "10s")
	CheckInterval string            // Health check interval (e.g., "5s")
}

// ServiceRegistry defines the interface for service registry operations
type ServiceRegistry interface {
	// Register registers a service with the registry
	Register(service ServiceInfo) error

	// Deregister removes a service from the registry
	Deregister(serviceID string) error

	// GetService returns the service information for a given service name
	GetService(serviceName string) ([]ServiceInfo, error)

	// GetServiceURL returns the URL for a service by name (e.g., "http://service-name:port")
	GetServiceURL(serviceName string) (string, error)

	// Close cleans up resources used by the registry
	Close() error
}

// Provider manages the current active service registry
type Provider struct {
	registry ServiceRegistry
	mu       sync.RWMutex
}

// Global provider instance
var provider = &Provider{}

// GetRegistry returns the current service registry
func GetRegistry() ServiceRegistry {
	provider.mu.RLock()
	defer provider.mu.RUnlock()

	if provider.registry == nil {
		return nil
	}

	return provider.registry
}

// InitRegistry initializes a new service registry based on the specified type
func InitRegistry(registryType DiscoveryType, config Config, log sflogger.Logger) error {
	var registry ServiceRegistry
	var err error

	switch registryType {
	case ConsulDiscovery:
		registry, err = NewConsulRegistry(config, log)
	case ZookeeperDiscovery:
		registry, err = NewZookeeperRegistry(config, log)
	case EtcdDiscovery:
		registry, err = NewEtcdRegistry(config, log)
	default:
		return fmt.Errorf("unsupported registry type: %s", registryType)
	}

	if err != nil {
		return err
	}

	provider.mu.Lock()
	provider.registry = registry
	provider.mu.Unlock()

	return nil
}
