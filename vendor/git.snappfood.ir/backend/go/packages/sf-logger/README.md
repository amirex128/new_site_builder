# SF-Logger

A powerful, extensible, structured logging package for Go applications with multiple output destinations and comprehensive configuration options.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Features](#features)
- [Logger Architecture](#logger-architecture)
- [Configuration Options](#configuration-options)
- [Log Levels](#log-levels)
- [API Reference](#api-reference)
- [Context Support](#context-support)
- [Category-Based Logging](#category-based-logging)
- [URL-Based Sink Registry](#url-based-sink-registry)
- [Built-in Sinks](#built-in-sinks)
- [Custom Sinks](#custom-sinks)
- [Logger Types](#logger-types)
- [Examples](#examples)

## Installation

```bash
go get git.snappfood.ir/backend/go/packages/sf-logger
```

## Quick Start

```go
package main

import (
    "git.snappfood.ir/backend/go/packages/sf-logger"
    "context"
)

func main() {
    // Create a simple logger with URL-based sink
    logger := sflogger.New(
        sflogger.WithAppName("MyService"),
        sflogger.WithLevel(sflogger.InfoLevel),
        sflogger.WithSinkURL("console://"),
    )
    
    // Log simple messages
    logger.Info("Application started", map[string]interface{}{
        "version": "1.0.0",
        "environment": "production",
    })
    
    // Log with categories
    logger.InfoWithCategory("General", "Startup", "Server initialized", 
        map[string]interface{}{
            "hostIp": "192.168.1.1",
            "port": "8080",
        })
        
    // Use the global logger
    sflogger.InitGlobalLogger(
        sflogger.WithAppName("MyService"),
        sflogger.WithLevel(sflogger.InfoLevel),
        sflogger.WithSinkURL("console://"),
    )
    
    sflogger.Info("Using global logger", nil)
}
```

## Features

- **Multiple Logger Implementations** - Choose between Zap and Zerolog
- **Structured Logging** - Log both messages and structured data as key-value pairs
- **URL-Based Sink System** - Configure sinks using URL syntax for better clarity
- **Multiple Output Destinations** - Console, Files, Elasticsearch, Graylog, Logstash, Loki
- **Extensible Sink Registry** - Register custom sinks with a simple interface
- **Log Levels** - Debug, Info, Warning, Error, and Fatal levels
- **Multiple Formatters** - JSON, Pretty JSON, Colored Text, Plain Text, CSV, and Custom JSON
- **Context Support** - Context-aware logging with automatic extraction of common fields
- **Category-Based Logging** - Organize logs by category and subcategory for better filtering
- **Global Logger** - Simple package-level logging functions
- **Log Rotation** - Automatic log file rotation with size and age limits
- **Functional Options API** - Clear, type-safe configuration with good defaults
- **Thread-Safe** - Safe for concurrent use
- **Stack Traces** - Automatic stack traces for error-level logs

## Logger Architecture

SF-Logger uses a layered architecture for flexibility and extensibility:

1. **Core Interface** - Defines the Logger interface with common methods
2. **Implementations** - Different logger implementations (Zap, Zerolog)
3. **Formatters** - Controls how logs are formatted (JSON, Text, etc.)
4. **Sinks** - Manages where logs are sent using a URL-based registry system
5. **Functional Options** - Provides a clean API for configuration

## Configuration Options

SF-Logger uses a functional options pattern for clean, flexible configuration:

```go
logger := sflogger.New(
    // Core settings
    sflogger.WithLoggerType(sflogger.ZapLoggerType),
    sflogger.WithLevel(sflogger.InfoLevel),
    sflogger.WithAppName("MyService"),
    
    // Formatting options
    sflogger.WithFormatter(sflogger.JSONFormatter),
    sflogger.WithTimeFormat("2006-01-02T15:04:05.000Z07:00"),
    
    // Development and debugging
    sflogger.WithDevelopment(true),
    sflogger.WithStacktrace(true),
    
    // URL-based sink configuration
    sflogger.WithSinkURL("console://?color=true"),
    sflogger.WithSinkURL("file:///logs/app.log?maxSize=10&maxAge=30&maxBackups=5&compress=true"),
    sflogger.WithSinkURL("elasticsearch://elastic.example.com:9200?index=app-logs&username=user&password=pass"),
)
```

### Available Configuration Options

| Option | Description |
|--------|-------------|
| WithLoggerType | Choose between ZapLoggerType and ZeroLoggerType |
| WithLevel | Set the minimum log level |
| WithAppName | Set the application name |
| WithFormatter | Set the log formatter type |
| WithDevelopment | Enable development mode |
| WithStacktrace | Enable stack traces for errors |
| WithTimeFormat | Set time format string |
| WithSinkURL | Add a sink using URL syntax |
| WithSinkURLs | Set multiple sinks using URL syntax |
| WithConsoleSink | Add a console sink with optional color |
| WithFileSink | Add a file sink with rotation settings |
| WithGraylogSink | Add a Graylog sink |
| WithLogstashSink | Add a Logstash sink |
| WithLokiSink | Add a Grafana Loki sink |
| WithElasticsearchSink | Add an Elasticsearch sink |
| WithMultiSink | Configure multiple sinks with routing options |
| WithCustomSink | Add a custom sink with specific scheme |

## Log Levels

SF-Logger supports the following log levels in order of increasing severity:

| Level | Description |
|-------|-------------|
| DebugLevel | Debug information, verbose output |
| InfoLevel | General operational information |
| WarningLevel | Warning conditions, non-critical issues |
| ErrorLevel | Error conditions, operation failed |
| FatalLevel | Critical errors causing program termination |

## API Reference

### Core Interface

```go
type Logger interface {
    // Core logging methods
    Debug(msg string, extra map[string]interface{})
    Info(msg string, extra map[string]interface{})
    Warn(msg string, extra map[string]interface{})
    Error(msg string, extra map[string]interface{})
    Fatal(msg string, extra map[string]interface{})

    // Formatted logging methods
    Debugf(template string, args ...interface{})
    Infof(template string, args ...interface{})
    Warnf(template string, args ...interface{})
    Errorf(template string, args ...interface{})
    Fatalf(template string, args ...interface{})

    // Context-aware logging methods
    DebugContext(ctx context.Context, msg string, extra map[string]interface{})
    InfoContext(ctx context.Context, msg string, extra map[string]interface{})
    WarnContext(ctx context.Context, msg string, extra map[string]interface{})
    ErrorContext(ctx context.Context, msg string, extra map[string]interface{})
    FatalContext(ctx context.Context, msg string, extra map[string]interface{})

    // Category-based logging methods
    DebugWithCategory(cat string, sub string, msg string, extra map[string]interface{})
    InfoWithCategory(cat string, sub string, msg string, extra map[string]interface{})
    WarnWithCategory(cat string, sub string, msg string, extra map[string]interface{})
    ErrorWithCategory(cat string, sub string, msg string, extra map[string]interface{})
    FatalWithCategory(cat string, sub string, msg string, extra map[string]interface{})
}
```

### Sink Interface

```go
type Sink interface {
    // Write sends a log entry to the sink
    Write(entry map[string]interface{}) error

    // Close cleans up resources used by the sink
    Close() error
    
    // Sync flushes any buffered log entries
    Sync() error
}
```

### Global Helpers

```go
// Initialize the global logger
sflogger.InitGlobalLogger(options...)

// Use global helper functions
sflogger.Debug("Debug message", extra)
sflogger.Info("Info message", extra)
sflogger.Warn("Warning message", extra)
sflogger.Error("Error message", extra)
sflogger.Fatal("Fatal message", extra)

// Format strings
sflogger.Debugf("User %s logged in", username)
sflogger.Infof("Process took %d ms", duration)

// With context
sflogger.InfoContext(ctx, "Request processed", extra)

// With categories
sflogger.InfoWithCategory("General", "Startup", "Service starting", extra)
```

## Context Support

SF-Logger provides context-aware logging to support distributed tracing and request handling:

```go
// Add values to context
ctx = sflogger.WithRequestID(ctx, "req-123")
ctx = sflogger.WithTraceID(ctx, "trace-456")
ctx = sflogger.WithUserID(ctx, "user-789")

// Log with context
logger.InfoContext(ctx, "Processing request", extra)
```

Context values are automatically extracted and added to log entries.

## Category-Based Logging

SF-Logger supports categorizing logs for better organization and filtering:

```go
logger.InfoWithCategory("General", "Startup", "Service starting", 
    map[string]interface{}{
        "appName": "MyService",
        "hostIp": "192.168.1.1",
    })
```

## URL-Based Sink Registry

SF-Logger features a powerful URL-based sink registry system that makes configuring log destinations cleaner and more flexible.

### Basic Usage

```go
logger := sflogger.New(
    sflogger.WithAppName("MyApp"),
    // Add sinks using URLs
    sflogger.WithSinkURL("console://?color=true"),
    sflogger.WithSinkURL("file:///logs/app.log?maxSize=10&maxAge=30"),
)
```

### How the Registry Works

The sink registry uses a factory pattern to create sink instances from URLs:

```go
// Register a sink factory
sflogger.RegisterSink("myscheme", func(u *url.URL) (sflogger.Sink, error) {
    // Parse URL and create sink
    return &MySinkImpl{}, nil
})

// Get a sink from a URL
sink, err := sflogger.GetSink("myscheme://example.com?param=value")
```

### Sink URL Format

Sink URLs have the following structure:

```
scheme://[host[:port]][/path][?param1=value1&param2=value2...]
```

- **scheme**: Identifies the sink type (e.g., console, file, elasticsearch)
- **host/port**: For network sinks, specifies the destination server
- **path**: For file-based sinks, specifies the file path
- **query parameters**: Configure sink-specific options

## Built-in Sinks

### Console Sink
```
console://[?color=true|false&stderr=true|false&timeFormat=format]
```

Parameters:
- `color`: Enable colored output (default: false)
- `stderr`: Output to stderr instead of stdout (default: false)
- `timeFormat`: Custom time format string

Example:
```go
sflogger.WithSinkURL("console://?color=true")
```

Helper function:
```go
sflogger.WithConsoleSink(true) // with color
```

### File Sink
```
file:///path/to/file[?maxSize=10&maxAge=30&maxBackups=5&compress=true|false]
```

Parameters:
- `maxSize`: Maximum file size in megabytes before rotation (default: 100)
- `maxAge`: Maximum age in days to keep files (default: 30)
- `maxBackups`: Maximum number of old files to keep (default: 5)
- `compress`: Whether to compress rotated files (default: true)

Example:
```go
sflogger.WithSinkURL("file:///logs/app.log?maxSize=10&maxAge=30&maxBackups=5&compress=true")
```

Helper function:
```go
sflogger.WithFileSink("/logs/app.log", 10, 30, 5, true)
```

### Elasticsearch Sink
```
elasticsearch://host:port[?index=logs&username=user&password=pass&batchSize=100&flushTime=5s]
```

Parameters:
- `index`: Elasticsearch index name (default: "logs")
- `username`: Authentication username
- `password`: Authentication password
- `batchSize`: Number of logs to batch before sending (default: 100)
- `flushTime`: How often to flush partial batches (default: "5s")
- `timeout`: HTTP request timeout (default: "10s")

Example:
```go
sflogger.WithSinkURL("elasticsearch://elasticsearch:9200?index=app-logs&username=user&password=pass")
```

Helper function:
```go
sflogger.WithElasticsearchSink("elasticsearch:9200", "app-logs", "user", "pass", 100)
```

### Graylog Sink
```
graylog://host:port          # UDP protocol (default)
graylog+tcp://host:port      # TCP protocol
graylog+http://host:port     # HTTP protocol
graylog+https://host:port    # HTTPS protocol
```

Parameters:
- `timeout`: Network timeout (default: "5s")

Example:
```go
sflogger.WithSinkURL("graylog://graylog:12201")
```

Helper function:
```go
sflogger.WithGraylogSink("graylog.example.com", 12201, true) // true for UDP
```

### Logstash Sink
```
logstash://host:port        # UDP protocol (default)
logstash+tcp://host:port    # TCP protocol
logstash+http://host:port   # HTTP protocol
logstash+https://host:port  # HTTPS protocol
```

Parameters:
- `timeout`: Network timeout (default: "5s")

Example:
```go
sflogger.WithSinkURL("logstash+tcp://logstash:5044")
```

Helper function:
```go
sflogger.WithLogstashSink("logstash.example.com", 5044, false) // false for TCP
```

### Loki Sink
```
loki://host:port[?apiKey=key]
loki+https://host:port[?apiKey=key]
```

Parameters:
- `apiKey`: Authentication API key
- `timeout`: Network timeout (default: "5s")

Example:
```go
sflogger.WithSinkURL("loki://loki:3100?apiKey=secret")
```

Helper function:
```go
sflogger.WithLokiSink("loki.example.com:3100", "apikey")
```

### Multi Sink
```
multi://[?failSafe=true|false&sink=sink1&sink=sink2...]
```

Parameters:
- `failSafe`: Whether to ignore write errors (default: true)
- `sink`: URLs of sub-sinks (can be specified multiple times)

Example:
```go
sflogger.WithSinkURL("multi://?failSafe=true&sink=console://&sink=file:///logs/app.log")
```

Helper function:
```go
sflogger.WithMultiSink([]string{
    "console://?color=true",
    "file:///logs/app.log",
}, true)
```

## Custom Sinks

You can extend the logger by registering custom sinks:

```go
package main

import (
    "net/url"
    "git.snappfood.ir/backend/go/packages/sf-logger"
)

// MyCustomSink implements the Sink interface
type MyCustomSink struct {
    apiKey string
    host   string
}

// Write sends a log entry to the sink
func (s *MyCustomSink) Write(entry map[string]interface{}) error {
    // Your implementation to send logs to a custom destination
    return nil
}

// Sync flushes any buffered log entries
func (s *MyCustomSink) Sync() error {
    // Your implementation
    return nil
}

// Close cleans up resources
func (s *MyCustomSink) Close() error {
    // Your implementation
    return nil
}

// Register the custom sink
func init() {
    sflogger.RegisterSink("myservice", func(u *url.URL) (sflogger.Sink, error) {
        // Parse URL parameters
        q := u.Query()
        apiKey := q.Get("apiKey")
        
        // Create and return the custom sink
        return &MyCustomSink{
            apiKey: apiKey,
            host:   u.Host,
        }, nil
    })
}

func main() {
    // Use the custom sink
    logger := sflogger.New(
        sflogger.WithAppName("MyApp"),
        sflogger.WithSinkURL("myservice://api.example.com?apiKey=secret"),
    )
    
    logger.Info("Using custom sink", nil)
}
```

## Logger Types

SF-Logger provides two logger implementations:

### Zap Logger (Default)

```go
logger := sflogger.New(
    sflogger.WithLoggerType(sflogger.ZapLoggerType),
)
```

Uber's Zap logger is high-performance and optimized for production use.

### Zero Logger

```go
logger := sflogger.New(
    sflogger.WithLoggerType(sflogger.ZeroLoggerType),
)
```

Zerolog is a fast, JSON-first logger with a slightly different API design.

## Examples

### Basic Application Logging

```go
package main

import (
    "git.snappfood.ir/backend/go/packages/sf-logger"
)

func main() {
    // Initialize logger with URL-based sinks
    logger := sflogger.New(
        sflogger.WithAppName("MyApp"),
        sflogger.WithSinkURL("console://"),
        sflogger.WithSinkURL("file:///logs/app.log?maxSize=10&maxAge=30&maxBackups=5&compress=true"),
    )
    
    // Startup logs
    logger.Info("Application starting", map[string]interface{}{
        "version": "1.0.0",
        "environment": "production",
    })
    
    // Business logic...
    
    // Logging errors
    if err := someOperation(); err != nil {
        logger.Error("Operation failed", map[string]interface{}{
            "error": err.Error(),
            "operation": "someOperation",
        })
    }
    
    // Shutdown logs
    logger.Info("Application shutting down", nil)
}
```

### Web Service with Multiple Sinks

```go
package main

import (
    "context"
    "net/http"
    "git.snappfood.ir/backend/go/packages/sf-logger"
)

var logger sflogger.Logger

func init() {
    // Initialize global logger with multiple sinks
    sflogger.InitGlobalLogger(
        sflogger.WithAppName("API-Service"),
        sflogger.WithSinkURL("multi://?failSafe=true" +
            "&sink=console://?color=true" +
            "&sink=file:///logs/api.log?maxSize=10&maxAge=30&maxBackups=5&compress=true" +
            "&sink=elasticsearch://elasticsearch:9200?index=api-logs"),
    )
    
    logger = sflogger.GetGlobalLogger()
}

func loggerMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Add request ID to context
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        // Create context with request info
        ctx := sflogger.WithRequestID(r.Context(), requestID)
        ctx = sflogger.WithUserID(ctx, getUserIDFromRequest(r))
        
        // Log the request
        logger.InfoContext(ctx, "Received request", map[string]interface{}{
            "method": r.Method,
            "path": r.URL.Path,
            "client_ip": r.RemoteAddr,
        })
        
        // Proceed with the request
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func main() {
    // Set up HTTP routes
    http.Handle("/api/", loggerMiddleware(apiHandler()))
    
    // Start server
    logger.InfoWithCategory("General", "Startup", "Server starting", 
        map[string]interface{}{
            "hostIp": "0.0.0.0",
            "port": "8080",
        })
        
    http.ListenAndServe(":8080", nil)
}
```

### Custom Sink with Special Behavior

```go
package main

import (
    "encoding/json"
    "net/url"
    "os"
    "sync"
    "git.snappfood.ir/backend/go/packages/sf-logger"
)

// LevelFilterSink writes only logs of specific levels
type LevelFilterSink struct {
    levels    map[string]bool
    underlying sflogger.Sink
    mutex     sync.Mutex
}

// NewLevelFilterSink creates a new level filter sink
func NewLevelFilterSink(levels []string, underlying sflogger.Sink) *LevelFilterSink {
    levelMap := make(map[string]bool)
    for _, level := range levels {
        levelMap[level] = true
    }
    
    return &LevelFilterSink{
        levels:    levelMap,
        underlying: underlying,
    }
}

// Write sends a log entry to the underlying sink if level matches
func (s *LevelFilterSink) Write(entry map[string]interface{}) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    // Check if this level should be logged
    if level, ok := entry["level"].(string); ok {
        if s.levels[level] {
            return s.underlying.Write(entry)
        }
    }
    
    // Skip this entry
    return nil
}

// Sync flushes the underlying sink
func (s *LevelFilterSink) Sync() error {
    return s.underlying.Sync()
}

// Close closes the underlying sink
func (s *LevelFilterSink) Close() error {
    return s.underlying.Close()
}

// Register the level filter sink factory
func init() {
    sflogger.RegisterSink("levelfilter", func(u *url.URL) (sflogger.Sink, error) {
        // Parse query parameters
        q := u.Query()
        
        // Get levels to include
        levels := q["level"]
        if len(levels) == 0 {
            levels = []string{"ERROR", "FATAL"} // Default to errors only
        }
        
        // Get the underlying sink URL
        targetURL := q.Get("target")
        if targetURL == "" {
            // Default to console
            targetURL = "console://"
        }
        
        // Create the underlying sink
        underlying, err := sflogger.GetSink(targetURL)
        if err != nil {
            return nil, err
        }
        
        // Create and return the filter sink
        return NewLevelFilterSink(levels, underlying), nil
    })
}

func main() {
    // Create a logger that only logs errors to file, but everything to console
    logger := sflogger.New(
        sflogger.WithAppName("FilteredApp"),
        sflogger.WithSinkURL("console://"),
        sflogger.WithSinkURL("levelfilter://?level=ERROR&level=FATAL&target=file:///logs/errors.log"),
    )
    
    // These will go to both console and file
    logger.Error("This is an error", nil)
    logger.Fatal("This is fatal", nil)
    
    // These will only go to console
    logger.Info("This is info", nil)
    logger.Debug("This is debug", nil)
} 