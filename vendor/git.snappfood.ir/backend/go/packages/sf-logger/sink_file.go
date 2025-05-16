package sflogger

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// FileSinkImpl implements Sink for file output with rotation
type FileSinkImpl struct {
	levelLoggers   map[string]*lumberjack.Logger
	allLogsLogger  *lumberjack.Logger
	encoders       map[string]*json.Encoder
	allLogsEncoder *json.Encoder
	mutex          sync.Mutex
	baseFilePath   string
	config         *fileSinkConfig
	currentDate    string
	currentChunk   int
	rotationState  *rotationState
}

// rotationState tracks when files are rotated
type rotationState struct {
	lastCheck time.Time
	mutex     sync.Mutex
}

// FileSinkOption defines the function signature for file sink options
type FileSinkOption func(*fileSinkConfig)

// fileSinkConfig contains options for configuring a file sink
type fileSinkConfig struct {
	FilePath   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
}

// NewFileSink creates a new file sink with rotation
func NewFileSink(opts ...FileSinkOption) Sink {
	config := &fileSinkConfig{
		FilePath:   "./logs/app.log",
		MaxSize:    100,
		MaxAge:     30,
		MaxBackups: 5,
		Compress:   true,
	}

	// Apply options
	for _, opt := range opts {
		opt(config)
	}

	// Get base file path without extension
	baseFilePath := config.FilePath
	ext := filepath.Ext(baseFilePath)
	if ext != "" {
		baseFilePath = baseFilePath[:len(baseFilePath)-len(ext)]
	}

	// Get current date for filename
	currentDate := time.Now().Format("2006-01-02")

	// Create the rotation state
	rotState := &rotationState{
		lastCheck: time.Now(),
	}

	// Create the file sink instance first
	sink := &FileSinkImpl{
		levelLoggers:  make(map[string]*lumberjack.Logger),
		encoders:      make(map[string]*json.Encoder),
		baseFilePath:  baseFilePath,
		config:        config,
		currentDate:   currentDate,
		currentChunk:  1,
		rotationState: rotState,
	}

	// Create level-specific loggers
	levels := []string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}
	for _, level := range levels {
		levelFilePath := fmt.Sprintf("%s_%s_%s_chunk%d.log", baseFilePath, currentDate, level, sink.currentChunk)
		ljLogger := &lumberjack.Logger{
			Filename:   levelFilePath,
			MaxSize:    config.MaxSize,
			MaxAge:     config.MaxAge,
			MaxBackups: config.MaxBackups,
			Compress:   config.Compress,
		}
		sink.levelLoggers[level] = ljLogger
		sink.encoders[level] = json.NewEncoder(ljLogger)
		sink.encoders[level].SetEscapeHTML(false)
	}

	// Create all logs logger
	allLogsPath := fmt.Sprintf("%s_%s_ALL_chunk%d.log", baseFilePath, currentDate, sink.currentChunk)
	allLogsLogger := &lumberjack.Logger{
		Filename:   allLogsPath,
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
		MaxBackups: config.MaxBackups,
		Compress:   config.Compress,
	}
	sink.allLogsLogger = allLogsLogger
	sink.allLogsEncoder = json.NewEncoder(allLogsLogger)
	sink.allLogsEncoder.SetEscapeHTML(false)

	return sink
}

// Write sends a log entry to the appropriate file(s)
func (s *FileSinkImpl) Write(entry map[string]interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if we need to rotate files (once per second at most)
	if time.Since(s.rotationState.lastCheck) > time.Second {
		s.rotationState.lastCheck = time.Now()

		// Check if date has changed - if so, start with a new chunk 1
		currentDate := time.Now().Format("2006-01-02")
		if currentDate != s.currentDate {
			s.currentDate = currentDate
			s.currentChunk = 1
			if err := s.recreateAllLoggers(); err != nil {
				return err
			}
		} else if s.shouldRotate() {
			// Same date but size exceeded - increment chunk
			s.currentChunk++
			if err := s.recreateAllLoggers(); err != nil {
				return err
			}
		}
	}

	// Get the log level
	level, ok := entry["level"].(string)
	if !ok {
		level = "INFO" // Default to INFO if level not specified
	}

	// Write to level-specific file if we have a logger for this level
	var err error
	if encoder, exists := s.encoders[level]; exists {
		err = encoder.Encode(entry)
		if err != nil {
			return err
		}
	}

	// Write to the all logs file
	return s.allLogsEncoder.Encode(entry)
}

// shouldRotate checks if any log file exceeds the maximum size
func (s *FileSinkImpl) shouldRotate() bool {
	// Get file size for each logger
	maxSizeBytes := int64(s.config.MaxSize * 1024 * 1024) // Convert MB to bytes

	// Check all logs file first
	if s.allLogsLogger != nil {
		// Get current file size
		info, err := os.Stat(s.allLogsLogger.Filename)
		if err == nil && info.Size() >= maxSizeBytes {
			return true
		}
	}

	// Check each level-specific log file
	for _, logger := range s.levelLoggers {
		if logger != nil {
			// Get current file size
			info, err := os.Stat(logger.Filename)
			if err == nil && info.Size() >= maxSizeBytes {
				return true
			}
		}
	}

	return false
}

// recreateAllLoggers recreates all loggers with current date and chunk number
func (s *FileSinkImpl) recreateAllLoggers() error {
	// Close existing loggers first
	s.closeLoggers()

	// Create level-specific loggers with new chunk number
	levels := []string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}
	for _, level := range levels {
		levelFilePath := fmt.Sprintf("%s_%s_%s_chunk%d.log", s.baseFilePath, s.currentDate, level, s.currentChunk)
		ljLogger := &lumberjack.Logger{
			Filename:   levelFilePath,
			MaxSize:    s.config.MaxSize,
			MaxAge:     s.config.MaxAge,
			MaxBackups: s.config.MaxBackups,
			Compress:   s.config.Compress,
		}
		s.levelLoggers[level] = ljLogger
		s.encoders[level] = json.NewEncoder(ljLogger)
		s.encoders[level].SetEscapeHTML(false)
	}

	// Create all-logs logger with new chunk number
	allLogsPath := fmt.Sprintf("%s_%s_ALL_chunk%d.log", s.baseFilePath, s.currentDate, s.currentChunk)
	allLogsLogger := &lumberjack.Logger{
		Filename:   allLogsPath,
		MaxSize:    s.config.MaxSize,
		MaxAge:     s.config.MaxAge,
		MaxBackups: s.config.MaxBackups,
		Compress:   s.config.Compress,
	}
	s.allLogsLogger = allLogsLogger
	s.allLogsEncoder = json.NewEncoder(allLogsLogger)
	s.allLogsEncoder.SetEscapeHTML(false)

	return nil
}

// closeLoggers closes all existing loggers
func (s *FileSinkImpl) closeLoggers() {
	// Close level loggers
	for _, logger := range s.levelLoggers {
		if logger != nil {
			logger.Close()
		}
	}

	// Close all logs logger
	if s.allLogsLogger != nil {
		s.allLogsLogger.Close()
	}
}

// Close closes all file sinks
func (s *FileSinkImpl) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.closeLoggers()
	return nil
}

// Sync ensures all log entries are written to disk
func (s *FileSinkImpl) Sync() error {
	// Lumberjack doesn't expose a Sync method
	return nil
}

// FileWithPath sets the file path for logs
func FileWithPath(path string) FileSinkOption {
	return func(c *fileSinkConfig) {
		c.FilePath = path
	}
}

// FileWithMaxSize sets the maximum size in megabytes before log rotation
func FileWithMaxSize(maxSize int) FileSinkOption {
	return func(c *fileSinkConfig) {
		c.MaxSize = maxSize
	}
}

// FileWithMaxAge sets the maximum number of days to retain old logs
func FileWithMaxAge(days int) FileSinkOption {
	return func(c *fileSinkConfig) {
		c.MaxAge = days
	}
}

// FileWithMaxBackups sets the maximum number of old log files to keep
func FileWithMaxBackups(count int) FileSinkOption {
	return func(c *fileSinkConfig) {
		c.MaxBackups = count
	}
}

// FileWithCompress enables or disables log compression
func FileWithCompress(compress bool) FileSinkOption {
	return func(c *fileSinkConfig) {
		c.Compress = compress
	}
}

// Register file sink with the registry
func init() {
	RegisterSink("file", func(u *url.URL) (Sink, error) {
		var opts []FileSinkOption

		// Get path from URL path
		path := u.Path
		if path != "" {
			// Remove leading slash if present
			if path[0] == '/' {
				path = path[1:]
			}
			// Make path absolute if not already
			if !filepath.IsAbs(path) {
				path = filepath.Join(".", path)
			}
			opts = append(opts, FileWithPath(path))
		}

		// Parse query parameters
		q := u.Query()

		if maxSize := q.Get("maxSize"); maxSize != "" {
			if size, err := strconv.Atoi(maxSize); err == nil {
				opts = append(opts, FileWithMaxSize(size))
			}
		}

		if maxAge := q.Get("maxAge"); maxAge != "" {
			if age, err := strconv.Atoi(maxAge); err == nil {
				opts = append(opts, FileWithMaxAge(age))
			}
		}

		if maxBackups := q.Get("maxBackups"); maxBackups != "" {
			if count, err := strconv.Atoi(maxBackups); err == nil {
				opts = append(opts, FileWithMaxBackups(count))
			}
		}

		if compress := q.Get("compress"); compress != "" {
			opts = append(opts, FileWithCompress(compress == "true"))
		}

		return NewFileSink(opts...), nil
	})
}
