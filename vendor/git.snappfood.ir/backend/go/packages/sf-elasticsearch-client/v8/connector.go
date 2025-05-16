package v8

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
)

// =============================================================================
// ServiceConnector
// =============================================================================

// ServiceConnectorOption is a function that configures the connection parameters
type ServiceConnectorOption func(*connectionParams)

// connectionParams holds all the connection parameters
type connectionParams struct {
	name        string
	host        string
	username    string
	password    string
	port        int
	insecureTLS bool
	logger      Logger
}

// ServiceConnector implements the serviceregistry.ServiceConnector interface
type ServiceConnector struct {
	options []interface{} // Can hold both ServiceConnectorOption and RegistryOption
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

	// Set logger if provided
	if params.logger != nil {
		c.logger = params.logger
	}

	// Validate required fields
	if params.name == "" || params.host == "" || params.port == 0 {
		if c.logger != nil {
			// Convert category types to strings for logging
			extraMap := make(map[string]interface{})
			extraMap["name"] = params.name
			extraMap["host"] = params.host
			extraMap["port"] = params.port

			c.logger.ErrorWithCategory(
				Category.System.General,
				SubCategory.Status.Error,
				"Missing required connection parameters",
				extraMap,
			)
		}
		return fmt.Errorf("missing required connection parameters")
	}

	// Build Elasticsearch config
	config := elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("http://%s:%d", params.host, params.port)},
	}

	// Add authentication if provided
	if params.username != "" && params.password != "" {
		config.Username = params.username
		config.Password = params.password
	}

	// Add TLS configuration if needed
	if params.insecureTLS {
		// Create custom transport with insecure TLS
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		config.Transport = transport
	}

	// Create options for the connection
	options := []RegistryOption{}

	// Add logger if available
	if c.logger != nil {
		options = append(options, WithLogger(c.logger))

		// Convert category types to strings for logging
		extraMap := make(map[string]interface{})
		extraMap["connection"] = params.name
		extraMap["host"] = params.host
		extraMap["port"] = params.port

		c.logger.InfoWithCategory(
			Category.System.General,
			SubCategory.Networking.Connection,
			"Configuring Elasticsearch connection",
			extraMap,
		)
	}

	// Add connection details
	options = append(options, WithConnectionDetails(params.name, config))

	// Add any options set via WithOptions
	for _, opt := range c.options {
		if registryOpt, ok := opt.(RegistryOption); ok {
			options = append(options, registryOpt)
		}
	}

	// Register the connection using the registry pattern
	err := RegisterConnection(options...)
	if err != nil && c.logger != nil {
		// Convert category types to strings for logging
		extraMap := make(map[string]interface{})
		extraMap["connection"] = params.name
		extraMap["error"] = err.Error()

		c.logger.ErrorWithCategory(
			Category.System.General,
			SubCategory.Networking.Connection,
			"Failed to register Elasticsearch connection",
			extraMap,
		)
	}
	return err
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

// Health checks if all Elasticsearch connections are alive
func (c *ServiceConnector) Health(ctx context.Context) error {
	if c.logger != nil {
		c.logger.InfoWithCategory(
			Category.System.Health,
			SubCategory.Status.Info,
			"Checking Elasticsearch connections health",
			nil,
		)
	}
	return Health(ctx)
}

// =============================================================================
// Factory Functions
// =============================================================================

// GetConnector returns a new ServiceConnector instance
func GetConnector() *ServiceConnector {
	return &ServiceConnector{
		options: []interface{}{},
	}
}

// =============================================================================
// Helper Functions
// =============================================================================

// ParsePort converts a string port to an integer
func ParsePort(portStr string) (int, error) {
	return strconv.Atoi(portStr)
}

// ParseInsecureTLS converts a string to a boolean for insecure TLS setting
func ParseInsecureTLS(insecureStr string) (bool, error) {
	return strconv.ParseBool(insecureStr)
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

func WithPort(port int) ServiceConnectorOption {
	return func(p *connectionParams) {
		p.port = port
	}
}

func WithUser(username string) ServiceConnectorOption {
	return func(p *connectionParams) {
		p.username = username
	}
}

func WithPassword(password string) ServiceConnectorOption {
	return func(p *connectionParams) {
		p.password = password
	}
}

func WithInsecureTLS(insecure bool) ServiceConnectorOption {
	return func(p *connectionParams) {
		p.insecureTLS = insecure
	}
}

func WithConnectorLogger(logger Logger) ServiceConnectorOption {
	return func(p *connectionParams) {
		p.logger = logger
	}
}
