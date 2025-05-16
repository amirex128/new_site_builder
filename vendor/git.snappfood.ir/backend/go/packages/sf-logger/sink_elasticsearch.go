package sflogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// ElasticsearchSinkImpl implements Sink for Elasticsearch
type ElasticsearchSinkImpl struct {
	client    *http.Client
	baseURL   string
	indexName string
	batchSize int
	flushTime time.Duration
	username  string
	password  string
	buffer    []map[string]interface{}
	mutex     sync.Mutex
	ticker    *time.Ticker
	done      chan struct{}
}

// ElasticsearchSinkOption defines the function signature for Elasticsearch sink options
type ElasticsearchSinkOption func(*elasticsearchSinkConfig)

// elasticsearchSinkConfig contains options for configuring an Elasticsearch sink
type elasticsearchSinkConfig struct {
	BaseURL   string
	IndexName string
	BatchSize int
	FlushTime time.Duration
	Username  string
	Password  string
	Timeout   time.Duration
}

// NewElasticsearchSink creates a new Elasticsearch sink
func NewElasticsearchSink(opts ...ElasticsearchSinkOption) Sink {
	config := &elasticsearchSinkConfig{
		IndexName: "logs",
		BatchSize: 100,
		FlushTime: 5 * time.Second,
		Timeout:   10 * time.Second,
	}

	// Apply options
	for _, opt := range opts {
		opt(config)
	}

	// Create HTTP client
	client := &http.Client{
		Timeout: config.Timeout,
	}

	// Create sink
	sink := &ElasticsearchSinkImpl{
		client:    client,
		baseURL:   config.BaseURL,
		indexName: config.IndexName,
		batchSize: config.BatchSize,
		flushTime: config.FlushTime,
		username:  config.Username,
		password:  config.Password,
		buffer:    make([]map[string]interface{}, 0, config.BatchSize),
		done:      make(chan struct{}),
	}

	// Start background flusher
	sink.ticker = time.NewTicker(config.FlushTime)
	go sink.flushPeriodically()

	return sink
}

// Write sends a log entry to Elasticsearch
func (s *ElasticsearchSinkImpl) Write(entry map[string]interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Add entry to buffer
	s.buffer = append(s.buffer, entry)

	// Flush if buffer is full
	if len(s.buffer) >= s.batchSize {
		return s.flush()
	}

	return nil
}

// flush sends buffered entries to Elasticsearch
func (s *ElasticsearchSinkImpl) flush() error {
	if len(s.buffer) == 0 {
		return nil
	}

	// Create bulk request
	var bulkData bytes.Buffer
	for _, entry := range s.buffer {
		// Add index action line
		indexAction := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": s.indexName,
			},
		}
		actionLine, err := json.Marshal(indexAction)
		if err != nil {
			continue
		}
		bulkData.Write(actionLine)
		bulkData.WriteString("\n")

		// Add entry data
		entryLine, err := json.Marshal(entry)
		if err != nil {
			continue
		}
		bulkData.Write(entryLine)
		bulkData.WriteString("\n")
	}

	// Send bulk request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/_bulk", s.baseURL, s.indexName), &bulkData)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-ndjson")
	if s.username != "" && s.password != "" {
		req.SetBasicAuth(s.username, s.password)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Reset buffer
	s.buffer = s.buffer[:0]

	if resp.StatusCode >= 400 {
		return fmt.Errorf("elasticsearch returned status code %d", resp.StatusCode)
	}

	return nil
}

// flushPeriodically flushes the buffer at regular intervals
func (s *ElasticsearchSinkImpl) flushPeriodically() {
	for {
		select {
		case <-s.ticker.C:
			s.mutex.Lock()
			_ = s.flush()
			s.mutex.Unlock()
		case <-s.done:
			return
		}
	}
}

// Close flushes any remaining entries and cleans up resources
func (s *ElasticsearchSinkImpl) Close() error {
	s.ticker.Stop()
	close(s.done)

	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.flush()
}

// Sync flushes any buffered entries
func (s *ElasticsearchSinkImpl) Sync() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.flush()
}

// ElasticsearchWithURL sets the Elasticsearch server URL
func ElasticsearchWithURL(url string) ElasticsearchSinkOption {
	return func(c *elasticsearchSinkConfig) {
		c.BaseURL = url
	}
}

// ElasticsearchWithIndexName sets the index name
func ElasticsearchWithIndexName(indexName string) ElasticsearchSinkOption {
	return func(c *elasticsearchSinkConfig) {
		c.IndexName = indexName
	}
}

// ElasticsearchWithBatchSize sets the batch size for bulk operations
func ElasticsearchWithBatchSize(size int) ElasticsearchSinkOption {
	return func(c *elasticsearchSinkConfig) {
		c.BatchSize = size
	}
}

// ElasticsearchWithFlushTime sets the flush interval
func ElasticsearchWithFlushTime(duration time.Duration) ElasticsearchSinkOption {
	return func(c *elasticsearchSinkConfig) {
		c.FlushTime = duration
	}
}

// ElasticsearchWithCredentials sets the username and password
func ElasticsearchWithCredentials(username, password string) ElasticsearchSinkOption {
	return func(c *elasticsearchSinkConfig) {
		c.Username = username
		c.Password = password
	}
}

// ElasticsearchWithTimeout sets the HTTP timeout
func ElasticsearchWithTimeout(timeout time.Duration) ElasticsearchSinkOption {
	return func(c *elasticsearchSinkConfig) {
		c.Timeout = timeout
	}
}

// Register Elasticsearch sink with the registry
func init() {
	RegisterSink("elasticsearch", func(u *url.URL) (Sink, error) {
		var opts []ElasticsearchSinkOption

		// Build base URL
		scheme := "http"
		if u.Scheme == "elasticsearch+https" {
			scheme = "https"
		}
		baseURL := fmt.Sprintf("%s://%s", scheme, u.Host)
		opts = append(opts, ElasticsearchWithURL(baseURL))

		// Parse query parameters
		q := u.Query()

		if indexName := q.Get("index"); indexName != "" {
			opts = append(opts, ElasticsearchWithIndexName(indexName))
		}

		if batchSize := q.Get("batchSize"); batchSize != "" {
			if size, err := strconv.Atoi(batchSize); err == nil && size > 0 {
				opts = append(opts, ElasticsearchWithBatchSize(size))
			}
		}

		if flushTime := q.Get("flushTime"); flushTime != "" {
			if duration, err := time.ParseDuration(flushTime); err == nil {
				opts = append(opts, ElasticsearchWithFlushTime(duration))
			}
		}

		if timeout := q.Get("timeout"); timeout != "" {
			if duration, err := time.ParseDuration(timeout); err == nil {
				opts = append(opts, ElasticsearchWithTimeout(duration))
			}
		}

		// Check for credentials
		if username := q.Get("username"); username != "" {
			password := q.Get("password")
			opts = append(opts, ElasticsearchWithCredentials(username, password))
		}

		return NewElasticsearchSink(opts...), nil
	})
}
