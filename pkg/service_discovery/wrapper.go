package servicediscovery

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
)

// InitServiceRegistry initializes the service registry
func InitServiceRegistry(log sflogger.Logger, serviceDiscoveryType string, config Config) error {
	// Convert string to DiscoveryType
	var discoveryType DiscoveryType
	switch serviceDiscoveryType {
	case "consul":
		discoveryType = ConsulDiscovery
	case "zookeeper":
		discoveryType = ZookeeperDiscovery
	case "etcd":
		discoveryType = EtcdDiscovery
	default:
		return fmt.Errorf("unsupported discovery type: %s", serviceDiscoveryType)
	}

	log.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Operation.Initialization, "Initializing service registry", map[string]interface{}{
		"type":    string(discoveryType),
		"address": config.Address,
	})

	// Initialize the registry
	err := InitRegistry(discoveryType, config, log)
	if err != nil {
		extraMap := map[string]interface{}{
			sflogger.ExtraKey.Error.ErrorMessage: err.Error(),
			"type":                               string(discoveryType),
		}
		log.ErrorWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Error, "Failed to initialize service registry", extraMap)
		return err
	}

	log.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Success, "Service registry initialized successfully", map[string]interface{}{
		"type": string(discoveryType),
	})

	return nil
}

// RegisterService registers a service with the service registry
func RegisterService(log sflogger.Logger, service ServiceInfo) error {
	registry := GetRegistry()
	if registry == nil {
		return fmt.Errorf("service registry not initialized")
	}

	log.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Operation.Registration, "Registering service", map[string]interface{}{
		"name":    service.Name,
		"id":      service.ID,
		"address": service.Address,
		"port":    service.Port,
	})

	err := registry.Register(service)
	if err != nil {
		extraMap := map[string]interface{}{
			sflogger.ExtraKey.Error.ErrorMessage: err.Error(),
			"name":                               service.Name,
			"id":                                 service.ID,
		}
		log.ErrorWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Error, "Failed to register service", extraMap)
		return err
	}

	log.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Success, "Service registered successfully", map[string]interface{}{
		"name": service.Name,
		"id":   service.ID,
	})

	return nil
}

// DeregisterService deregisters a service from the service registry
func DeregisterService(log sflogger.Logger, serviceID string) error {
	registry := GetRegistry()
	if registry == nil {
		return fmt.Errorf("service registry not initialized")
	}

	log.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Operation.Registration, "Deregistering service", map[string]interface{}{
		"id": serviceID,
	})

	err := registry.Deregister(serviceID)
	if err != nil {
		extraMap := map[string]interface{}{
			sflogger.ExtraKey.Error.ErrorMessage: err.Error(),
			"id":                                 serviceID,
		}
		log.ErrorWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Error, "Failed to deregister service", extraMap)
		return err
	}

	log.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Success, "Service deregistered successfully", map[string]interface{}{
		"id": serviceID,
	})

	return nil
}

// GetServiceURL gets the URL for a service
func GetServiceURL(log sflogger.Logger, serviceName string) (string, error) {
	registry := GetRegistry()
	if registry == nil {
		return "", fmt.Errorf("service registry not initialized")
	}

	log.DebugWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Operation.Discovery, "Getting service URL", map[string]interface{}{
		"name": serviceName,
	})

	url, err := registry.GetServiceURL(serviceName)
	if err != nil {
		extraMap := map[string]interface{}{
			sflogger.ExtraKey.Error.ErrorMessage: err.Error(),
			"name":                               serviceName,
		}
		log.ErrorWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Error, "Failed to get service URL", extraMap)
		return "", err
	}

	log.DebugWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Success, "Service URL retrieved successfully", map[string]interface{}{
		"name": serviceName,
		"url":  url,
	})

	return url, nil
}

// CloseServiceRegistry closes the service registry
func CloseServiceRegistry(log sflogger.Logger) error {
	registry := GetRegistry()
	if registry == nil {
		return nil
	}

	log.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Operation.Shutdown, "Closing service registry", nil)

	err := registry.Close()
	if err != nil {
		extraMap := map[string]interface{}{
			sflogger.ExtraKey.Error.ErrorMessage: err.Error(),
		}
		log.ErrorWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Error, "Failed to close service registry", extraMap)
		return err
	}

	log.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Success, "Service registry closed successfully", nil)

	return nil
}

// GenerateServiceID generates a unique service ID
func GenerateServiceID(serviceName string) string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return fmt.Sprintf("%s-%s-%d", serviceName, hostname, os.Getpid())
}

// ServiceRegistryConfig holds all configuration needed for service registry initialization
type ServiceRegistryConfig struct {
	// Registry configuration
	Registry RegistryConfig

	// Service configuration
	Service ServiceConfig

	// Behavior configuration
	SetupGracefulShutdown bool // Whether to set up graceful shutdown
}

// RegistryConfig holds configuration for the service registry provider
type RegistryConfig struct {
	Type     string // Type of service registry (consul, zookeeper, etcd)
	Address  string // Address of the service registry
	Token    string // Authentication token for the service registry
	Username string // Username for the service registry
	Password string // Password for the service registry
}

// ServiceConfig holds configuration for the service being registered
type ServiceConfig struct {
	Name        string      // Name of this service
	Port        string      // Port of this service
	Address     string      // Address of this service
	Tags        string      // Tags for this service (comma-separated)
	Scheme      string      // Service URL scheme (http/https)
	HealthCheck HealthCheck // Health check configuration
}

// HealthCheck holds configuration for service health checks
type HealthCheck struct {
	URL      string // Health check URL for this service
	TTL      string // Health check TTL for this service
	Interval string // Health check interval
}

// InitAndRegisterService initializes the service registry and registers a service with it
// This function handles all the configuration parsing and initialization in one place
func InitAndRegisterService(log sflogger.Logger, config ServiceRegistryConfig) (string, error) {
	// Skip if service registry type is empty
	if config.Registry.Type == "" {
		log.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Info, "Service registry disabled, skipping initialization", nil)
		return "", nil
	}

	// Validate registry address
	if config.Registry.Address == "" {
		// Set default addresses based on registry type
		switch config.Registry.Type {
		case "consul":
			config.Registry.Address = "localhost:8500"
		case "zookeeper":
			config.Registry.Address = "localhost:2181"
		case "etcd":
			config.Registry.Address = "localhost:2379"
		default:
			config.Registry.Address = "localhost:8500"
		}

		log.WarnWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Warning, "Registry address not provided, using default", map[string]interface{}{
			"type":    config.Registry.Type,
			"address": config.Registry.Address,
		})
	}

	// Parse port string to int
	port := 8080 // Default port
	if config.Service.Port != "" {
		portInt, err := strconv.Atoi(config.Service.Port)
		if err == nil && portInt > 0 {
			port = portInt
		} else {
			extraMap := map[string]interface{}{
				sflogger.ExtraKey.Error.ErrorMessage: err.Error(),
			}
			log.WarnWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Warning, "Invalid port number, using default port 8080", extraMap)
		}
	}

	// Use a default address if not provided
	serviceAddress := config.Service.Address
	if serviceAddress == "" {
		// Try to get the hostname
		hostname, err := os.Hostname()
		if err == nil {
			serviceAddress = hostname
		} else {
			serviceAddress = "localhost"
		}
		log.WarnWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Warning, "Service address not provided, using hostname", map[string]interface{}{
			"hostname": hostname,
		})
	}

	// Parse tags from comma-separated string
	tags := []string{}
	if config.Service.Tags != "" {
		tags = strings.Split(config.Service.Tags, ",")
	}

	// Initialize registry configuration
	regConfig := Config{
		Address: config.Registry.Address,
		Options: map[string]interface{}{},
	}

	// Only add non-empty options
	if config.Registry.Token != "" {
		regConfig.Options["token"] = config.Registry.Token
	}
	if config.Registry.Username != "" {
		regConfig.Options["username"] = config.Registry.Username
	}
	if config.Registry.Password != "" {
		regConfig.Options["password"] = config.Registry.Password
	}

	// Initialize registry
	err := InitServiceRegistry(log, config.Registry.Type, regConfig)
	if err != nil {
		return "", fmt.Errorf("failed to initialize service registry: %w", err)
	}

	// Generate service ID
	serviceName := config.Service.Name
	if serviceName == "" {
		serviceName = "service"
		log.WarnWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Warning, "Service name not provided, using default name", map[string]interface{}{
			"defaultName": serviceName,
		})
	}
	serviceID := GenerateServiceID(serviceName)

	// Set default check interval if not provided
	checkInterval := "10s"
	if config.Service.HealthCheck.Interval != "" {
		checkInterval = config.Service.HealthCheck.Interval
	}

	// Register service
	serviceInfo := ServiceInfo{
		Name:          serviceName,
		ID:            serviceID,
		Address:       serviceAddress,
		Port:          port,
		Tags:          tags,
		Meta:          map[string]string{},
		CheckURL:      config.Service.HealthCheck.URL,
		CheckTTL:      config.Service.HealthCheck.TTL,
		CheckInterval: checkInterval,
	}

	// Add scheme to metadata if provided
	if config.Service.Scheme != "" {
		serviceInfo.Meta["scheme"] = config.Service.Scheme
	}

	err = RegisterService(log, serviceInfo)
	if err != nil {
		return "", fmt.Errorf("failed to register service: %w", err)
	}

	// Set up graceful shutdown if requested
	if config.SetupGracefulShutdown {
		setupGracefulShutdown(log, serviceID)
	}

	return serviceID, nil
}

// setupGracefulShutdown sets up a signal handler_init to deregister the service on shutdown
func setupGracefulShutdown(log sflogger.Logger, serviceID string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Operation.Shutdown, "Shutting down service...", nil)

		// Deregister service
		_ = DeregisterService(log, serviceID)

		// Close service registry
		_ = CloseServiceRegistry(log)

		os.Exit(0)
	}()
}
