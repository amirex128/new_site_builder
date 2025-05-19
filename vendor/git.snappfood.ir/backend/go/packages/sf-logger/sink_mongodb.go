package sflogger

import (
	"context"
	"fmt"
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
}

// MongoDBSinkOption defines the function signature for MongoDB sink options
type MongoDBSinkOption func(*mongoDBSinkConfig)

// mongoDBSinkConfig contains options for configuring a MongoDB sink
type mongoDBSinkConfig struct {
	URI        string
	Database   string
	Collection string
	BatchSize  int
	FlushTime  time.Duration
	Timeout    time.Duration
}

// NewMongoDBSink creates a new MongoDB sink
func NewMongoDBSink(opts ...MongoDBSinkOption) Sink {
	config := &mongoDBSinkConfig{
		Database:   "logs",
		Collection: "logs",
		BatchSize:  100,
		FlushTime:  5 * time.Second,
		Timeout:    10 * time.Second,
	}

	// Apply options
	for _, opt := range opts {
		opt(config)
	}

	// Create context with timeout for MongoDB operations
	ctx, cancel := context.WithCancel(context.Background())

	// Create MongoDB client
	clientOpts := options.Client().ApplyURI(config.URI).SetTimeout(config.Timeout)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		// In a real implementation, we might want to handle this error better
		// For now, we'll just panic
		panic(fmt.Sprintf("failed to connect to MongoDB: %v", err))
	}

	// Ping MongoDB to verify connection
	pingCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		client.Disconnect(ctx)
		panic(fmt.Sprintf("failed to ping MongoDB: %v", err))
	}

	// Get collection reference
	collection := client.Database(config.Database).Collection(config.Collection)

	// Create sink
	sink := &MongoDBSinkImpl{
		client:        client,
		database:      config.Database,
		collection:    config.Collection,
		batchSize:     config.BatchSize,
		flushTime:     config.FlushTime,
		buffer:        make([]map[string]interface{}, 0, config.BatchSize),
		done:          make(chan struct{}),
		collectionRef: collection,
		ctx:           ctx,
		cancelFunc:    cancel,
	}

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
		return fmt.Errorf("failed to insert documents into MongoDB: %w", err)
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

	// Disconnect from MongoDB
	if s.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if disconnectErr := s.client.Disconnect(ctx); disconnectErr != nil {
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

// Register MongoDB sink with the registry
func init() {
	RegisterSink("mongodb", func(u *url.URL) (Sink, error) {
		var opts []MongoDBSinkOption

		// Build connection URI
		// MongoDB connection strings are in the format: mongodb://user:password@host:port/database
		// We need to reconstruct this from the parsed URL
		mongoURI := fmt.Sprintf("mongodb://%s", u.Host)

		// Handle authentication if present in the URL
		if u.User != nil {
			username := u.User.Username()
			password, _ := u.User.Password()
			// If we have auth info, reconstruct the URI with auth
			mongoURI = fmt.Sprintf("mongodb://%s:%s@%s", username, password, u.Host)
		}

		if u.Path != "" {
			// Remove leading slash from path
			dbName := u.Path
			if len(dbName) > 0 && dbName[0] == '/' {
				dbName = dbName[1:]
			}
			if dbName != "" {
				opts = append(opts, MongoDBWithDatabase(dbName))
			}
		}
		opts = append(opts, MongoDBWithURI(mongoURI))

		// Parse query parameters
		q := u.Query()

		if collection := q.Get("collection"); collection != "" {
			opts = append(opts, MongoDBWithCollection(collection))
		}

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

		return NewMongoDBSink(opts...), nil
	})
}
