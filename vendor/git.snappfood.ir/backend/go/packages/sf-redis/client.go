package sfredis

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.elastic.co/apm/v2"
)

// =============================================================================
// Types
// =============================================================================

// SfRedis is a client for SfRedis operations
type SfRedis struct {
	client redis.UniversalClient
	ctx    context.Context
	logger Logger
}

// OperationContext holds the context and span for an operation
type OperationContext struct {
	ctx    context.Context
	span   *apm.Span
	logger Logger
}

// =============================================================================
// Factory Functions
// =============================================================================

// MustClient creates a new SfRedis client for the named connection
func MustClient(ctx context.Context, name string) *SfRedis {
	client, err := getConnection(name)
	var logger Logger
	globalRegistry.mu.RLock()
	logger = globalRegistry.logger
	globalRegistry.mu.RUnlock()
	if err != nil {
		if logger != nil {
			logger.ErrorWithCategory(Category.Database.Database, SubCategory.Status.Error, "Failed to get connection in MustClient", map[string]interface{}{
				"connection_name": name,
				"error":           err.Error(),
			})
		}
		return nil
	}
	return &SfRedis{
		client: client,
		ctx:    ctx,
		logger: logger,
	}
}

// SafeClient returns a SfRedis client or error for the named connection
func SafeClient(ctx context.Context, name string) (*SfRedis, error) {
	client, err := getConnection(name)
	var logger Logger
	globalRegistry.mu.RLock()
	logger = globalRegistry.logger
	globalRegistry.mu.RUnlock()
	if err != nil {
		if logger != nil {
			logger.ErrorWithCategory(Category.Database.Database, SubCategory.Status.Error, "Failed to get connection in SafeClient", map[string]interface{}{
				"connection_name": name,
				"error":           err.Error(),
			})
		}
		return nil, err
	}
	return &SfRedis{
		client: client,
		ctx:    ctx,
		logger: logger,
	}, nil
}

// =============================================================================
// Context and Span Helpers
// =============================================================================

// withSpan creates a new span for an operation
func (r *SfRedis) withSpan(ctx context.Context, name string) *OperationContext {
	// Use provided context if available, otherwise use the default context
	if ctx == nil {
		ctx = r.ctx
	}

	// Create a new span with SF-Redis prefix
	span, ctx := apm.StartSpan(ctx, "SF-Redis."+name, "redis")
	return &OperationContext{
		ctx:    ctx,
		span:   span,
		logger: r.logger,
	}
}

// endSpan ends the span and handles any error
func (oc *OperationContext) endSpan(err error) {
	if err != nil {
		oc.span.Context.SetLabel("error", err.Error())

		// Log error if logger is available
		if oc.logger != nil {
			oc.logger.ErrorWithCategory(Category.Database.Database, Category.Error.Error, "Redis operation error", map[string]interface{}{
				"operation": oc.span.Name,
				"error":     err.Error(),
			})
		}
	}
	oc.span.End()
}

// =============================================================================
// SfRedis Methods
// =============================================================================

// Client returns the underlying SfRedis client
func (r *SfRedis) Client() redis.UniversalClient {
	return r.client
}

// ClientOption is a function to customize SfRedis client options
type ClientOption func(*redis.Options)

// WithTLS returns a new ClientOption to enable TLS
func WithTLS() ClientOption {
	return func(o *redis.Options) {
		o.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
}

// =============================================================================
// Basic Redis Operations
// =============================================================================

// Get gets the value of a key
func (r *SfRedis) Get(ctx context.Context, key string) (string, error) {
	oc := r.withSpan(ctx, "Get")
	defer oc.endSpan(nil)

	if r.logger != nil {
		r.logger.DebugWithCategory(Category.Database.Database, Category.System.Startup, "Getting Redis key", map[string]interface{}{
			"key": key,
		})
	}

	value, err := r.client.Get(oc.ctx, key).Result()
	oc.endSpan(err)
	return value, err
}

// Set sets the value of a key with an expiration time
func (r *SfRedis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	oc := r.withSpan(ctx, "Set")
	defer oc.endSpan(nil)

	if r.logger != nil {
		r.logger.DebugWithCategory(Category.Database.Database, Category.System.Startup, "Setting Redis key", map[string]interface{}{
			"key": key,
			"ttl": ttl.String(),
		})
	}

	err := r.client.Set(oc.ctx, key, value, ttl).Err()
	oc.endSpan(err)
	return err
}

// Delete deletes one or more keys
func (r *SfRedis) Delete(ctx context.Context, keys []string) (int64, error) {
	oc := r.withSpan(ctx, "Delete")
	defer oc.endSpan(nil)

	return r.client.Del(oc.ctx, keys...).Result()
}

// Exists checks if a key exists
func (r *SfRedis) Exists(ctx context.Context, key string) (bool, error) {
	oc := r.withSpan(ctx, "Exists")
	defer oc.endSpan(nil)

	result, err := r.client.Exists(oc.ctx, key).Result()
	return result > 0, err
}

// Expire sets a key's time to live in seconds
func (r *SfRedis) Expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	oc := r.withSpan(ctx, "Expire")
	defer oc.endSpan(nil)

	return r.client.Expire(oc.ctx, key, ttl).Result()
}

// TTL gets the time to live for a key in seconds
func (r *SfRedis) TTL(ctx context.Context, key string) (time.Duration, error) {
	oc := r.withSpan(ctx, "TTL")
	defer oc.endSpan(nil)

	return r.client.TTL(oc.ctx, key).Result()
}

// =============================================================================
// Lock Operations
// =============================================================================

// LockAcquire implements a distributed lock with SfRedis
func (r *SfRedis) LockAcquire(ctx context.Context, key string, value string, ttl time.Duration, retryDelay time.Duration, maxRetries int) (bool, error) {
	oc := r.withSpan(ctx, "LockAcquire")
	defer oc.endSpan(nil)

	lockKey := fmt.Sprintf("lock:%s", key)

	// Try to acquire the lock
	for i := 0; i <= maxRetries; i++ {
		// Use SET NX to ensure atomic lock acquisition
		locked, err := r.client.SetNX(oc.ctx, lockKey, value, ttl).Result()
		if err != nil {
			return false, err
		}

		if locked {
			return true, nil
		}

		// If we've hit max retries, fail
		if i == maxRetries {
			return false, errors.New("failed to acquire lock after max retries")
		}

		// Wait and retry
		time.Sleep(retryDelay)
	}

	return false, errors.New("failed to acquire lock")
}

// LockRelease releases a distributed lock
func (r *SfRedis) LockRelease(ctx context.Context, key string, value string) (bool, error) {
	oc := r.withSpan(ctx, "LockRelease")
	defer oc.endSpan(nil)

	lockKey := fmt.Sprintf("lock:%s", key)

	// Use Lua script to ensure atomic release of lock only if we own it
	// This prevents race conditions between checking lock owner and deleting it
	script := redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`)

	result, err := script.Run(oc.ctx, r.client, []string{lockKey}, value).Int64()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}

// =============================================================================
// JSON Operations
// =============================================================================

// SetJSON sets a JSON value for a key
func (r *SfRedis) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	oc := r.withSpan(ctx, "SetJSON")
	defer oc.endSpan(nil)

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(oc.ctx, key, data, expiration).Err()
}

// GetJSON gets a JSON value for a key and unmarshals it
func (r *SfRedis) GetJSON(ctx context.Context, key string, dest interface{}) error {
	oc := r.withSpan(ctx, "GetJSON")
	defer oc.endSpan(nil)

	data, err := r.client.Get(oc.ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// SetJSONIfNotExists sets a JSON value for a key only if it doesn't exist
func (r *SfRedis) SetJSONIfNotExists(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	oc := r.withSpan(ctx, "SetJSONIfNotExists")
	defer oc.endSpan(nil)

	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	return r.client.SetNX(oc.ctx, key, data, expiration).Result()
}

// HashSetJSON sets a JSON value for a hash field
func (r *SfRedis) HashSetJSON(ctx context.Context, key, field string, value interface{}) error {
	oc := r.withSpan(ctx, "HashSetJSON")
	defer oc.endSpan(nil)

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = r.client.HSet(oc.ctx, key, field, data).Result()
	return err
}

// HashGetJSON gets a JSON value from a hash field and unmarshals it
func (r *SfRedis) HashGetJSON(ctx context.Context, key, field string, dest interface{}) error {
	oc := r.withSpan(ctx, "HashGetJSON")
	defer oc.endSpan(nil)

	data, err := r.client.HGet(oc.ctx, key, field).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// =============================================================================
// Batch Operations
// =============================================================================

// BatchGet gets multiple keys efficiently using pipeline
func (r *SfRedis) BatchGet(ctx context.Context, keys []string) (map[string]string, error) {
	oc := r.withSpan(ctx, "BatchGet")
	defer oc.endSpan(nil)

	if len(keys) == 0 {
		return make(map[string]string), nil
	}

	pipe := r.client.Pipeline()
	cmds := make(map[string]*redis.StringCmd)

	// Queue all get commands
	for _, key := range keys {
		cmds[key] = pipe.Get(oc.ctx, key)
	}

	// Execute pipeline
	_, err := pipe.Exec(oc.ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	// Process results
	result := make(map[string]string, len(keys))
	for key, cmd := range cmds {
		val, err := cmd.Result()
		if err == nil {
			result[key] = val
		}
	}

	return result, nil
}

// BatchGetJSON gets multiple JSON values efficiently using pipeline
func (r *SfRedis) BatchGetJSON(ctx context.Context, keyDestPairs map[string]interface{}) error {
	oc := r.withSpan(ctx, "BatchGetJSON")
	defer oc.endSpan(nil)

	if len(keyDestPairs) == 0 {
		return nil
	}

	pipe := r.client.Pipeline()
	cmds := make(map[string]*redis.StringCmd)

	// Queue all get commands
	for key := range keyDestPairs {
		cmds[key] = pipe.Get(oc.ctx, key)
	}

	// Execute pipeline
	_, err := pipe.Exec(oc.ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	// Process results
	for key, cmd := range cmds {
		if dest, ok := keyDestPairs[key]; ok {
			val, err := cmd.Result()
			if err == nil {
				err = json.Unmarshal([]byte(val), dest)
				if err != nil {
					return fmt.Errorf("error unmarshaling key %s: %w", key, err)
				}
			}
		}
	}

	return nil
}

// BatchSet sets multiple key-value pairs efficiently using pipeline
func (r *SfRedis) BatchSet(ctx context.Context, keyValuePairs map[string]interface{}, ttl time.Duration) error {
	oc := r.withSpan(ctx, "BatchSet")
	defer oc.endSpan(nil)

	if len(keyValuePairs) == 0 {
		return nil
	}

	pipe := r.client.Pipeline()

	// Queue all set commands
	for key, value := range keyValuePairs {
		pipe.Set(oc.ctx, key, value, ttl)
	}

	// Execute pipeline
	_, err := pipe.Exec(oc.ctx)
	return err
}

// BatchDelete deletes multiple keys efficiently using pipeline
func (r *SfRedis) BatchDelete(ctx context.Context, keys []string) (int64, error) {
	oc := r.withSpan(ctx, "BatchDelete")
	defer oc.endSpan(nil)

	if len(keys) == 0 {
		return 0, nil
	}

	pipe := r.client.Pipeline()
	cmd := pipe.Del(oc.ctx, keys...)

	// Execute pipeline
	_, err := pipe.Exec(oc.ctx)
	if err != nil {
		return 0, err
	}

	return cmd.Result()
}

// BatchExists checks existence of multiple keys efficiently using pipeline
func (r *SfRedis) BatchExists(ctx context.Context, keys []string) (map[string]bool, error) {
	oc := r.withSpan(ctx, "BatchExists")
	defer oc.endSpan(nil)

	if len(keys) == 0 {
		return make(map[string]bool), nil
	}

	pipe := r.client.Pipeline()
	cmds := make(map[string]*redis.IntCmd)

	// Queue all exists commands
	for _, key := range keys {
		cmds[key] = pipe.Exists(oc.ctx, key)
	}

	// Execute pipeline
	_, err := pipe.Exec(oc.ctx)
	if err != nil {
		return nil, err
	}

	// Process results
	result := make(map[string]bool, len(keys))
	for key, cmd := range cmds {
		count, err := cmd.Result()
		if err == nil {
			result[key] = count > 0
		}
	}

	return result, nil
}

// getContext returns the first non-nil context from opts, or the default context
func getContext(opts ...context.Context) context.Context {
	if len(opts) > 0 && opts[0] != nil {
		return opts[0]
	}
	return nil // will be replaced by default context in withSpan
}

// =============================================================================
// Cache Operations
// =============================================================================

// GetOrSet gets a value from cache or executes the function to get it
func (r *SfRedis) GetOrSet(ctx context.Context, key string, ttl time.Duration, fn func() (interface{}, error)) (string, error) {
	oc := r.withSpan(ctx, "GetOrSet")
	defer oc.endSpan(nil)

	// Try to get from cache first
	val, err := r.client.Get(oc.ctx, key).Result()
	if err == nil {
		return val, nil
	}

	// If key doesn't exist, execute the function
	if errors.Is(err, redis.Nil) {
		result, err := fn()
		if err != nil {
			return "", err
		}

		// Store the result in cache
		err = r.client.Set(oc.ctx, key, result, ttl).Err()
		if err != nil {
			return "", err
		}

		// Convert result to string if possible
		switch v := result.(type) {
		case string:
			return v, nil
		case []byte:
			return string(v), nil
		default:
			return fmt.Sprintf("%v", result), nil
		}
	}

	return "", err
}

// GetJSONOrSet gets a JSON value from cache or executes the function to get it
func (r *SfRedis) GetJSONOrSet(ctx context.Context, key string, ttl time.Duration, dest interface{}, fn func() (interface{}, error)) error {
	oc := r.withSpan(ctx, "GetJSONOrSet")
	defer oc.endSpan(nil)

	if dest == nil {
		return errors.New("destination cannot be nil")
	}

	// Try to get from cache first
	err := r.GetJSON(oc.ctx, key, dest)
	if err == nil {
		return nil
	}

	// If key doesn't exist, execute the function
	if errors.Is(err, redis.Nil) {
		result, err := fn()
		if err != nil {
			return fmt.Errorf("failed to execute callback function: %w", err)
		}

		if result == nil {
			return errors.New("function returned nil result")
		}

		// Store the result in cache
		err = r.SetJSON(oc.ctx, key, result, ttl)
		if err != nil {
			return fmt.Errorf("failed to store result in cache: %w", err)
		}

		// Convert from the function's result to the destination
		resultBytes, err := json.Marshal(result)
		if err != nil {
			return fmt.Errorf("failed to marshal result to JSON: %w", err)
		}

		err = json.Unmarshal(resultBytes, dest)
		if err != nil {
			return fmt.Errorf("failed to unmarshal result into destination: %w", err)
		}

		return nil
	}

	return fmt.Errorf("failed to get value from cache: %w", err)
}

// =============================================================================
// Key Pattern Operations
// =============================================================================

// KeyGetByPattern gets all keys matching a pattern
func (r *SfRedis) KeyGetByPattern(ctx context.Context, pattern string) ([]string, error) {
	oc := r.withSpan(ctx, "KeyGetByPattern")
	defer oc.endSpan(nil)

	return r.client.Keys(oc.ctx, pattern).Result()
}

// KeyDeleteByPattern removes all keys matching a pattern
func (r *SfRedis) KeyDeleteByPattern(ctx context.Context, pattern string) (int64, error) {
	oc := r.withSpan(ctx, "KeyDeleteByPattern")
	defer oc.endSpan(nil)

	keys, err := r.KeyGetByPattern(oc.ctx, pattern)
	if err != nil {
		return 0, err
	}

	if len(keys) == 0 {
		return 0, nil
	}

	return r.client.Del(oc.ctx, keys...).Result()
}

// KeyScan iterates over keys matching a pattern
func (r *SfRedis) KeyScan(ctx context.Context, pattern string, count int64) ([]string, error) {
	oc := r.withSpan(ctx, "KeyScan")
	defer oc.endSpan(nil)

	var cursor uint64
	var keys []string

	for {
		var batch []string
		var err error

		batch, cursor, err = r.client.Scan(oc.ctx, cursor, pattern, count).Result()
		if err != nil {
			return nil, err
		}

		keys = append(keys, batch...)

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}

// =============================================================================
// Utility Operations
// =============================================================================

// SetWithRetry sets a key with a retry mechanism
func (r *SfRedis) SetWithRetry(ctx context.Context, key string, value interface{}, expiration time.Duration, maxRetries int) error {
	oc := r.withSpan(ctx, "SetWithRetry")
	defer oc.endSpan(nil)

	var err error
	for i := 0; i < maxRetries; i++ {
		err = r.client.Set(oc.ctx, key, value, expiration).Err()
		if err == nil {
			return nil
		}
		time.Sleep(time.Millisecond * time.Duration(50*(i+1)))
	}
	return fmt.Errorf("failed to set key after %d retries: %w", maxRetries, err)
}

// GetMemoryUsage returns the memory usage of a key
func (r *SfRedis) GetMemoryUsage(ctx context.Context, key string) (int64, error) {
	oc := r.withSpan(ctx, "GetMemoryUsage")
	defer oc.endSpan(nil)

	return r.client.MemoryUsage(oc.ctx, key).Result()
}

// DeleteManyWithBatch deletes multiple keys in batches to avoid blocking
func (r *SfRedis) DeleteManyWithBatch(ctx context.Context, keys []string, batchSize int) (int64, error) {
	oc := r.withSpan(ctx, "DeleteManyWithBatch")
	defer oc.endSpan(nil)

	if len(keys) == 0 {
		return 0, nil
	}

	var totalDeleted int64
	for i := 0; i < len(keys); i += batchSize {
		end := i + batchSize
		if end > len(keys) {
			end = len(keys)
		}

		batch := keys[i:end]
		deleted, err := r.client.Del(oc.ctx, batch...).Result()
		if err != nil {
			return totalDeleted, err
		}

		totalDeleted += deleted
	}

	return totalDeleted, nil
}

// StringGetSet gets a value from cache or executes the function to get it
func (r *SfRedis) StringGetSet(ctx context.Context, key string, ttl time.Duration, fn func() (string, error)) (string, error) {
	oc := r.withSpan(ctx, "StringGetSet")
	defer oc.endSpan(nil)

	// Try to get from cache first
	result, err := r.client.Get(oc.ctx, key).Result()
	if err == nil {
		return result, nil
	}

	// If key doesn't exist, execute the function
	if errors.Is(err, redis.Nil) {
		result, err := fn()
		if err != nil {
			return "", fmt.Errorf("failed to execute callback function: %w", err)
		}

		// Store the result in cache
		if err := r.client.Set(oc.ctx, key, result, ttl).Err(); err != nil {
			return result, fmt.Errorf("failed to store result in cache: %w", err)
		}

		return result, nil
	}

	return "", fmt.Errorf("failed to get value from cache: %w", err)
}
