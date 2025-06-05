package servicediscovery

import (
	"fmt"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"

	"github.com/hashicorp/consul/api"
)

// ConsulRegistry implements ServiceRegistry using Consul
type ConsulRegistry struct {
	client *api.Client
	config Config
	log    sflogger.Logger
}

// NewConsulRegistry creates a new Consul-based service registry
func NewConsulRegistry(config Config, log sflogger.Logger) (ServiceRegistry, error) {
	consulConfig := api.DefaultConfig()

	// Set address if provided
	if config.Address != "" {
		consulConfig.Address = config.Address
	}

	// Set token if provided
	if config.Token != "" {
		consulConfig.Token = config.Token
	}

	// Apply additional options from the config
	if timeout, ok := config.Options["timeout"].(string); ok && timeout != "" {
		duration, err := time.ParseDuration(timeout)
		if err == nil && duration > 0 {
			consulConfig.WaitTime = duration
		}
	}

	// Set TLS configuration if provided
	if tlsEnabled, ok := config.Options["tls_enabled"].(bool); ok && tlsEnabled {
		tlsConfig := api.TLSConfig{}

		if caFile, ok := config.Options["ca_file"].(string); ok && caFile != "" {
			tlsConfig.CAFile = caFile
		}
		if certFile, ok := config.Options["cert_file"].(string); ok && certFile != "" {
			tlsConfig.CertFile = certFile
		}
		if keyFile, ok := config.Options["key_file"].(string); ok && keyFile != "" {
			tlsConfig.KeyFile = keyFile
		}
		if insecureSkipVerify, ok := config.Options["insecure_skip_verify"].(bool); ok {
			tlsConfig.InsecureSkipVerify = insecureSkipVerify
		}

		consulConfig.TLSConfig = tlsConfig
	}

	// Create a new client
	client, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %w", err)
	}

	// Verify connection by pinging the agent
	_, err = client.Agent().Self()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Consul agent: %w", err)
	}

	log.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Success, "Connected to Consul service registry", map[string]interface{}{
		"address": config.Address,
	})

	return &ConsulRegistry{
		client: client,
		config: config,
		log:    log,
	}, nil
}

// Register registers a service with Consul
func (c *ConsulRegistry) Register(service ServiceInfo) error {
	if service.Name == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	if service.ID == "" {
		return fmt.Errorf("service ID cannot be empty")
	}

	c.log.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Operation.Registration, "Registering service with Consul", map[string]interface{}{
		"name":    service.Name,
		"id":      service.ID,
		"address": service.Address,
		"port":    service.Port,
	})

	registration := &api.AgentServiceRegistration{
		ID:      service.ID,
		Name:    service.Name,
		Address: service.Address,
		Port:    service.Port,
		Tags:    service.Tags,
		Meta:    service.Meta,
	}

	// Add health check if configured
	if service.CheckURL != "" {
		check := &api.AgentServiceCheck{
			HTTP:          service.CheckURL,
			Interval:      service.CheckInterval,
			TLSSkipVerify: true,
		}

		// Set timeout to 5s by default, or half the interval if specified
		timeout := "5s"
		if service.CheckInterval != "" {
			// Try to parse interval and set timeout to half of it
			if interval, err := time.ParseDuration(service.CheckInterval); err == nil {
				halfInterval := interval / 2
				if halfInterval > time.Second {
					timeout = halfInterval.String()
				}
			}
		}
		check.Timeout = timeout

		registration.Check = check
	} else if service.CheckTTL != "" {
		check := &api.AgentServiceCheck{
			TTL: service.CheckTTL,
		}
		registration.Check = check
	}

	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		extraMap := map[string]interface{}{
			sflogger.ExtraKey.Error.ErrorMessage: err.Error(),
			"name":                               service.Name,
			"id":                                 service.ID,
		}
		c.log.ErrorWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Error, "Failed to register service with Consul", extraMap)
		return fmt.Errorf("failed to register service with Consul: %w", err)
	}

	c.log.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Success, "Service registered successfully with Consul", map[string]interface{}{
		"name": service.Name,
		"id":   service.ID,
	})

	return nil
}

// Deregister removes a service from Consul
func (c *ConsulRegistry) Deregister(serviceID string) error {
	if serviceID == "" {
		return fmt.Errorf("service ID cannot be empty")
	}

	c.log.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Operation.Registration, "Deregistering service from Consul", map[string]interface{}{
		"id": serviceID,
	})

	err := c.client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		c.log.ErrorWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Error, "Failed to deregister service from Consul", map[string]interface{}{
			"id":                                 serviceID,
			sflogger.ExtraKey.Error.ErrorMessage: err.Error(),
		})
		return fmt.Errorf("failed to deregister service from Consul: %w", err)
	}

	c.log.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Success, "Service deregistered successfully from Consul", map[string]interface{}{
		"id": serviceID,
	})

	return nil
}

// GetService returns service instances by name
func (c *ConsulRegistry) GetService(serviceName string) ([]ServiceInfo, error) {
	if serviceName == "" {
		return nil, fmt.Errorf("service name cannot be empty")
	}

	c.log.DebugWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Operation.Discovery, "Getting service instances from Consul", map[string]interface{}{
		"name": serviceName,
	})

	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		c.log.ErrorWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Error, "Failed to get service from Consul", map[string]interface{}{
			"name":                               serviceName,
			sflogger.ExtraKey.Error.ErrorMessage: err.Error(),
		})
		return nil, fmt.Errorf("failed to get service %s: %w", serviceName, err)
	}

	var serviceInfos []ServiceInfo
	for _, service := range services {
		serviceInfos = append(serviceInfos, ServiceInfo{
			Name:    service.Service.Service,
			ID:      service.Service.ID,
			Address: service.Service.Address,
			Port:    service.Service.Port,
			Tags:    service.Service.Tags,
			Meta:    service.Service.Meta,
		})
	}

	c.log.DebugWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Success, "Found service instances from Consul", map[string]interface{}{
		"name":  serviceName,
		"count": len(serviceInfos),
	})

	return serviceInfos, nil
}

// GetServiceURL returns the URL for a service by name
func (c *ConsulRegistry) GetServiceURL(serviceName string) (string, error) {
	if serviceName == "" {
		return "", fmt.Errorf("service name cannot be empty")
	}

	services, err := c.GetService(serviceName)
	if err != nil {
		return "", err
	}

	if len(services) == 0 {
		c.log.WarnWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Warning, "No healthy instances found for service", map[string]interface{}{
			"name": serviceName,
		})
		return "", fmt.Errorf("no healthy instances found for service: %s", serviceName)
	}

	// Use the first service instance
	service := services[0]

	// Construct service URL
	scheme := "http"
	if value, ok := service.Meta["scheme"]; ok && value != "" {
		scheme = value
	}

	serviceURL := fmt.Sprintf("%s://%s:%d", scheme, service.Address, service.Port)
	c.log.DebugWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Success, "Service URL constructed", map[string]interface{}{
		"name": serviceName,
		"url":  serviceURL,
	})

	return serviceURL, nil
}

// Close closes the Consul registry client
func (c *ConsulRegistry) Close() error {
	c.log.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Operation.Shutdown, "Closing Consul service registry", nil)
	return nil
}
