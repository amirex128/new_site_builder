package sflogger

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"sync"
)

// ConsoleSinkImpl outputs logs to stdout or stderr
type ConsoleSinkImpl struct {
	writer     io.Writer
	encoder    *json.Encoder
	colorized  bool
	timeFormat string
	mutex      sync.Mutex
}

// ConsoleSinkOption defines the function signature for console sink options
type ConsoleSinkOption func(*consoleSinkConfig)

// consoleSinkConfig contains options for configuring a console sink
type consoleSinkConfig struct {
	UseStderr  bool
	Colorized  bool
	TimeFormat string
}

// NewConsoleSink creates a new console sink
func NewConsoleSink(opts ...ConsoleSinkOption) Sink {
	config := &consoleSinkConfig{
		UseStderr:  false,
		Colorized:  false,
		TimeFormat: "2006-01-02T15:04:05.000Z07:00",
	}

	// Apply options
	for _, opt := range opts {
		opt(config)
	}

	var writer io.Writer
	if config.UseStderr {
		writer = os.Stderr
	} else {
		writer = os.Stdout
	}

	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "")

	return &ConsoleSinkImpl{
		writer:     writer,
		encoder:    encoder,
		colorized:  config.Colorized,
		timeFormat: config.TimeFormat,
	}
}

// Write sends a log entry to the console
func (s *ConsoleSinkImpl) Write(entry map[string]interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Format timestamp if present
	if ts, ok := entry["timestamp"].(string); ok {
		// Here you could format the timestamp according to s.timeFormat
		entry["timestamp"] = ts
	}

	return s.encoder.Encode(entry)
}

// Close is a no-op for console sink
func (s *ConsoleSinkImpl) Close() error {
	return nil
}

// Sync is a no-op for console sink
func (s *ConsoleSinkImpl) Sync() error {
	// Stdout/Stderr doesn't need explicit syncing
	return nil
}

// ConsoleWithStderr configures the console sink to use stderr
func ConsoleWithStderr() ConsoleSinkOption {
	return func(c *consoleSinkConfig) {
		c.UseStderr = true
	}
}

// ConsoleWithColor configures the console sink to use colors
func ConsoleWithColor() ConsoleSinkOption {
	return func(c *consoleSinkConfig) {
		c.Colorized = true
	}
}

// ConsoleWithTimeFormat configures the time format for the console sink
func ConsoleWithTimeFormat(format string) ConsoleSinkOption {
	return func(c *consoleSinkConfig) {
		c.TimeFormat = format
	}
}

// Register console sink with the registry
func init() {
	RegisterSink("console", func(u *url.URL) (Sink, error) {
		var opts []ConsoleSinkOption

		// Parse query parameters
		q := u.Query()

		if q.Get("stderr") == "true" {
			opts = append(opts, ConsoleWithStderr())
		}

		if q.Get("color") == "true" {
			opts = append(opts, ConsoleWithColor())
		}

		if tf := q.Get("timeFormat"); tf != "" {
			opts = append(opts, ConsoleWithTimeFormat(tf))
		}

		return NewConsoleSink(opts...), nil
	})
}
