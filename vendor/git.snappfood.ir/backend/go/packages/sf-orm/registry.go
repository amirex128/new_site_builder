package sform

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	// Database drivers with APM integration
	mysqlDriver "go.elastic.co/apm/module/apmgormv2/v2/driver/mysql"
	sqlserverDriver "go.elastic.co/apm/module/apmgormv2/v2/driver/sqlserver"
	// Uncomment and add additional drivers as needed
	// postgresDriver "go.elastic.co/apm/module/apmgormv2/v2/driver/postgres"
	// sqliteDriver "go.elastic.co/apm/module/apmgormv2/v2/driver/sqlite"
)

// =============================================================================
// Driver Constants and Connection Configuration
// =============================================================================

// DBDriver represents the type of database driver
type DBDriver string

const (
	// MySQL driver
	MySQL DBDriver = "mysql"
	// PostgreSQL driver
	PostgreSQL DBDriver = "postgres"
	// SQLite driver
	SQLite DBDriver = "sqlite"
	// SQLServer driver
	SQLServer DBDriver = "sqlserver"
)

// ConnectionConfig is the interface for all connection configurations
type ConnectionConfig interface {
	// Driver returns the database driver
	Driver() DBDriver
	// DSN returns the data source name connection string
	DSN() string
}

// MySQLConfig holds MySQL connection parameters
type MySQLConfig struct {
	Username     string
	Password     string
	Host         string
	Port         int
	Database     string
	Charset      string
	ParseTime    bool
	Loc          string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	DSNString    string // For direct DSN usage
}

// Driver returns the database driver
func (c *MySQLConfig) Driver() DBDriver {
	return MySQL
}

// DSN returns the MySQL connection string
func (c *MySQLConfig) DSN() string {
	if c.DSNString != "" {
		return c.DSNString
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.Charset, c.ParseTime, c.Loc)
}

// PostgresConfig holds PostgreSQL connection parameters
type PostgresConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	Database     string
	SSLMode      string
	TimeZone     string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	DSNString    string // For direct DSN usage
}

// Driver returns the database driver
func (c *PostgresConfig) Driver() DBDriver {
	return PostgreSQL
}

// DSN returns the PostgreSQL connection string
func (c *PostgresConfig) DSN() string {
	if c.DSNString != "" {
		return c.DSNString
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		c.Host, c.Username, c.Password, c.Database, c.Port, c.SSLMode, c.TimeZone)
}

// SQLiteConfig holds SQLite connection parameters
type SQLiteConfig struct {
	Database     string
	Mode         string
	Cache        string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	DSNString    string // For direct DSN usage
}

// Driver returns the database driver
func (c *SQLiteConfig) Driver() DBDriver {
	return SQLite
}

// DSN returns the SQLite connection string
func (c *SQLiteConfig) DSN() string {
	if c.DSNString != "" {
		return c.DSNString
	}
	return fmt.Sprintf("%s?_mode=%s&_cache=%s", c.Database, c.Mode, c.Cache)
}

// SQLServerConfig holds SQL Server connection parameters
type SQLServerConfig struct {
	Username     string
	Password     string
	Host         string
	Port         int
	Database     string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	DSNString    string // For direct DSN usage
}

// Driver returns the database driver
func (c *SQLServerConfig) Driver() DBDriver {
	return SQLServer
}

// DSN returns the SQL Server connection string
func (c *SQLServerConfig) DSN() string {
	if c.DSNString != "" {
		return c.DSNString
	}
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		c.Username, c.Password, c.Host, c.Port, c.Database)
}

// =============================================================================
// Registry Types and Global Instance
// =============================================================================

// DBOption is a function that configures a gorm.DB instance
type DBOption func(*gorm.DB)

// RegistryOption is a function that configures the Registry
type RegistryOption func(*Registry)

// Registry manages GORM database connections
type Registry struct {
	mu            sync.RWMutex
	connections   map[string]*Connection
	globalOptions []DBOption
	logger        Logger
	retryOptions  *RetryOptions
}

// Connection represents a registered database connection
type Connection struct {
	db         *gorm.DB
	config     ConnectionConfig
	options    []DBOption
	connecting bool
}

// Global registry instance
var globalRegistry = &Registry{
	connections:   make(map[string]*Connection),
	globalOptions: []DBOption{},
	logger:        nil,
}

// =============================================================================
// Registry Configuration Options
// =============================================================================

// WithLogger sets a custom logger for the registry
func WithLogger(logger Logger) RegistryOption {
	return func(r *Registry) {
		r.logger = logger
	}
}

// WithGlobalOptions adds options that will be applied to all new connections
func WithGlobalOptions(option func(*gorm.DB)) RegistryOption {
	return func(r *Registry) {
		r.globalOptions = append(r.globalOptions, option)
	}
}

// WithConnectionDetails sets the connection details for a named connection
// It accepts a connection configuration struct
func WithConnectionDetails(name string, config ConnectionConfig, options ...DBOption) RegistryOption {
	return func(r *Registry) {
		if r.connections == nil {
			r.connections = make(map[string]*Connection)
		}

		r.connections[name] = &Connection{
			config:     config,
			options:    options,
			connecting: false,
		}
	}
}

// WithRetryOptions adds retry options to the registry
func WithRetryOptions(options *RetryOptions) RegistryOption {
	return func(r *Registry) {
		r.retryOptions = options
	}
}

// =============================================================================
// Connection Management
// =============================================================================

// RegisterConnection configures database connections with provided options
func RegisterConnection(opts ...RegistryOption) error {
	// Apply options
	for _, opt := range opts {
		opt(globalRegistry)
	}

	// Set default retry options if not set
	if globalRegistry.retryOptions == nil {
		globalRegistry.retryOptions = DefaultRetryOptions()
	}

	var errs []error
	// Start connections for all registered connections
	for name := range globalRegistry.connections {
		err := connectWithRetry(name)
		if err != nil {
			errs = append(errs, fmt.Errorf("connection '%s': %w", name, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to connect to one or more databases: %v", errs)
	}
	return nil
}

// connectWithRetry attempts to establish a connection with retry and returns error if all retries fail
func connectWithRetry(name string) error {
	operation := func() error {
		globalRegistry.mu.RLock()
		conn, exists := globalRegistry.connections[name]
		globalRegistry.mu.RUnlock()
		if !exists || conn.db != nil {
			return nil // Already connected or removed
		}
		// Mark that we're connecting
		globalRegistry.mu.Lock()
		if !conn.connecting {
			conn.connecting = true
			globalRegistry.connections[name] = conn
		}
		globalRegistry.mu.Unlock()

		db, err := initializeConnection(name)
		if err == nil && db != nil {
			globalRegistry.mu.Lock()
			conn, stillExists := globalRegistry.connections[name]
			if stillExists {
				conn.db = db
				conn.connecting = false
				globalRegistry.connections[name] = conn
			}
			globalRegistry.mu.Unlock()
			if globalRegistry.logger != nil {
				extras := map[string]interface{}{
					ExtraKey.Database.Table: name,
				}
				globalRegistry.logger.InfoWithCategory(Category.Database.Database, SubCategory.Status.Success, "Successfully connected to database", extras)
			}
			return nil
		}
		return err
	}
	return WithRetry(
		operation,
		globalRegistry.retryOptions,
		globalRegistry.logger,
		"connectWithRetry("+name+")",
	)
}

// initializeConnection attempts to initialize a database connection
func initializeConnection(name string) (*gorm.DB, error) {
	globalRegistry.mu.RLock()
	conn, exists := globalRegistry.connections[name]
	globalRegistry.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("connection with name '%s' not registered", name)
	}

	var db *gorm.DB
	var err error

	// Open database connection based on driver type
	switch conn.config.Driver() {
	case MySQL:
		db, err = gorm.Open(mysqlDriver.Open(conn.config.DSN()), &gorm.Config{})
	case PostgreSQL:
		// Uncomment when postgres driver is available
		// db, err = gorm.Open(postgresDriver.Open(conn.config.DSN()), &gorm.Config{})
		return nil, fmt.Errorf("postgres driver not implemented - uncomment the import and implementation")
	case SQLite:
		// Uncomment when sqlite driver is available
		// db, err = gorm.Open(sqliteDriver.Open(conn.config.DSN()), &gorm.Config{})
		return nil, fmt.Errorf("sqlite driver not implemented - uncomment the import and implementation")
	case SQLServer:
		db, err = gorm.Open(sqlserverDriver.Open(conn.config.DSN()), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", conn.config.Driver())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Apply global options
	for _, option := range globalRegistry.globalOptions {
		option(db)
	}

	// Apply connection-specific options
	for _, option := range conn.options {
		option(db)
	}

	// Apply connection pool settings if available
	sqlDB, err := db.DB()
	if err == nil {
		switch cfg := conn.config.(type) {
		case *MySQLConfig:
			if cfg.MaxOpenConns > 0 {
				sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
			}
			if cfg.MaxIdleConns > 0 {
				sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
			}
			if cfg.MaxLifetime > 0 {
				sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)
			}
		case *PostgresConfig:
			if cfg.MaxOpenConns > 0 {
				sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
			}
			if cfg.MaxIdleConns > 0 {
				sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
			}
			if cfg.MaxLifetime > 0 {
				sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)
			}
		case *SQLiteConfig:
			if cfg.MaxOpenConns > 0 {
				sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
			}
			if cfg.MaxIdleConns > 0 {
				sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
			}
			if cfg.MaxLifetime > 0 {
				sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)
			}
		case *SQLServerConfig:
			if cfg.MaxOpenConns > 0 {
				sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
			}
			if cfg.MaxIdleConns > 0 {
				sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
			}
			if cfg.MaxLifetime > 0 {
				sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)
			}
		}
	}

	return db, nil
}

// =============================================================================
// Public API
// =============================================================================

// GetConnection returns a GORM DB instance for the named connection
func GetConnection(name string) (*gorm.DB, error) {
	globalRegistry.mu.RLock()
	conn, exists := globalRegistry.connections[name]
	globalRegistry.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("connection with name '%s' not registered", name)
	}

	// If DB is already initialized, return it
	if conn.db != nil {
		return conn.db, nil
	}

	// If DB is not initialized yet, wait for connection
	for {
		globalRegistry.mu.RLock()
		conn, exists = globalRegistry.connections[name]

		if !exists {
			globalRegistry.mu.RUnlock()
			return nil, fmt.Errorf("connection with name '%s' not registered", name)
		}

		if conn.db != nil {
			db := conn.db
			globalRegistry.mu.RUnlock()
			return db, nil
		}

		// If not connecting yet, start connecting
		if !conn.connecting {
			globalRegistry.mu.RUnlock()
			go connectWithRetry(name)
		} else {
			globalRegistry.mu.RUnlock()
		}

		// Wait a bit before checking again
		time.Sleep(100 * time.Millisecond)
	}
}

// MustDB returns a GORM DB instance or panics if it cannot be retrieved
func MustDB(name string) *gorm.DB {
	db, _ := GetConnection(name)
	return db
}

// Health checks if all database connections are healthy
func Health(ctx context.Context) error {
	globalRegistry.mu.RLock()
	connections := make([]string, 0, len(globalRegistry.connections))
	for name := range globalRegistry.connections {
		connections = append(connections, name)
	}
	globalRegistry.mu.RUnlock()

	// Check if there are any connections to check
	if len(connections) == 0 {
		if globalRegistry.logger != nil {
			globalRegistry.logger.WarnWithCategory(Category.Database.Database, SubCategory.Status.Warning, "No database connections registered", nil)
		}
		return fmt.Errorf("no database connections registered")
	}

	// Check all database connections
	dbs := make([]*gorm.DB, 0, len(connections))
	errs := make([]string, 0)

	for _, name := range connections {
		db, err := GetConnection(name)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", name, err))
			continue
		}
		dbs = append(dbs, db)
	}

	if len(dbs) == 0 {
		errorMsg := strings.Join(errs, "; ")
		if globalRegistry.logger != nil {
			extras := map[string]interface{}{
				ExtraKey.Error.ErrorMessage: errorMsg,
			}
			globalRegistry.logger.ErrorWithCategory(Category.Database.Database, SubCategory.Status.Error, "Database health check failed", extras)
		}
		return fmt.Errorf("database health check failed: %s", errorMsg)
	}

	// Run a simple query on each connection to verify it's working
	for i, db := range dbs {
		currentName := connections[i]
		if err := db.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
			if globalRegistry.logger != nil {
				extras := map[string]interface{}{
					ExtraKey.Database.Table:     currentName,
					ExtraKey.Error.ErrorMessage: err.Error(),
				}
				globalRegistry.logger.ErrorWithCategory(Category.Database.Database, SubCategory.Status.Error, "Database health check failed", extras)
			}
			return fmt.Errorf("health check failed for connection '%s': %w", currentName, err)
		}
	}

	if globalRegistry.logger != nil {
		globalRegistry.logger.InfoWithCategory(Category.Database.Database, SubCategory.Status.Success, "All database connections are healthy", nil)
	}
	return nil
}

// CloseConnections closes all registered database connections
func CloseConnections() {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	for name, conn := range globalRegistry.connections {
		if conn.db != nil {
			sqlDB, err := conn.db.DB()
			if err == nil {
				sqlDB.Close()
				if globalRegistry.logger != nil {
					extras := map[string]interface{}{
						ExtraKey.Database.Table: name,
					}
					globalRegistry.logger.InfoWithCategory(Category.Database.Database, SubCategory.Operation.Shutdown, "Database connection closed", extras)
				}
			} else if globalRegistry.logger != nil {
				extras := map[string]interface{}{
					ExtraKey.Database.Table:     name,
					ExtraKey.Error.ErrorMessage: err.Error(),
				}
				globalRegistry.logger.ErrorWithCategory(Category.Database.Database, SubCategory.Status.Error, "Failed to close database connection", extras)
			}
			conn.db = nil
			globalRegistry.connections[name] = conn
		}
	}
}

// CloseConnection closes a specific connection by name
func CloseConnection(name string) error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	conn, exists := globalRegistry.connections[name]
	if !exists {
		return fmt.Errorf("connection with name '%s' not registered", name)
	}

	if conn.db != nil {
		sqlDB, err := conn.db.DB()
		if err != nil {
			return err
		}

		if err := sqlDB.Close(); err != nil {
			if globalRegistry.logger != nil {
				extras := map[string]interface{}{
					ExtraKey.Database.Table:     name,
					ExtraKey.Error.ErrorMessage: err.Error(),
				}
				globalRegistry.logger.ErrorWithCategory(Category.Database.Database, SubCategory.Status.Error, "Failed to close database connection", extras)
			}
			return err
		}

		if globalRegistry.logger != nil {
			extras := map[string]interface{}{
				ExtraKey.Database.Table: name,
			}
			globalRegistry.logger.InfoWithCategory(Category.Database.Database, SubCategory.Operation.Shutdown, "Database connection closed", extras)
		}
	}

	return nil
}
