package sflogger

import "fmt"

// Option is a function that configures the logger
type Option func(*Config)

// WithLoggerType sets the logger implementation type
func WithLoggerType(loggerType LoggerType) Option {
	return func(c *Config) {
		c.LoggerType = loggerType
	}
}

// WithLevel sets the minimum log level
func WithLevel(level Level) Option {
	return func(c *Config) {
		c.Level = level
	}
}

// WithAppName sets the application name
func WithAppName(appName string) Option {
	return func(c *Config) {
		c.AppName = appName
	}
}

// WithFilePath is deprecated: use WithFileSink instead which provides more comprehensive file configuration
// Kept for backwards compatibility
func WithFilePath(filePath string) Option {
	return func(c *Config) {
		c.FilePath = filePath
	}
}

// WithDevelopment enables development mode
func WithDevelopment(enable bool) Option {
	return func(c *Config) {
		c.Development = enable
	}
}

// WithConsole is deprecated: use WithConsoleSink instead which provides console sink configuration
// Kept for backwards compatibility
func WithConsole(enable bool) Option {
	return func(c *Config) {
		c.EnableConsole = enable
	}
}

// WithStacktrace enables stack traces for errors
func WithStacktrace(enable bool) Option {
	return func(c *Config) {
		c.EnableStacktrace = enable
	}
}

// WithFormatter sets the log formatter
func WithFormatter(formatter FormatterType) Option {
	return func(c *Config) {
		c.Formatter = formatter
	}
}

// WithMaxSize sets the maximum size in megabytes of the log file before it gets rotated
func WithMaxSize(size int) Option {
	return func(c *Config) {
		c.MaxSize = size
	}
}

// WithMaxAge sets the maximum number of days to retain old log files
func WithMaxAge(age int) Option {
	return func(c *Config) {
		c.MaxAge = age
	}
}

// WithMaxBackups sets the maximum number of old log files to retain
func WithMaxBackups(backups int) Option {
	return func(c *Config) {
		c.MaxBackups = backups
	}
}

// WithCompression enables log file compression
func WithCompression(enable bool) Option {
	return func(c *Config) {
		c.Compress = enable
	}
}

// WithTimeFormat sets the time format string
func WithTimeFormat(format string) Option {
	return func(c *Config) {
		c.TimeFormat = format
	}
}

// WithSinkURL adds a sink URL to the logger configuration
func WithSinkURL(sinkURL string) Option {
	return func(c *Config) {
		c.SinkURLs = append(c.SinkURLs, sinkURL)
	}
}

// WithSinkURLs replaces all sink URLs in the logger configuration
func WithSinkURLs(sinkURLs []string) Option {
	return func(c *Config) {
		c.SinkURLs = sinkURLs
	}
}

// WithConsoleSink adds a console sink to the logger
func WithConsoleSink(color bool) Option {
	return func(c *Config) {
		url := "console://"
		if color {
			url += "?color=true"
		}
		c.SinkURLs = append(c.SinkURLs, url)
	}
}

// WithFileSink adds a file sink to the logger
func WithFileSink(filePath string, maxSize, maxAge, maxBackups int, compress bool) Option {
	return func(c *Config) {
		c.FilePath = filePath
		c.MaxSize = maxSize
		c.MaxAge = maxAge
		c.MaxBackups = maxBackups
		c.Compress = compress

		url := fmt.Sprintf("file://%s?maxSize=%d&maxAge=%d&maxBackups=%d&compress=%t",
			filePath, maxSize, maxAge, maxBackups, compress)
		c.SinkURLs = append(c.SinkURLs, url)
	}
}

// WithGraylogSink adds a Graylog sink to the logger
func WithGraylogSink(host string, port int, useUDP bool) Option {
	return func(c *Config) {
		proto := "graylog+tcp"
		if useUDP {
			proto = "graylog"
		}
		url := fmt.Sprintf("%s://%s:%d", proto, host, port)
		c.SinkURLs = append(c.SinkURLs, url)
	}
}

// WithLogstashSink adds a Logstash sink to the logger
func WithLogstashSink(host string, port int, useUDP bool) Option {
	return func(c *Config) {
		proto := "logstash"
		if !useUDP {
			proto = "logstash+tcp"
		}
		url := fmt.Sprintf("%s://%s:%d", proto, host, port)
		c.SinkURLs = append(c.SinkURLs, url)
	}
}

// WithLokiSink adds a Loki sink to the logger
func WithLokiSink(url string, apiKey string) Option {
	return func(c *Config) {
		lokiURL := fmt.Sprintf("loki://%s", url)
		if apiKey != "" {
			lokiURL += fmt.Sprintf("?apiKey=%s", apiKey)
		}
		c.SinkURLs = append(c.SinkURLs, lokiURL)
	}
}

// WithElasticsearchSink adds an Elasticsearch sink to the logger
func WithElasticsearchSink(url string, indexName string, username, password string, batchSize int) Option {
	return func(c *Config) {
		esURL := fmt.Sprintf("elasticsearch://%s?index=%s", url, indexName)
		if username != "" {
			esURL += fmt.Sprintf("&username=%s&password=%s", username, password)
		}
		if batchSize > 0 {
			esURL += fmt.Sprintf("&batchSize=%d", batchSize)
		}
		c.SinkURLs = append(c.SinkURLs, esURL)
	}
}

// WithMultiSink adds a multi-sink configuration with multiple destinations
func WithMultiSink(sinkURLs []string, failSafe bool) Option {
	return func(c *Config) {
		multiURL := "multi://?failSafe="
		if failSafe {
			multiURL += "true"
		} else {
			multiURL += "false"
		}

		for _, sinkURL := range sinkURLs {
			multiURL += "&sink=" + sinkURL
		}

		c.SinkURLs = append(c.SinkURLs, multiURL)
	}
}

// WithCustomSink adds a custom sink URL directly
func WithCustomSink(scheme, rawURL string) Option {
	return func(c *Config) {
		sinkURL := fmt.Sprintf("%s://%s", scheme, rawURL)
		c.SinkURLs = append(c.SinkURLs, sinkURL)
	}
}
