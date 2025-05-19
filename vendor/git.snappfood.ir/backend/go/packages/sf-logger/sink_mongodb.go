package sflogger

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBSinkImpl implements Sink for MongoDB
type MongoDBSinkImpl struct {
	client        *mongo.Client
	database      string
	collection    string
	batchSize     int
	flushTime     time.Duration
	buffer        []map[string]interface{}
	mutex         sync.Mutex
	ticker        *time.Ticker
	done          chan struct{}
	collectionRef *mongo.Collection
	ctx           context.Context
	cancelFunc    context.CancelFunc
	connected     bool
	failSilently  bool
}

// MongoDBSinkOption defines the function signature for MongoDB sink options
type MongoDBSinkOption func(*mongoDBSinkConfig)

// mongoDBSinkConfig contains options for configuring a MongoDB sink
type mongoDBSinkConfig struct {
	URI          string
	Database     string
	Collection   string
	BatchSize    int
	FlushTime    time.Duration
	Timeout      time.Duration
	FailSilently bool
}

// NewMongoDBSink creates a new MongoDB sink
func NewMongoDBSink(opts ...MongoDBSinkOption) Sink {
	config := &mongoDBSinkConfig{
		Database:     "logs",
		Collection:   "logs",
		BatchSize:    100,
		FlushTime:    5 * time.Second,
		Timeout:      10 * time.Second,
		FailSilently: true,
	}

	// Apply options
	for _, opt := range opts {
		opt(config)
	}

	// Create context with timeout for MongoDB operations
	ctx, cancel := context.WithCancel(context.Background())

	// Create sink instance first
	sink := &MongoDBSinkImpl{
		database:     config.Database,
		collection:   config.Collection,
		batchSize:    config.BatchSize,
		flushTime:    config.FlushTime,
		buffer:       make([]map[string]interface{}, 0, config.BatchSize),
		done:         make(chan struct{}),
		ctx:          ctx,
		cancelFunc:   cancel,
		connected:    false,
		failSilently: config.FailSilently,
	}

	// Try to connect to MongoDB
	clientOpts := options.Client().ApplyURI(config.URI).SetTimeout(config.Timeout)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		if !config.FailSilently {
			log.Printf("Failed to connect to MongoDB: %v", err)
		}
		// Start the ticker anyway so we don't block if MongoDB becomes available later
		sink.ticker = time.NewTicker(config.FlushTime)
		go sink.flushPeriodically()
		return sink
	}

	// Ping MongoDB to verify connection
	pingCtx, pingCancel := context.WithTimeout(ctx, config.Timeout)
	defer pingCancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		if !config.FailSilently {
			log.Printf("Failed to ping MongoDB: %v", err)
		}
		// Start the ticker anyway so we don't block if MongoDB becomes available later
		sink.ticker = time.NewTicker(config.FlushTime)
		go sink.flushPeriodically()
		return sink
	}

	// Get collection reference
	sink.client = client
	sink.collectionRef = client.Database(config.Database).Collection(config.Collection)
	sink.connected = true

	// Start background flusher
	sink.ticker = time.NewTicker(config.FlushTime)
	go sink.flushPeriodically()

	return sink
}

// Write sends a log entry to MongoDB
func (s *MongoDBSinkImpl) Write(entry map[string]interface{}) error {
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

// flush sends buffered entries to MongoDB
func (s *MongoDBSinkImpl) flush() error {
	if len(s.buffer) == 0 {
		return nil
	}

	// If not connected, just clear the buffer and return
	if !s.connected || s.client == nil || s.collectionRef == nil {
		s.buffer = s.buffer[:0]
		return nil
	}

	// Convert buffer to interface slice for MongoDB
	documents := make([]interface{}, len(s.buffer))
	for i, entry := range s.buffer {
		// Add timestamp if not present
		if _, ok := entry["timestamp"]; !ok {
			entry["timestamp"] = time.Now()
		}
		documents[i] = entry
	}

	// Create context with timeout for insert operation
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	// Insert documents
	_, err := s.collectionRef.InsertMany(ctx, documents)
	if err != nil {
		if !s.failSilently {
			log.Printf("Failed to insert documents into MongoDB: %v", err)
		}
		// Clear buffer anyway to prevent memory buildup
		s.buffer = s.buffer[:0]
		return err
	}

	// Reset buffer
	s.buffer = s.buffer[:0]

	return nil
}

// flushPeriodically flushes the buffer at regular intervals
func (s *MongoDBSinkImpl) flushPeriodically() {
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
func (s *MongoDBSinkImpl) Close() error {
	s.ticker.Stop()
	close(s.done)

	s.mutex.Lock()
	err := s.flush()
	s.mutex.Unlock()

	// Disconnect from MongoDB if connected
	if s.connected && s.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if disconnectErr := s.client.Disconnect(ctx); disconnectErr != nil && !s.failSilently {
			log.Printf("Failed to disconnect from MongoDB: %v", disconnectErr)
			if err == nil {
				err = disconnectErr
			}
		}
	}

	// Cancel the context
	s.cancelFunc()

	return err
}

// Sync flushes any buffered entries
func (s *MongoDBSinkImpl) Sync() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.flush()
}

// MongoDBWithURI sets the MongoDB connection URI
func MongoDBWithURI(uri string) MongoDBSinkOption {
	return func(c *mongoDBSinkConfig) {
		c.URI = uri
	}
}

// MongoDBWithDatabase sets the database name
func MongoDBWithDatabase(database string) MongoDBSinkOption {
	return func(c *mongoDBSinkConfig) {
		c.Database = database
	}
}

// MongoDBWithCollection sets the collection name
func MongoDBWithCollection(collection string) MongoDBSinkOption {
	return func(c *mongoDBSinkConfig) {
		c.Collection = collection
	}
}

// MongoDBWithBatchSize sets the batch size for bulk operations
func MongoDBWithBatchSize(size int) MongoDBSinkOption {
	return func(c *mongoDBSinkConfig) {
		c.BatchSize = size
	}
}

// MongoDBWithFlushTime sets the flush interval
func MongoDBWithFlushTime(duration time.Duration) MongoDBSinkOption {
	return func(c *mongoDBSinkConfig) {
		c.FlushTime = duration
	}
}

// MongoDBWithTimeout sets the connection timeout
func MongoDBWithTimeout(timeout time.Duration) MongoDBSinkOption {
	return func(c *mongoDBSinkConfig) {
		c.Timeout = timeout
	}
}

// MongoDBWithFailSilently sets whether to fail silently or log errors
func MongoDBWithFailSilently(failSilently bool) MongoDBSinkOption {
	return func(c *mongoDBSinkConfig) {
		c.FailSilently = failSilently
	}
}

// Register MongoDB sink with the registry
func init() {
	RegisterSink("mongodb", func(u *url.URL) (Sink, error) {
		var opts []MongoDBSinkOption

		// Extract database and collection from URL path and query
		dbName := ""
		if u.Path != "" {
			// Remove leading slash from path
			dbName = u.Path
			if len(dbName) > 0 && dbName[0] == '/' {
				dbName = dbName[1:]
			}
		}

		// Parse query parameters
		q := u.Query()
		collName := q.Get("collection")
		if collName == "" {
			collName = "logs"
		}

		// Build connection URI with database included
		// Format: mongodb://[username:password@]host[:port][/database][?options]
		var mongoURI string
		if u.User != nil {
			// Has authentication
			username := u.User.Username()
			password, _ := u.User.Password()
			mongoURI = fmt.Sprintf("mongodb://%s:%s@%s", username, password, u.Host)
		} else {
			// No authentication
			mongoURI = fmt.Sprintf("mongodb://%s", u.Host)
		}

		// Add database to URI if provided
		if dbName != "" {
			mongoURI += "/" + dbName
		}

		opts = append(opts, MongoDBWithURI(mongoURI))

		// Set database and collection
		if dbName != "" {
			opts = append(opts, MongoDBWithDatabase(dbName))
		}

		opts = append(opts, MongoDBWithCollection(collName))

		// Parse other query parameters
		if batchSize := q.Get("batchSize"); batchSize != "" {
			if size, err := strconv.Atoi(batchSize); err == nil && size > 0 {
				opts = append(opts, MongoDBWithBatchSize(size))
			}
		}

		if flushTime := q.Get("flushTime"); flushTime != "" {
			if duration, err := time.ParseDuration(flushTime); err == nil {
				opts = append(opts, MongoDBWithFlushTime(duration))
			}
		}

		if timeout := q.Get("timeout"); timeout != "" {
			if duration, err := time.ParseDuration(timeout); err == nil {
				opts = append(opts, MongoDBWithTimeout(duration))
			}
		}

		if failSilently := q.Get("failSilently"); failSilently != "" {
			silent := failSilently == "true" || failSilently == "1"
			opts = append(opts, MongoDBWithFailSilently(silent))
		}

		return NewMongoDBSink(opts...), nil
	})
}
