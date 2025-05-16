package servicediscovery

import (
	"context"
	"encoding/json"
	"fmt"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdRegistry implements ServiceRegistry using etcd
type EtcdRegistry struct {
	client *clientv3.Client
	config Config
	ctx    context.Context
	cancel context.CancelFunc
}

// NewEtcdRegistry creates a new etcd-based service registry
func NewEtcdRegistry(config Config, log sflogger.Logger) (ServiceRegistry, error) {
	endpoints := []string{config.Address}
	if config.Address == "" {
		endpoints = []string{"localhost:2379"} // Default etcd address
	}

	// Apply any additional endpoints from the config
	if additionalEndpoints, ok := config.Options["endpoints"].([]string); ok && len(additionalEndpoints) > 0 {
		endpoints = additionalEndpoints
	}

	// Set up etcd client options
	etcdConfig := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}

	// Apply credentials if provided
	if username, ok := config.Options["username"].(string); ok {
		etcdConfig.Username = username
	}
	if password, ok := config.Options["password"].(string); ok {
		etcdConfig.Password = password
	}

	// Create a new client
	client, err := clientv3.New(etcdConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &EtcdRegistry{
		client: client,
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// serviceKey returns the key for a service in etcd
func (e *EtcdRegistry) serviceKey(serviceName, serviceID string) string {
	return fmt.Sprintf("/services/%s/%s", serviceName, serviceID)
}

// Register registers a service with etcd
func (e *EtcdRegistry) Register(service ServiceInfo) error {
	key := e.serviceKey(service.Name, service.ID)

	// Serialize service info to JSON
	serviceData, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service data: %w", err)
	}

	// Determine lease TTL if health check is configured
	var leaseID clientv3.LeaseID
	var keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse

	if service.CheckTTL != "" {
		// Parse the TTL duration
		ttl, err := time.ParseDuration(service.CheckTTL)
		if err != nil {
			return fmt.Errorf("invalid TTL format: %w", err)
		}

		// Create a lease
		lease, err := e.client.Grant(e.ctx, int64(ttl.Seconds()))
		if err != nil {
			return fmt.Errorf("failed to create lease: %w", err)
		}

		leaseID = lease.ID

		// Set up keep-alive
		keepAliveCh, err = e.client.KeepAlive(e.ctx, leaseID)
		if err != nil {
			return fmt.Errorf("failed to keep lease alive: %w", err)
		}

		// Monitor the keep-alive channel in a goroutine
		go func() {
			for {
				select {
				case <-e.ctx.Done():
					return
				case ka, ok := <-keepAliveCh:
					if !ok {
						return
					}
					// Use ka for logging if needed
					_ = ka
				}
			}
		}()
	}

	// Put the service data with or without a lease
	var opts []clientv3.OpOption
	if leaseID != 0 {
		opts = append(opts, clientv3.WithLease(leaseID))
	}

	_, err = e.client.Put(e.ctx, key, string(serviceData), opts...)
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	return nil
}

// Deregister removes a service from etcd
func (e *EtcdRegistry) Deregister(serviceID string) error {
	// Since we need the service name to form the key, we have to list
	// all services and find the one with matching ID
	resp, err := e.client.Get(e.ctx, "/services/", clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("failed to list services: %w", err)
	}

	for _, kv := range resp.Kvs {
		var serviceInfo ServiceInfo
		if err := json.Unmarshal(kv.Value, &serviceInfo); err != nil {
			continue
		}

		if serviceInfo.ID == serviceID {
			_, err := e.client.Delete(e.ctx, string(kv.Key))
			if err != nil {
				return fmt.Errorf("failed to deregister service: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("service with ID %s not found", serviceID)
}

// GetService returns service instances by name
func (e *EtcdRegistry) GetService(serviceName string) ([]ServiceInfo, error) {
	key := fmt.Sprintf("/services/%s/", serviceName)
	resp, err := e.client.Get(e.ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get service %s: %w", serviceName, err)
	}

	var serviceInfos []ServiceInfo
	for _, kv := range resp.Kvs {
		var serviceInfo ServiceInfo
		if err := json.Unmarshal(kv.Value, &serviceInfo); err != nil {
			continue
		}
		serviceInfos = append(serviceInfos, serviceInfo)
	}

	return serviceInfos, nil
}

// GetServiceURL returns the URL for a service by name
func (e *EtcdRegistry) GetServiceURL(serviceName string) (string, error) {
	services, err := e.GetService(serviceName)
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
func (e *EtcdRegistry) Close() error {
	e.cancel() // Cancel the context to stop all background goroutines
	return e.client.Close()
}
