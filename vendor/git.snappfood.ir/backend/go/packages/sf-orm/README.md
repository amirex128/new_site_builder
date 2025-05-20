# SF-ORM

A resilient GORM database connection registry for Go.

[![Go Reference](https://pkg.go.dev/badge/git.snappfood.ir/backend/go/packages/sf-orm.svg)](https://pkg.go.dev/git.snappfood.ir/backend/go/packages/sf-orm)
[![Go Report Card](https://goreportcard.com/badge/git.snappfood.ir/backend/go/packages/sf-orm)](https://goreportcard.com/report/git.snappfood.ir/backend/go/packages/sf-orm)

This library provides a registry for managing multiple GORM database connections in Go applications:

- **Resilient Connection Management**: Automatic retry and reconnection with exponential backoff, panic recovery, and error logging
- **Multiple Database Support**: MySQL, PostgreSQL, SQLite, SQL Server
- **Type-Safe Configuration**: Strongly typed connection configs with validation
- **Service Registry Integration**: Works with SF Service Registry
- **APM Integration**: Built-in support for Elastic APM
- **Health Checks**: Comprehensive health check functionality

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Retry Options and Exponential Backoff](#retry-options-and-exponential-backoff)
- [Connection Registry](#connection-registry)
- [Connection Configuration](#connection-configuration)
- [Service Connector](#service-connector)
- [Health Checks](#health-checks)
- [Features](#features)
- [Database Operations](#database-operations)
- [Connection Options](#connection-options)
- [API Reference](#api-reference)
- [Implementation Notes](#implementation-notes)
- [License](#license)

## Installation

```bash
go get git.snappfood.ir/backend/go/packages/sf-orm
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "git.snappfood.ir/backend/go/packages/sf-orm"
)

type User struct {
    ID   uint   `gorm:"primarykey"`
    Name string
}

func main() {
    // Create MySQL configuration directly as a struct
    mysqlConfig := &sform.MySQLConfig{
        Username:     "user",
        Password:     "password",
        Host:         "127.0.0.1",
        Port:         3306,
        Database:     "dbname",
        Charset:      "utf8mb4",
        ParseTime:    true,
        Loc:          "Local",
        MaxOpenConns: 10,
        MaxIdleConns: 5,
        MaxLifetime:  5 * time.Minute,
    }

    // Create PostgreSQL configuration directly as a struct
    pgConfig := &sform.PostgresConfig{
        Username:     "postgres",
        Password:     "postgres",
        Host:         "localhost",
        Port:         5432,
        Database:     "testdb",
        SSLMode:      "disable",
        TimeZone:     "Asia/Tehran",
    }

    // Register your database connections with meaningful names and options
    err := sform.RegisterConnection(
        sform.WithLogger(&MyLogger{}),
        sform.WithGlobalOptions(func(db *gorm.DB) {
            db.Debug()
        }),
        sform.WithConnectionDetails("main", mysqlConfig),
        sform.WithConnectionDetails("secondary", pgConfig),
        sform.WithRetryOptions(&sform.RetryOptions{
            MaxRetries:     5,
            InitialBackoff: time.Second,
            MaxBackoff:     15 * time.Second,
            BackoffFactor:  1.5,
        }),
    )
    
    if err != nil {
        log.Fatalf("Failed to register database connection: %v", err)
    }

    // Get the database connection by name
    db, err := sform.GetConnection("main")
    if err != nil {
        log.Fatalf("Failed to get database connection: %v", err)
    }
    
    // Use GORM as usual with the connection
    var users []User
    db.Find(&users)
    
    fmt.Printf("Found %d users\n", len(users))
}
```

## Retry Options and Exponential Backoff

SF-ORM supports robust retry logic with exponential backoff and panic recovery for all external data source connections. You can configure retry behavior globally for all connections using `WithRetryOptions`:

```go
err := sform.RegisterConnection(
    sform.WithConnectionDetails("main", mysqlConfig),
    sform.WithRetryOptions(&sform.RetryOptions{
        MaxRetries:     5,                // Maximum number of retry attempts
        InitialBackoff: time.Second,      // Initial waiting time between retries
        MaxBackoff:     15 * time.Second, // Maximum waiting time between retries
        BackoffFactor:  1.5,              // Exponential backoff multiplier
    }),
)
```

- **MaxRetries**: Maximum number of retry attempts before giving up
- **InitialBackoff**: Initial delay before the first retry
- **MaxBackoff**: Maximum delay between retries
- **BackoffFactor**: Multiplier for exponential backoff

All connection attempts are protected with panic recovery and errors are logged using the provided logger.

### ServiceConnector with Retry Options

You can also use retry options with the ServiceConnector:

```go
connector := sform.GetConnector()
err := connector.RegisterConnection(
    sform.WithConnectionName("main"),
    sform.WithDriver(sform.MySQL),
    sform.WithHost("localhost"),
    sform.WithPort(3306),
    sform.WithUser("user"),
    sform.WithPassword("password"),
    sform.WithDatabase("mydb"),
    sform.WithRetryOptions(&sform.RetryOptions{
        MaxRetries:     5,
        InitialBackoff: time.Second,
        MaxBackoff:     10 * time.Second,
        BackoffFactor:  2.0,
    }),
)
```

## Connection Registry

The connection registry is the core of SF-ORM. It allows you to register database connections by name in one place (typically in main.go), then reuse them throughout your application without creating new connections each time.

```go
package main

import (
    "log"
    "time"
    
    "git.snappfood.ir/backend/go/packages/sf-orm"
)

func main() {
    // Create MySQL configuration directly
    mainConfig := &sform.MySQLConfig{
        Username:     "user",
        Password:     "password",
        Host:         "127.0.0.1",
        Port:         3306,
        Database:     "dbname",
        Charset:      "utf8mb4",
        ParseTime:    true,
        Loc:          "Local",
    }
    
    // Create PostgreSQL configuration directly
    analyticsConfig := &sform.PostgresConfig{
        Username:     "postgres",
        Password:     "postgres",
        Host:         "localhost",
        Port:         5432,
        Database:     "analytics",
        SSLMode:      "disable",
        TimeZone:     "UTC",
    }
    
    // Register connections with options
    err := sform.RegisterConnection(
        sform.WithLogger(&MyLogger{}),
        sform.WithGlobalOptions(func(db *gorm.DB) {
            db.Debug()
        }),
        sform.WithConnectionDetails("main", mainConfig),
        sform.WithConnectionDetails("analytics", analyticsConfig),
        sform.WithRetryOptions(&sform.RetryOptions{
            MaxRetries:     5,
            InitialBackoff: time.Second,
            MaxBackoff:     15 * time.Second,
            BackoffFactor:  1.5,
        }),
    )
    
    if err != nil {
        log.Fatalf("Failed to register connection: %v", err)
    }
    
    // Start application...
    
    // Close all connections when application exits
    defer sform.CloseConnections()
}
```

## Connection Configuration

SF-ORM supports multiple database types with specific configuration structs for each:

### MySQL Configuration

```go
// Create a MySQL configuration directly
mysqlConfig := &sform.MySQLConfig{
    Username:     "user",
    Password:     "password",
    Host:         "127.0.0.1",
    Port:         3306,
    Database:     "dbname",
    Charset:      "utf8mb4",      // Default: utf8mb4
    ParseTime:    true,           // Default: true
    Loc:          "Local",        // Default: Local
    MaxOpenConns: 10,             // Connection pool settings
    MaxIdleConns: 5,
    MaxLifetime:  5 * time.Minute,
}

// Register with the connection registry
sform.RegisterConnection(
    sform.WithConnectionDetails("mysql_conn", mysqlConfig),
)
```

### PostgreSQL Configuration

```go
// Create a PostgreSQL configuration directly
pgConfig := &sform.PostgresConfig{
    Username:     "postgres",
    Password:     "postgres",
    Host:         "localhost",
    Port:         5432,
    Database:     "mydb",
    SSLMode:      "disable",         // Default: disable
    TimeZone:     "UTC",             // Default: UTC
    MaxOpenConns: 10,
    MaxIdleConns: 5,
    MaxLifetime:  5 * time.Minute,
}

// Register with the connection registry
sform.RegisterConnection(
    sform.WithConnectionDetails("pg_conn", pgConfig),
)
```

### SQLite Configuration

```go
// Create a SQLite configuration directly
sqliteConfig := &sform.SQLiteConfig{
    Database:     "/path/to/db.sqlite",
    Mode:         "rwc",            // Default: rwc
    Cache:        "shared",         // Default: shared
}

// Register with the connection registry
sform.RegisterConnection(
    sform.WithConnectionDetails("sqlite_conn", sqliteConfig),
)
```

### SQL Server Configuration

```go
// Create a SQL Server configuration directly
sqlServerConfig := &sform.SQLServerConfig{
    Username:     "sa",
    Password:     "yourStrong(!)Password",
    Host:         "localhost",
    Port:         1433,
    Database:     "master",
}

// Register with the connection registry
sform.RegisterConnection(
    sform.WithConnectionDetails("sqlserver_conn", sqlServerConfig),
)
```

## Service Connector

SF-ORM provides a ServiceConnector that implements the serviceregistry.ServiceConnector interface, making it easy to integrate with the SF Service Registry.

```go
package main

import (
    "log"
    
    "git.snappfood.ir/backend/go/packages/sf-service-registry"
    "git.snappfood.ir/backend/go/packages/sf-orm"
)

func main() {
    // Get the connector from sf-orm
    connector := sform.GetConnector()
    
    // Register services with the service registry
    err := serviceregistry.RegisterServices(
        connector, // MySQL
        nil,       // No Redis
        nil,       // No Elastic
        nil,       // No RabbitMQ
        nil,       // No SQL Server
    )
    if err != nil {
        log.Fatalf("Failed to register services: %v", err)
    }
    
    // Now you can use the database connections
    db, err := sform.GetConnection("read_db")
    if err != nil {
        log.Fatalf("Failed to get database connection: %v", err)
    }
    
    // Use the connection...
}
```

### Customizing the Service Connector

The service connector now supports typed connection parameters:

```go
// Get the connector
connector := sform.GetConnector()

// Configure MySQL connection with typesafe options
err := connector.RegisterConnection(
    sform.WithConnectionName("read_db"),
    sform.WithDriver(sform.MySQL),
    sform.WithHost("localhost"),
    sform.WithPort(3306),
    sform.WithUser("user"),
    sform.WithPassword("password"),
    sform.WithDatabase("mydb"),
    sform.WithCharset("utf8mb4"),
    sform.WithParseTime(true),
    sform.WithMaxOpenConns(10),
    sform.WithMaxIdleConns(5),
    sform.WithMaxLifetime(30), // seconds
)

// Configure PostgreSQL connection
err = connector.RegisterConnection(
    sform.WithConnectionName("analytics_db"),
    sform.WithDriver(sform.PostgreSQL),
    sform.WithHost("pg-server"),
    sform.WithPort(5432),
    sform.WithUser("postgres"),
    sform.WithPassword("postgres"),
    sform.WithDatabase("analytics"),
    sform.WithSSLMode("disable"),
    sform.WithTimeZone("UTC"),
)
```

## Health Checks

SF-ORM provides a Health function that checks if all database connections are alive:

```go
package main

import (
    "context"
    "log"
    "time"
    
    "git.snappfood.ir/backend/go/packages/sf-orm"
)

func main() {
    // Register your connections...
    
    // Check health of all connections
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    err := sform.Health(ctx)
    if err != nil {
        log.Fatalf("Health check failed: %v", err)
    }
    
    log.Println("All database connections are healthy")
}
```

### Health Checks with Gin Framework

You can integrate health checks with a Gin web server:

```go
package main

import (
    "context"
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    
    "git.snappfood.ir/backend/go/packages/sf-orm"
)

func main() {
    // Register your connections...
    
    // Set up Gin router
    r := gin.Default()
    
    // Health check endpoint
    r.GET("/health", func(c *gin.Context) {
        // Create a context with timeout for health checks
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        
        // Check database connections
        if err := sform.Health(ctx); err != nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{
                "status": "unhealthy",
                "error": err.Error(),
            })
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "status": "healthy",
        })
    })
    
    // Start server
    r.Run(":8080")
}
```

## Features

- **Named Connection Registry** - Register multiple database connections by name
- **Type-Safe Configuration** - Strong typed configuration objects for each database type
- **Flexible Connection Setup** - Use either structured configs or connection strings
- **Connection Pool Settings** - Configure max open/idle connections and max lifetime
- **Lazy Connection Initialization** - Connections are only established when first used
- **Connection-Specific Options** - Set custom options for each connection
- **Global Options** - Apply settings to all connections
- **Thread-Safe Implementation** - Proper locking for concurrent access
- **Multiple Database Support** - MySQL, PostgreSQL, SQLite, SQL Server with constants
- **Connection Management** - Close and reopen connections as needed
- **APM Integration** - Built-in support for Elastic APM
- **Health Checks** - Comprehensive health check functionality
- **Service Registry Integration** - Works with SF Service Registry
- **Resilient Connection Management** - Automatic retry and reconnection with exponential backoff, panic recovery, and error logging

## Database Operations

When you have registered your connections, you can use them anywhere in your application:

```go
package repository

import (
    "git.snappfood.ir/backend/go/packages/sf-orm"
)

type User struct {
    ID   uint   `gorm:"primarykey"`
    Name string
}

type UserRepository struct {
    // No need to store the connection, get it when needed
}

func NewUserRepository() *UserRepository {
    return &UserRepository{}
}

func (r *UserRepository) FindByID(id uint) (*User, error) {
    // Get the database connection by name
    db, err := sform.GetConnection("main")
    if err != nil {
        return nil, err
    }
    
    var user User
    result := db.First(&user, id)
    if result.Error != nil {
        return nil, result.Error
    }
    
    return &user, nil
}

func (r *UserRepository) Create(user *User) error {
    // Use MustDB if you're sure the connection exists and want to panic if not
    db := sform.MustDB("main")
    
    result := db.Create(user)
    return result.Error
}

func (r *UserRepository) FindWithRawSQL() ([]User, error) {
    db, err := sform.GetConnection("main")
    if err != nil {
        return nil, err
    }
    
    var users []User
    db.Raw("SELECT * FROM users WHERE active = ?", true).Scan(&users)
    
    return users, nil
}
```

## Connection Options

```go
// Create database configuration directly
mysqlConfig := &sform.MySQLConfig{
    Username:     "user",
    Password:     "password",
    Host:         "127.0.0.1",
    Port:         3306,
    Database:     "dbname",
    // Configure connection pool
    MaxOpenConns: 10,
    MaxIdleConns: 5,
    MaxLifetime:  5 * time.Minute,
}

// Register the connection
sform.RegisterConnection(
    sform.WithConnectionDetails("main", mysqlConfig),
    sform.WithGlobalOptions(func(db *gorm.DB) {
        db.Debug()
    }),
)

// Configure GORM logger
customLogger := logger.New(
    log.New(os.Stdout, "\r\n", log.LstdFlags),
    logger.Config{
        SlowThreshold:             200 * time.Millisecond,
        LogLevel:                  logger.Info,
        IgnoreRecordNotFoundError: false,
        Colorful:                  true,
    },
)

// Create debug database config directly
debugConfig := &sform.MySQLConfig{
    Username:     "user",
    Password:     "password",
    Host:         "127.0.0.1",
    Port:         3306,
    Database:     "debug_db",
}

sform.RegisterConnection(
    sform.WithLogger(customLogger),
    sform.WithConnectionDetails("debug-db", debugConfig),
    sform.WithGlobalOptions(func(db *gorm.DB) {
        db.Debug()
    }),
)
```

## API Reference

### Database Drivers

- `MySQL DBDriver = "mysql"` - MySQL/MariaDB driver
- `PostgreSQL DBDriver = "postgres"` - PostgreSQL driver
- `SQLite DBDriver = "sqlite"` - SQLite driver
- `SQLServer DBDriver = "sqlserver"` - SQL Server driver

### Connection Configuration

- `MySQLConfig` - MySQL configuration struct
- `PostgresConfig` - PostgreSQL configuration struct
- `SQLiteConfig` - SQLite configuration struct
- `SQLServerConfig` - SQL Server configuration struct

### Connection Registration and Management

- `RegisterConnection(opts ...RegistryOption) error`
- `GetConnection(name string) (*gorm.DB, error)`
- `MustDB(name string) *gorm.DB`
- `CloseConnection(name string) error`
- `CloseConnections()`
- `Health(ctx context.Context) error`

### Registry Options

- `WithLogger(logger Logger) RegistryOption`
- `WithGlobalOptions(options ...DBOption) RegistryOption`
- `WithConnectionDetails(name string, driverOrConfig interface{}, dsnOrOptions ...interface{}) RegistryOption`
- `WithRetryOptions(options *RetryOptions) RegistryOption` — Set global retry options for all connections

### Service Connector

- `GetConnector() *ServiceConnector`
- `ServiceConnector.WithOptions(opts ...RegistryOption) *ServiceConnector`
- `ServiceConnector.Health(ctx context.Context) error`

### Service Connector Options

- `WithConnectionName(name string) ServiceConnectorOption`
- `WithHost(host string) ServiceConnectorOption`
- `WithPort(port int) ServiceConnectorOption`
- `WithUser(user string) ServiceConnectorOption`
- `WithPassword(password string) ServiceConnectorOption`
- `WithDatabase(database string) ServiceConnectorOption`
- `WithDriver(driver DBDriver) ServiceConnectorOption`
- `WithCharset(charset string) ServiceConnectorOption`
- `WithParseTime(parseTime bool) ServiceConnectorOption`
- `WithLocale(loc string) ServiceConnectorOption`
- `WithSSLMode(sslMode string) ServiceConnectorOption`
- `WithTimeZone(timeZone string) ServiceConnectorOption`
- `WithMaxOpenConns(maxOpen int) ServiceConnectorOption`
- `WithMaxIdleConns(maxIdle int) ServiceConnectorOption`
- `WithMaxLifetime(maxLife int) ServiceConnectorOption`

### RetryOptions struct

```
type RetryOptions struct {
    MaxRetries     int           // Maximum number of retry attempts
    InitialBackoff time.Duration // Initial waiting time between retries
    MaxBackoff     time.Duration // Maximum waiting time between retries
    BackoffFactor  float64       // Exponential backoff multiplier
}
```

## Implementation Notes

To use the registry with specific database drivers, you need to uncomment and correctly import the database drivers you need in your code:

```go
import (
    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/driver/sqlite"
    "gorm.io/driver/sqlserver"
)
```

The following database drivers are supported:

- `MySQL` - MySQL/MariaDB
- `PostgreSQL` - PostgreSQL
- `SQLite` - SQLite
- `SQLServer` - Microsoft SQL Server

### Configuration

Services are configured in `services.json`:

```json
{
  "mysql": [
    {
      "name": "read_db",
      "env_host": "host",
      "env_port": "port",
      "env_database": "db",
      "env_user": "user",
      "env_password": "password"
    }
  ],
  "redis": [
    {
      "name": "search_redis",
      "env_host": "host",
      "env_password": "password",
      "env_database": "db",
      "kind": "normal"
    }
  ],
  "elastic": [
    {
      "name": "base_elastic",
      "env_host": "host",
      "env_port": "port",
      "env_insecure_tls": "insecure_tls",
      "env_username": "username",
      "env_password": "password"
    }
  ]
}
```

The system uses Viper to read environment variables based on the configuration. For example, for a MySQL service named `read_db`, it will look for environment variables like:
- `read_db_host`
- `read_db_port`
- `read_db_user`
- `read_db_password`
- `read_db_db`

## License

MIT