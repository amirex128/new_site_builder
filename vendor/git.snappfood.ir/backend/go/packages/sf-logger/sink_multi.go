package sflogger

import (
	"fmt"
	"net/url"
	"sync"
)

// MultiSinkImpl sends logs to multiple destinations with optional filtering
type MultiSinkImpl struct {
	sinks     []Sink
	filters   []func(map[string]interface{}) bool
	mutex     sync.RWMutex
	failSafe  bool
	routeFunc func(map[string]interface{}) int
}

// MultiSinkOption defines the function signature for multi-sink options
type MultiSinkOption func(*multiSinkConfig)

// multiSinkConfig contains options for configuring a multi-sink
type multiSinkConfig struct {
	Sinks     []Sink
	Filters   []func(map[string]interface{}) bool
	FailSafe  bool
	RouteFunc func(map[string]interface{}) int
}

// NewMultiSink creates a new multi-sink
func NewMultiSink(opts ...MultiSinkOption) Sink {
	config := &multiSinkConfig{
		Sinks:    make([]Sink, 0),
		Filters:  make([]func(map[string]interface{}) bool, 0),
		FailSafe: true,
	}

	// Apply options
	for _, opt := range opts {
		opt(config)
	}

	// Create a default route function if none provided
	routeFunc := config.RouteFunc
	if routeFunc == nil {
		// Default routing sends to all sinks
		routeFunc = func(entry map[string]interface{}) int {
			return -1 // -1 means send to all sinks
		}
	}

	return &MultiSinkImpl{
		sinks:     config.Sinks,
		filters:   config.Filters,
		failSafe:  config.FailSafe,
		routeFunc: routeFunc,
	}
}

// Write sends a log entry to multiple sinks based on routing
func (s *MultiSinkImpl) Write(entry map[string]interface{}) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if len(s.sinks) == 0 {
		return nil
	}

	// Apply filters
	for _, filter := range s.filters {
		if !filter(entry) {
			return nil
		}
	}

	// Determine which sink(s) to use
	sinkIndex := s.routeFunc(entry)

	var errs []error
	if sinkIndex < 0 || sinkIndex >= len(s.sinks) {
		// Send to all sinks
		for _, sink := range s.sinks {
			if err := sink.Write(entry); err != nil && !s.failSafe {
				errs = append(errs, err)
			}
		}
	} else {
		// Send to specific sink
		if err := s.sinks[sinkIndex].Write(entry); err != nil && !s.failSafe {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to write to %d sinks", len(errs))
	}
	return nil
}

// Close closes all sinks
func (s *MultiSinkImpl) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var errs []error
	for _, sink := range s.sinks {
		if err := sink.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to close %d sinks", len(errs))
	}
	return nil
}

// Sync flushes all sinks
func (s *MultiSinkImpl) Sync() error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var errs []error
	for _, sink := range s.sinks {
		if err := sink.Sync(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to sync %d sinks", len(errs))
	}
	return nil
}

// AddSink adds a sink to the multi-sink
func (s *MultiSinkImpl) AddSink(sink Sink) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sinks = append(s.sinks, sink)
}

// RemoveSink removes a sink from the multi-sink
func (s *MultiSinkImpl) RemoveSink(sink Sink) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i, existingSink := range s.sinks {
		if existingSink == sink {
			s.sinks = append(s.sinks[:i], s.sinks[i+1:]...)
			break
		}
	}
}

// MultiWithSinks adds sinks to the multi-sink
func MultiWithSinks(sinks ...Sink) MultiSinkOption {
	return func(c *multiSinkConfig) {
		c.Sinks = append(c.Sinks, sinks...)
	}
}

// MultiWithFilter adds a filter to the multi-sink
func MultiWithFilter(filter func(map[string]interface{}) bool) MultiSinkOption {
	return func(c *multiSinkConfig) {
		c.Filters = append(c.Filters, filter)
	}
}

// MultiWithFailSafe sets whether the sink should ignore write errors
func MultiWithFailSafe(failSafe bool) MultiSinkOption {
	return func(c *multiSinkConfig) {
		c.FailSafe = failSafe
	}
}

// MultiWithRouter sets a custom routing function
func MultiWithRouter(routeFunc func(map[string]interface{}) int) MultiSinkOption {
	return func(c *multiSinkConfig) {
		c.RouteFunc = routeFunc
	}
}

// Level-based routing function
func LevelBasedRouter(levelMap map[string]int) func(map[string]interface{}) int {
	return func(entry map[string]interface{}) int {
		if level, ok := entry["level"]; ok {
			if sinkIndex, ok := levelMap[fmt.Sprintf("%v", level)]; ok {
				return sinkIndex
			}
		}
		return -1 // default to all sinks
	}
}

// Category-based routing function
func CategoryBasedRouter(categoryMap map[string]int) func(map[string]interface{}) int {
	return func(entry map[string]interface{}) int {
		if category, ok := entry["category"]; ok {
			if sinkIndex, ok := categoryMap[fmt.Sprintf("%v", category)]; ok {
				return sinkIndex
			}
		}
		return -1 // default to all sinks
	}
}

// Register multi-sink with the registry
func init() {
	RegisterSink("multi", func(u *url.URL) (Sink, error) {
		var opts []MultiSinkOption

		// Parse query parameters
		q := u.Query()

		// Add sinks from URLs
		sinkURLs := q["sink"]
		for _, sinkURL := range sinkURLs {
			sink, err := GetSink(sinkURL)
			if err == nil && sink != nil {
				opts = append(opts, MultiWithSinks(sink))
			}
		}

		// Set fail-safe option
		if failSafe := q.Get("failSafe"); failSafe != "" {
			opts = append(opts, MultiWithFailSafe(failSafe == "true"))
		}

		return NewMultiSink(opts...), nil
	})
}
