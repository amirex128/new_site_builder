package sflogger

import (
	"fmt"
	"net/url"
	"sync"
)

// SinkType represents different log destinations
type SinkType string

const (
	// ConsoleSinkType outputs logs to stdout
	ConsoleSinkType SinkType = "console"

	// FileSinkType outputs logs to files with rotation
	FileSinkType SinkType = "file"

	// GraylogSinkType sends logs to Graylog via GELF
	GraylogSinkType SinkType = "graylog"

	// LogstashSinkType sends logs to Logstash
	LogstashSinkType SinkType = "logstash"

	// LokiSinkType sends logs to Grafana Loki
	LokiSinkType SinkType = "loki"

	// ElasticsearchSinkType sends logs to Elasticsearch
	ElasticsearchSinkType SinkType = "elasticsearch"
)

// Protocol defines the communication protocol for remote sinks
type Protocol string

// String returns the string representation of the protocol
func (p Protocol) String() string {
	return string(p)
}

const (
	// UDP protocol
	UDP Protocol = "udp"

	// TCP protocol
	TCP Protocol = "tcp"

	// HTTP protocol
	HTTP Protocol = "http"

	// HTTPS protocol
	HTTPS Protocol = "https"
)

// SinkFactory is a function that creates a sink from a URL
type SinkFactory func(*url.URL) (Sink, error)

var (
	registry      = make(map[string]SinkFactory)
	registryMutex sync.RWMutex
)

// RegisterSink registers a sink factory for a particular URL scheme
func RegisterSink(scheme string, factory SinkFactory) error {
	registryMutex.Lock()
	defer registryMutex.Unlock()

	if scheme == "" {
		return fmt.Errorf("sink scheme cannot be empty")
	}

	if factory == nil {
		return fmt.Errorf("sink factory cannot be nil")
	}

	if _, exists := registry[scheme]; exists {
		return fmt.Errorf("sink scheme '%s' already registered", scheme)
	}

	registry[scheme] = factory
	return nil
}

// GetSink retrieves a sink by URL
func GetSink(sinkURL string) (Sink, error) {
	parsedURL, err := url.Parse(sinkURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sink URL: %w", err)
	}

	registryMutex.RLock()
	factory, exists := registry[parsedURL.Scheme]
	registryMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no sink registered for scheme '%s'", parsedURL.Scheme)
	}

	return factory(parsedURL)
}

// UnregisterSink removes a sink from the registry
func UnregisterSink(scheme string) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	delete(registry, scheme)
}
