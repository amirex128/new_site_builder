package sfrabbitmq

import (
	"context"
	"fmt"
	"reflect"
)

// =============================================================================
// ServiceConnector
// =============================================================================

// ServiceConnectorOption is a function that configures the connection parameters
type ServiceConnectorOption func(*connectionParams)

// connectionParams holds all the connection parameters
type connectionParams struct {
	name     string
	host     string
	user     string
	password string
	vhost    string
	port     int
}

// ServiceConnector implements the serviceregistry.ServiceConnector interface
type ServiceConnector struct {
	options []interface{} // Can hold both ServiceConnectorOption and RegistryOption
}

// RegisterConnection implements the serviceregistry.ServiceConnector interface
func (c *ServiceConnector) RegisterConnection(opts ...interface{}) error {
	// Initialize connection parameters
	params := &connectionParams{}

	// Apply all options
	for _, opt := range opts {
		switch o := opt.(type) {
		case ServiceConnectorOption:
			o(params)
		case RegistryOption:
			c.options = append(c.options, o)
		}
	}

	// Validate required fields
	if params.name == "" || params.host == "" || params.user == "" || params.password == "" || params.vhost == "" || params.port == 0 {
		return fmt.Errorf("missing required connection parameters")
	}

	// Create options for the connection
	options := []RegistryOption{
		WithConnectionDetails(params.name, params.host, params.port, params.user, params.password),
	}

	// Add any options set via WithOptions
	for _, opt := range c.options {
		if registryOpt, ok := opt.(RegistryOption); ok {
			options = append(options, registryOpt)
		}
	}

	// Check if we have a logger option
	hasLogger := false
	for _, opt := range opts {
		// Check if this is a WithLogger option by comparing function pointers
		// This relies on the fact that WithLogger will appear in the function name
		if fn := reflect.ValueOf(opt).String(); fn != "" {
			if fnName := reflect.ValueOf(opt).Pointer(); fnName == reflect.ValueOf(WithLogger).Pointer() {
				hasLogger = true
				break
			}
		}
	}

	if !hasLogger {
		if globalRegistry.logger != nil {
			extras := make(map[string]interface{})
			extras["connection_name"] = params.name
			globalRegistry.logger.WarnWithCategory(Category.Infrastructure.Network, SubCategory.Operation.Initialization, "No logger provided when registering connection", extras)
		}
	}

	// Register the connection using the registry pattern
	return RegisterConnection(options...)
}

// WithOptions adds multiple custom options to the ServiceConnector
func (c *ServiceConnector) WithOptions(opts ...interface{}) *ServiceConnector {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		c.options = append(c.options, opt)
	}
	return c
}

// Health checks if all RabbitMQ connections are alive
func (c *ServiceConnector) Health(ctx context.Context) error {
	if globalRegistry.logger != nil {
		globalRegistry.logger.InfoWithCategory(Category.Infrastructure.Network, SubCategory.Operation.Initialization, "Checking RabbitMQ health", nil)
	}

	err := Health(ctx)

	if err != nil && globalRegistry.logger != nil {
		extras := make(map[string]interface{})
		extras[ExtraKey.Error.ErrorMessage] = err.Error()
		globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "RabbitMQ health check failed", extras)
	} else if globalRegistry.logger != nil {
		globalRegistry.logger.InfoWithCategory(Category.Infrastructure.Network, SubCategory.Status.Success, "RabbitMQ health check passed", nil)
	}

	return err
}

// GetConnector returns a new ServiceConnector instance
func GetConnector() *ServiceConnector {
	return &ServiceConnector{
		options: []interface{}{},
	}
}

// Option functions for service connector
func WithConnectionName(name string) ServiceConnectorOption {
	return func(p *connectionParams) {
		p.name = name
	}
}

func WithHost(host string) ServiceConnectorOption {
	return func(p *connectionParams) {
		p.host = host
	}
}

func WithUser(user string) ServiceConnectorOption {
	return func(p *connectionParams) {
		p.user = user
	}
}

func WithPassword(password string) ServiceConnectorOption {
	return func(p *connectionParams) {
		p.password = password
	}
}

func WithVHost(vhost string) ServiceConnectorOption {
	return func(p *connectionParams) {
		p.vhost = vhost
	}
}

func WithPort(port int) ServiceConnectorOption {
	return func(p *connectionParams) {
		p.port = port
	}
}

// StartOutboxProcessor starts the background worker to process failed messages
func (c *ServiceConnector) StartOutboxProcessor() {
	if globalRegistry.logger != nil {
		globalRegistry.logger.InfoWithCategory(Category.API.Messaging, SubCategory.Operation.Initialization, "Starting outbox processor", nil)
	}
	globalRegistry.StartOutboxProcessor()
}
