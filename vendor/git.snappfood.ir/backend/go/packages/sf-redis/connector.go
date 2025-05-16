package sfredis

import (
	"context"
	"fmt"
	"strings"
)

// =============================================================================
// ServiceConnector
// =============================================================================

// ServiceConnectorOption is a function that configures the connection parameters
type ServiceConnectorOption func(*connectionParams)

// connectionParams holds all the connection parameters
type connectionParams struct {
	name      string
	host      string
	password  string
	database  int
	port      int
	connType  string
	addresses []string
}

// ServiceConnector implements the serviceregistry.ServiceConnector interface
type ServiceConnector struct {
	options []RegistryOption
	logger  Logger
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
	if params.name == "" || params.password == "" || params.connType == "" {
		if c.logger != nil {
			c.logger.ErrorWithCategory(Category.Database.Database, Category.Error.Error, "Missing required connection parameters", map[string]interface{}{
				"name":      params.name,
				"conn_type": params.connType,
			})
		}
		return fmt.Errorf("missing required connection parameters")
	}

	// Create options based on connection type
	var options []RegistryOption
	if params.connType == "cluster" {
		if len(params.addresses) == 0 {
			if c.logger != nil {
				c.logger.ErrorWithCategory(Category.Database.Database, Category.Error.Error, "Cluster addresses required for cluster connection", map[string]interface{}{
					"connection_name": params.name,
				})
			}
			return fmt.Errorf("cluster addresses required for cluster connection")
		}
		// For cluster, join all addresses with commas
		addr := strings.Join(params.addresses, ",")
		options = []RegistryOption{
			WithConnectionDetails(params.name, addr, params.password, params.database),
		}

		if c.logger != nil {
			c.logger.InfoWithCategory(Category.Database.Database, Category.System.Startup, "Registering cluster connection", map[string]interface{}{
				"connection_name": params.name,
				"addresses":       addr,
			})
		}
	} else {
		if params.host == "" {
			if c.logger != nil {
				c.logger.ErrorWithCategory(Category.Database.Database, Category.Error.Error, "Host required for normal connection", map[string]interface{}{
					"connection_name": params.name,
				})
			}
			return fmt.Errorf("host required for normal connection")
		}
		options = []RegistryOption{
			WithConnectionDetails(params.name, params.host, params.password, params.database),
		}

		if c.logger != nil {
			c.logger.InfoWithCategory(Category.Database.Database, Category.System.Startup, "Registering single node connection", map[string]interface{}{
				"connection_name": params.name,
				"host":            params.host,
				"database":        params.database,
			})
		}
	}

	// Add logger option if logger is set
	if c.logger != nil {
		options = append(options, WithLogger(c.logger))
	}

	// Add any options set via WithOptions
	for _, opt := range c.options {
		options = append(options, opt)
	}

	// Register the connection using the registry pattern
	return RegisterConnection(options...)
}

// WithOptions adds multiple custom options to the ServiceConnector
func (c *ServiceConnector) WithOptions(opts ...RegistryOption) *ServiceConnector {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		c.options = append(c.options, opt)
	}
	return c
}

// WithServiceLogger sets a logger for the ServiceConnector
func (c *ServiceConnector) WithServiceLogger(logger Logger) *ServiceConnector {
	c.logger = logger
	return c
}

// Health checks if all SfRedis connections are alive
func (c *ServiceConnector) Health(ctx context.Context) error {
	err := Health(ctx)
	if err != nil && c.logger != nil {
		c.logger.ErrorWithCategory(Category.Database.Database, Category.Error.Error, "Redis health check failed", map[string]interface{}{
			"error": err.Error(),
		})
	}
	return err
}

// =============================================================================
// Factory Functions
// =============================================================================

// GetConnector returns a new ServiceConnector instance
func GetConnector() *ServiceConnector {
	return &ServiceConnector{
		options: []RegistryOption{},
		logger:  nil,
	}
}

// Option functions for service connector
func WithConnectionName(name string) ServiceConnectorOption {
	return func(params *connectionParams) {
		params.name = name
	}
}

func WithHost(host string) ServiceConnectorOption {
	return func(params *connectionParams) {
		params.host = host
	}
}

func WithPassword(password string) ServiceConnectorOption {
	return func(params *connectionParams) {
		params.password = password
	}
}

func WithConnectionType(connType string) ServiceConnectorOption {
	return func(params *connectionParams) {
		params.connType = connType
	}
}

func WithRedisDatabase(database int) ServiceConnectorOption {
	return func(params *connectionParams) {
		params.database = database
	}
}

func WithClusterAddresses(addresses []string) ServiceConnectorOption {
	return func(params *connectionParams) {
		params.addresses = addresses
	}
}
