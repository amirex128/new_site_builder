package sform

import (
	"context"
	"fmt"
)

// =============================================================================
// ServiceConnector
// =============================================================================

// ServiceConnectorOption is a function that configures the connection parameters
type ServiceConnectorOption func(*ConnectorParams)

// ConnectorParams holds all the connection parameters
type ConnectorParams struct {
	Name      string
	Host      string
	Username  string
	Password  string
	Database  string
	Port      int
	Driver    DBDriver
	Charset   string
	ParseTime bool
	Loc       string
	SSLMode   string
	TimeZone  string
	MaxOpen   int
	MaxIdle   int
	MaxLife   int
}

// ServiceConnector implements the serviceregistry.ServiceConnector interface
type ServiceConnector struct {
	options []RegistryOption
}

// RegisterConnection implements the serviceregistry.ServiceConnector interface
func (c *ServiceConnector) RegisterConnection(opts ...interface{}) error {
	// Initialize connection parameters
	params := &ConnectorParams{
		Driver:    MySQL,
		Charset:   "utf8mb4",
		ParseTime: true,
		Loc:       "Local",
		SSLMode:   "disable",
		TimeZone:  "UTC",
	}

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
	if params.Name == "" || params.Host == "" || params.Username == "" || params.Password == "" || params.Database == "" || params.Port == 0 {
		return fmt.Errorf("missing required connection parameters")
	}

	// Create the appropriate connection configuration based on driver
	var config ConnectionConfig

	switch params.Driver {
	case MySQL:
		mysqlConfig := &MySQLConfig{
			Username:     params.Username,
			Password:     params.Password,
			Host:         params.Host,
			Port:         params.Port,
			Database:     params.Database,
			Charset:      params.Charset,
			ParseTime:    params.ParseTime,
			Loc:          params.Loc,
			MaxOpenConns: params.MaxOpen,
			MaxIdleConns: params.MaxIdle,
		}
		config = mysqlConfig
	case PostgreSQL:
		pgConfig := &PostgresConfig{
			Username:     params.Username,
			Password:     params.Password,
			Host:         params.Host,
			Port:         params.Port,
			Database:     params.Database,
			SSLMode:      params.SSLMode,
			TimeZone:     params.TimeZone,
			MaxOpenConns: params.MaxOpen,
			MaxIdleConns: params.MaxIdle,
		}
		config = pgConfig
	case SQLServer:
		sqlServerConfig := &SQLServerConfig{
			Username:     params.Username,
			Password:     params.Password,
			Host:         params.Host,
			Port:         params.Port,
			Database:     params.Database,
			MaxOpenConns: params.MaxOpen,
			MaxIdleConns: params.MaxIdle,
		}
		config = sqlServerConfig
	case SQLite:
		sqliteConfig := &SQLiteConfig{
			Database:     params.Database,
			MaxOpenConns: params.MaxOpen,
			MaxIdleConns: params.MaxIdle,
		}
		config = sqliteConfig
	default:
		return fmt.Errorf("unsupported database driver: %s", params.Driver)
	}

	// Create options for the connection
	options := []RegistryOption{
		WithConnectionDetails(params.Name, config),
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

// Health checks if all database connections are alive
func (c *ServiceConnector) Health(ctx context.Context) error {
	return Health(ctx)
}

// =============================================================================
// Factory Functions
// =============================================================================

// GetConnector returns a new ServiceConnector instance
func GetConnector() *ServiceConnector {
	return &ServiceConnector{
		options: []RegistryOption{},
	}
}

// Option functions for service connector
func WithConnectionName(name string) ServiceConnectorOption {
	return func(params *ConnectorParams) {
		params.Name = name
	}
}

func WithHost(host string) ServiceConnectorOption {
	return func(params *ConnectorParams) {
		params.Host = host
	}
}

func WithPort(port int) ServiceConnectorOption {
	return func(params *ConnectorParams) {
		params.Port = port
	}
}

func WithUser(user string) ServiceConnectorOption {
	return func(params *ConnectorParams) {
		params.Username = user
	}
}

func WithPassword(password string) ServiceConnectorOption {
	return func(params *ConnectorParams) {
		params.Password = password
	}
}

func WithDatabase(database string) ServiceConnectorOption {
	return func(params *ConnectorParams) {
		params.Database = database
	}
}

func WithDriver(driver DBDriver) ServiceConnectorOption {
	return func(params *ConnectorParams) {
		params.Driver = driver
	}
}

func WithCharset(charset string) ServiceConnectorOption {
	return func(params *ConnectorParams) {
		params.Charset = charset
	}
}

func WithParseTime(parseTime bool) ServiceConnectorOption {
	return func(params *ConnectorParams) {
		params.ParseTime = parseTime
	}
}

func WithLocale(loc string) ServiceConnectorOption {
	return func(params *ConnectorParams) {
		params.Loc = loc
	}
}

func WithSSLMode(sslMode string) ServiceConnectorOption {
	return func(params *ConnectorParams) {
		params.SSLMode = sslMode
	}
}

func WithTimeZone(timeZone string) ServiceConnectorOption {
	return func(params *ConnectorParams) {
		params.TimeZone = timeZone
	}
}

func WithMaxOpenConns(maxOpen int) ServiceConnectorOption {
	return func(params *ConnectorParams) {
		params.MaxOpen = maxOpen
	}
}

func WithMaxIdleConns(maxIdle int) ServiceConnectorOption {
	return func(params *ConnectorParams) {
		params.MaxIdle = maxIdle
	}
}

func WithMaxLifetime(maxLife int) ServiceConnectorOption {
	return func(params *ConnectorParams) {
		params.MaxLife = maxLife
	}
}
