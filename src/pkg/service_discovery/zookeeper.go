package servicediscovery

import (
	"encoding/json"
	"fmt"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"path"
	"time"

	"github.com/go-zookeeper/zk"
)

// ZookeeperRegistry implements ServiceRegistry using Zookeeper
type ZookeeperRegistry struct {
	conn   *zk.Conn
	config Config
}

// NewZookeeperRegistry creates a new Zookeeper-based service registry
func NewZookeeperRegistry(config Config, log sflogger.Logger) (ServiceRegistry, error) {
	// Set default address if not provided
	servers := []string{config.Address}
	if config.Address == "" {
		servers = []string{"localhost:2181"} // Default Zookeeper address
	}

	// Apply any additional servers from the config
	if additionalServers, ok := config.Options["servers"].([]string); ok && len(additionalServers) > 0 {
		servers = additionalServers
	}

	// Connect to Zookeeper
	conn, _, err := zk.Connect(servers, time.Second*5)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Zookeeper: %w", err)
	}

	// Create base path
	baseServicePath := "/services"
	exists, _, err := conn.Exists(baseServicePath)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to check if base path exists: %w", err)
	}

	if !exists {
		_, err = conn.Create(baseServicePath, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			conn.Close()
			return nil, fmt.Errorf("failed to create base path: %w", err)
		}
	}

	return &ZookeeperRegistry{
		conn:   conn,
		config: config,
	}, nil
}

// serviceKey returns the key for a service in Zookeeper
func (z *ZookeeperRegistry) serviceKey(serviceName, serviceID string) string {
	return path.Join("/services", serviceName, serviceID)
}

// servicePath returns the path for a service type in Zookeeper
func (z *ZookeeperRegistry) servicePath(serviceName string) string {
	return path.Join("/services", serviceName)
}

// Register registers a service with Zookeeper
func (z *ZookeeperRegistry) Register(service ServiceInfo) error {
	// Ensure the service path exists
	servicePath := z.servicePath(service.Name)
	exists, _, err := z.conn.Exists(servicePath)
	if err != nil {
		return fmt.Errorf("failed to check if service path exists: %w", err)
	}

	if !exists {
		_, err = z.conn.Create(servicePath, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return fmt.Errorf("failed to create service path: %w", err)
		}
	}

	// Create service node
	serviceKey := z.serviceKey(service.Name, service.ID)

	// Serialize service info to JSON
	serviceData, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service data: %w", err)
	}

	// Determine flags and TTL
	flags := int32(0)

	if service.CheckTTL != "" {
		flags = zk.FlagEphemeral

		// Parse TTL to validate it, even though we don't use the value directly
		_, err = time.ParseDuration(service.CheckTTL)
		if err != nil {
			return fmt.Errorf("invalid TTL format: %w", err)
		}
	}

	// Check if node already exists
	exists, stat, err := z.conn.Exists(serviceKey)
	if err != nil {
		return fmt.Errorf("failed to check if service exists: %w", err)
	}

	if exists {
		// Update existing node
		_, err = z.conn.Set(serviceKey, serviceData, stat.Version)
		if err != nil {
			return fmt.Errorf("failed to update service: %w", err)
		}
	} else {
		// Create new node
		_, err = z.conn.Create(serviceKey, serviceData, flags, zk.WorldACL(zk.PermAll))
		if err != nil {
			return fmt.Errorf("failed to register service: %w", err)
		}
	}

	return nil
}

// Deregister removes a service from Zookeeper
func (z *ZookeeperRegistry) Deregister(serviceID string) error {
	// We need to find the service path since we need the service name
	// First, get all service types
	services, _, err := z.conn.Children("/services")
	if err != nil {
		return fmt.Errorf("failed to list services: %w", err)
	}

	// Check each service type for the service ID
	for _, serviceName := range services {
		servicePath := z.servicePath(serviceName)
		serviceInstances, _, err := z.conn.Children(servicePath)
		if err != nil {
			continue
		}

		for _, instance := range serviceInstances {
			if instance == serviceID {
				serviceKey := z.serviceKey(serviceName, serviceID)
				err := z.conn.Delete(serviceKey, -1)
				if err != nil && err != zk.ErrNoNode {
					return fmt.Errorf("failed to deregister service: %w", err)
				}
				return nil
			}
		}
	}

	return fmt.Errorf("service with ID %s not found", serviceID)
}

// GetService returns service instances by name
func (z *ZookeeperRegistry) GetService(serviceName string) ([]ServiceInfo, error) {
	servicePath := z.servicePath(serviceName)

	// Check if the service exists
	exists, _, err := z.conn.Exists(servicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to check if service exists: %w", err)
	}

	if !exists {
		return []ServiceInfo{}, nil
	}

	// Get all service instances
	instances, _, err := z.conn.Children(servicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get service instances for %s: %w", serviceName, err)
	}

	var serviceInfos []ServiceInfo
	for _, instance := range instances {
		serviceKey := z.serviceKey(serviceName, instance)
		data, _, err := z.conn.Get(serviceKey)
		if err != nil {
			continue
		}

		var serviceInfo ServiceInfo
		if err := json.Unmarshal(data, &serviceInfo); err != nil {
			continue
		}
		serviceInfos = append(serviceInfos, serviceInfo)
	}

	return serviceInfos, nil
}

// GetServiceURL returns the URL for a service by name
func (z *ZookeeperRegistry) GetServiceURL(serviceName string) (string, error) {
	services, err := z.GetService(serviceName)
	if err != nil {
		return "", err
	}

	if len(services) == 0 {
		return "", fmt.Errorf("no healthy instances found for service: %s", serviceName)
	}

	// Use the first service instance
	service := services[0]

	// Construct service URL
	scheme := "http"
	if value, ok := service.Meta["scheme"]; ok {
		scheme = value
	}

	return fmt.Sprintf("%s://%s:%d", scheme, service.Address, service.Port), nil
}

// Close releases resources used by the registry
func (z *ZookeeperRegistry) Close() error {
	z.conn.Close()
	return nil
}
