# sf-logger

A flexible, pluggable, and structured logging library for Go, built on top of [Uber Zap](https://github.com/uber-go/zap). Supports multiple output sinks (console, file, MongoDB, Elasticsearch), log levels, JSON or colored text formatting, and rich context/categorization.

## Features

- Multiple output sinks: Console, File (with rotation), MongoDB, Elasticsearch
- Structured logging with context and extra fields
- Log levels: Debug, Info, Warn, Error, Fatal
- JSON or colored text formatting
- Category and subcategory support for log organization
- Thread-safe, high-performance

## Installation

```
go get github.com/yourusername/sf-logger
```

## Usage

### Basic Setup

```go
import "github.com/yourusername/sf-logger"

func main() {
    logger := sflogger.RegisterLogger(
        sflogger.WithAppName("my-app"),
        sflogger.WithLevel(sflogger.InfoLevel),
        sflogger.WithFormatter(sflogger.ColoredTextFormatter),
        sflogger.WithStacktrace(true),
    )

    logger.Info("Hello, world!", map[string]interface{}{"foo": "bar"})
}
```

### File Sink Example

```go
logger := sflogger.RegisterLogger(
    sflogger.WithAppName("my-app"),
    sflogger.WithFileSink("logs/app.log", 10, 7, 3, true), // path, maxSizeMB, maxAgeDays, maxBackups, compress
)
```

### MongoDB Sink Example

```go
logger := sflogger.RegisterLogger(
    sflogger.WithAppName("my-app"),
    sflogger.WithMongoDBSink(
        "localhost", 27017, "logsdb", "logs", "user", "pass", 5, false, // host, port, db, collection, user, pass, flushSec, compress
    ),
)
```

### Elasticsearch Sink Example

```go
logger := sflogger.RegisterLogger(
    sflogger.WithAppName("my-app"),
    sflogger.WithElasticSearchSink(
        "http://localhost:9200", // url
        "elastic",               // username
        "changeme",              // password
        "logs-index",            // index name
        5,                        // flushSec
    ),
)
```

### Logging with Categories and Context

```go
logger.InfoWithCategory(
    sflogger.Category.API.HTTP, // category
    sflogger.SubCategory.API.Request, // subcategory
    "Received API request",
    map[string]interface{}{"request_id": "abc123"},
)

logger.ErrorContext(ctx, "Something went wrong", map[string]interface{}{"err": err})
```

## API

See `logger.go` for the full `Logger` interface, including:
- `Debug`, `Info`, `Warn`, `Error`, `Fatal`
- `Debugf`, `Infof`, ... (formatted)
- `DebugContext`, ... (with context)
- `DebugWithCategory`, ... (with category/subcategory)

## Configuration Options

- `WithAppName(name string)`
- `WithLevel(level Level)`
- `WithFormatter(formatter FormatterType)`
- `WithStacktrace(enabled bool)`
- `WithFileSink(path string, maxSizeMB, maxAgeDays, maxBackups int, compress bool)`
- `WithMongoDBSink(host string, port int, database, collection, username, password string, flushSec int, compress bool)`
- `WithElasticSearchSink(url string, username string, password string, index string, flushSec int)`

## License

MIT 