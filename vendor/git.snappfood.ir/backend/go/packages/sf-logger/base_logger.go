package sflogger

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Config holds the logger configuration
type Config struct {
	// Core settings
	Level            Level
	FilePath         string
	AppName          string
	Development      bool
	EnableConsole    bool
	EnableStacktrace bool
	MaxSize          int // MB
	MaxAge           int // days
	MaxBackups       int
	Compress         bool

	// Logger implementation selection
	LoggerType LoggerType

	// Formatter configuration
	Formatter    FormatterType
	TimeFormat   string
	CSVSeparator string
	NoColors     bool

	// Log sinks configuration
	SinkURLs []string
}

// LoggerType specifies which logger implementation to use
type LoggerType string

const (
	ZapLoggerType  LoggerType = "zap"
	ZeroLoggerType LoggerType = "zerolog"
)

// DefaultConfig returns the default logging configuration
func DefaultConfig() Config {
	return Config{
		Level:            InfoLevel,
		FilePath:         "./logs/",
		AppName:          "App",
		Development:      false,
		EnableConsole:    false,
		EnableStacktrace: true,
		MaxSize:          1,
		MaxAge:           20,
		MaxBackups:       5,
		Compress:         true,
		LoggerType:       ZapLoggerType,
		Formatter:        JSONFormatter,
		TimeFormat:       "2006-01-02T15:04:05.000Z07:00",
		CSVSeparator:     ",",
		NoColors:         false,
		SinkURLs:         []string{},
	}
}

// New creates a new logger with options
func New(opts ...Option) Logger {
	config := DefaultConfig()

	// Apply all options
	for _, opt := range opts {
		opt(&config)
	}

	// Create base logger according to configuration
	var logger Logger
	switch config.LoggerType {
	case ZeroLoggerType:
		logger = newZeroLogger(config)
	default:
		logger = newZapLogger(config)
	}

	// Configure sinks
	if len(config.SinkURLs) > 0 {
		sinks := make([]Sink, 0, len(config.SinkURLs))

		// First, try to create sinks from URLs
		for _, sinkURL := range config.SinkURLs {
			sink, err := GetSink(sinkURL)
			if err == nil && sink != nil {
				sinks = append(sinks, sink)
			}
		}

		// If no sinks were created from URLs, create a default console sink
		if len(sinks) == 0 && len(config.SinkURLs) > 0 {
			sink, _ := GetSink("console://")
			if sink != nil {
				sinks = append(sinks, sink)
			}
		}

		// If we have sinks, wrap the logger with a SinkLogger
		if len(sinks) > 0 {
			logger = NewSinkLogger(logger, sinks...)
		}
	} else {
		// Add default sinks based on configuration
		var sinks []Sink

		// Add console sink if enabled
		if config.EnableConsole {
			sink, _ := GetSink("console://")
			if sink != nil {
				sinks = append(sinks, sink)
			}
		}

		// Add file sink if filepath is provided
		if config.FilePath != "" {
			url := fmt.Sprintf("file://%s?maxSize=%d&maxAge=%d&maxBackups=%d&compress=%t",
				config.FilePath,
				config.MaxSize,
				config.MaxAge,
				config.MaxBackups,
				config.Compress)
			sink, _ := GetSink(url)
			if sink != nil {
				sinks = append(sinks, sink)
			}
		}

		// If we have sinks, wrap the logger with a SinkLogger
		if len(sinks) > 0 {
			logger = NewSinkLogger(logger, sinks...)
		}
	}

	return logger
}

// SinkLogger wraps a Logger and sends logs to sinks
type SinkLogger struct {
	baseLogger Logger
	sinks      []Sink
	mutex      sync.RWMutex
}

// NewSinkLogger creates a new logger that sends logs to sinks
func NewSinkLogger(baseLogger Logger, sinks ...Sink) *SinkLogger {
	return &SinkLogger{
		baseLogger: baseLogger,
		sinks:      sinks,
	}
}

// AddSink adds a new sink to the logger
func (l *SinkLogger) AddSink(sink Sink) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.sinks = append(l.sinks, sink)
}

// RemoveSink removes a sink from the logger
func (l *SinkLogger) RemoveSink(sink Sink) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	for i, s := range l.sinks {
		if s == sink {
			l.sinks = append(l.sinks[:i], l.sinks[i+1:]...)
			break
		}
	}
}

// CloseSinks closes all sinks
func (l *SinkLogger) CloseSinks() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	for _, sink := range l.sinks {
		sink.Close()
	}
}

// Standard logging methods

func (l *SinkLogger) Debug(msg string, extra map[string]interface{}) {
	l.baseLogger.Debug(msg, extra)
	l.sendToSinks(msg, "DEBUG", extra)
}

func (l *SinkLogger) Debugf(template string, args ...interface{}) {
	l.baseLogger.Debugf(template, args...)
}

func (l *SinkLogger) DebugContext(ctx context.Context, msg string, extra map[string]interface{}) {
	l.baseLogger.DebugContext(ctx, msg, extra)
	if extra == nil {
		extra = make(map[string]interface{})
	}
	extra = addContextFields(ctx, extra)
	l.sendToSinks(msg, "DEBUG", extra)
}

func (l *SinkLogger) Info(msg string, extra map[string]interface{}) {
	l.baseLogger.Info(msg, extra)
	l.sendToSinks(msg, "INFO", extra)
}

func (l *SinkLogger) Infof(template string, args ...interface{}) {
	l.baseLogger.Infof(template, args...)
}

func (l *SinkLogger) InfoContext(ctx context.Context, msg string, extra map[string]interface{}) {
	l.baseLogger.InfoContext(ctx, msg, extra)
	if extra == nil {
		extra = make(map[string]interface{})
	}
	extra = addContextFields(ctx, extra)
	l.sendToSinks(msg, "INFO", extra)
}

func (l *SinkLogger) Warn(msg string, extra map[string]interface{}) {
	l.baseLogger.Warn(msg, extra)
	l.sendToSinks(msg, "WARN", extra)
}

func (l *SinkLogger) Warnf(template string, args ...interface{}) {
	l.baseLogger.Warnf(template, args...)
}

func (l *SinkLogger) WarnContext(ctx context.Context, msg string, extra map[string]interface{}) {
	l.baseLogger.WarnContext(ctx, msg, extra)
	if extra == nil {
		extra = make(map[string]interface{})
	}
	extra = addContextFields(ctx, extra)
	l.sendToSinks(msg, "WARN", extra)
}

func (l *SinkLogger) Error(msg string, extra map[string]interface{}) {
	l.baseLogger.Error(msg, extra)
	l.sendToSinks(msg, "ERROR", extra)
}

func (l *SinkLogger) Errorf(template string, args ...interface{}) {
	l.baseLogger.Errorf(template, args...)
}

func (l *SinkLogger) ErrorContext(ctx context.Context, msg string, extra map[string]interface{}) {
	l.baseLogger.ErrorContext(ctx, msg, extra)
	if extra == nil {
		extra = make(map[string]interface{})
	}
	extra = addContextFields(ctx, extra)
	l.sendToSinks(msg, "ERROR", extra)
}

func (l *SinkLogger) Fatal(msg string, extra map[string]interface{}) {
	l.baseLogger.Fatal(msg, extra)
	l.sendToSinks(msg, "FATAL", extra)
}

func (l *SinkLogger) Fatalf(template string, args ...interface{}) {
	l.baseLogger.Fatalf(template, args...)
}

func (l *SinkLogger) FatalContext(ctx context.Context, msg string, extra map[string]interface{}) {
	l.baseLogger.FatalContext(ctx, msg, extra)
	if extra == nil {
		extra = make(map[string]interface{})
	}
	extra = addContextFields(ctx, extra)
	l.sendToSinks(msg, "FATAL", extra)
}

// Category-based logging methods

func (l *SinkLogger) DebugWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	l.baseLogger.DebugWithCategory(cat, sub, msg, extra)
	l.sendCategoryToSinks(cat, sub, msg, "DEBUG", extra)
}

func (l *SinkLogger) InfoWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	l.baseLogger.InfoWithCategory(cat, sub, msg, extra)
	l.sendCategoryToSinks(cat, sub, msg, "INFO", extra)
}

func (l *SinkLogger) WarnWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	l.baseLogger.WarnWithCategory(cat, sub, msg, extra)
	l.sendCategoryToSinks(cat, sub, msg, "WARN", extra)
}

func (l *SinkLogger) ErrorWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	l.baseLogger.ErrorWithCategory(cat, sub, msg, extra)
	l.sendCategoryToSinks(cat, sub, msg, "ERROR", extra)
}

func (l *SinkLogger) FatalWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	l.baseLogger.FatalWithCategory(cat, sub, msg, extra)
	l.sendCategoryToSinks(cat, sub, msg, "FATAL", extra)
}

// Helper methods for sending logs to sinks

// sendToSinks sends a log entry to all configured sinks
func (l *SinkLogger) sendToSinks(msg, level string, extra map[string]interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	if len(l.sinks) == 0 {
		return
	}

	// Create log entry
	entry := map[string]interface{}{
		"level":     level,
		"msg":       msg,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// Add extra fields
	if extra != nil {
		for k, v := range extra {
			entry[k] = v
		}
	}

	// Send to all sinks
	for _, sink := range l.sinks {
		if err := sink.Write(entry); err != nil {
			// Log the error to the base logger
			errFields := map[string]interface{}{
				"sink_error": err.Error(),
			}
			l.baseLogger.Error("Failed to write to sink", errFields)
		}
	}
}

// sendCategoryToSinks sends a categorized log entry to all sinks
func (l *SinkLogger) sendCategoryToSinks(cat string, sub string, msg, level string, extra map[string]interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	if len(l.sinks) == 0 {
		return
	}

	// Create log entry
	entry := map[string]interface{}{
		"level":       level,
		"msg":         msg,
		"timestamp":   time.Now().Format(time.RFC3339),
		"Category":    cat,
		"SubCategory": sub,
	}

	// Add extra fields
	if extra != nil {
		for k, v := range extra {
			entry[k] = v
		}
	}

	// Send to all sinks
	for _, sink := range l.sinks {
		if err := sink.Write(entry); err != nil {
			// Log the error to the base logger
			errFields := map[string]interface{}{
				"sink_error":  err.Error(),
				"category":    cat,
				"subcategory": sub,
			}
			l.baseLogger.Error("Failed to write to sink", errFields)
		}
	}
}

// Global logger instance with mutex protection
var (
	globalLogger Logger
	loggerMutex  sync.RWMutex
	initialized  bool
)

// InitGlobalLogger initializes the global logger with options
func InitGlobalLogger(opts ...Option) Logger {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	globalLogger = New(opts...)
	initialized = true

	return globalLogger
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() Logger {
	loggerMutex.RLock()
	if initialized {
		defer loggerMutex.RUnlock()
		return globalLogger
	}
	loggerMutex.RUnlock()

	// Initialize with defaults if not initialized
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	if !initialized {
		globalLogger = New(
			WithConsole(true),
			WithFilePath(""),
		)
		initialized = true
	}

	return globalLogger
}

// Global helper functions
func Debug(msg string, extra map[string]interface{}) {
	GetGlobalLogger().Debug(msg, extra)
}

func Debugf(template string, args ...interface{}) {
	GetGlobalLogger().Debugf(template, args...)
}

func DebugContext(ctx context.Context, msg string, extra map[string]interface{}) {
	GetGlobalLogger().DebugContext(ctx, msg, extra)
}

func Info(msg string, extra map[string]interface{}) {
	GetGlobalLogger().Info(msg, extra)
}

func Infof(template string, args ...interface{}) {
	GetGlobalLogger().Infof(template, args...)
}

func InfoContext(ctx context.Context, msg string, extra map[string]interface{}) {
	GetGlobalLogger().InfoContext(ctx, msg, extra)
}

func Warn(msg string, extra map[string]interface{}) {
	GetGlobalLogger().Warn(msg, extra)
}

func Warnf(template string, args ...interface{}) {
	GetGlobalLogger().Warnf(template, args...)
}

func WarnContext(ctx context.Context, msg string, extra map[string]interface{}) {
	GetGlobalLogger().WarnContext(ctx, msg, extra)
}

func Error(msg string, extra map[string]interface{}) {
	GetGlobalLogger().Error(msg, extra)
}

func Errorf(template string, args ...interface{}) {
	GetGlobalLogger().Errorf(template, args...)
}

func ErrorContext(ctx context.Context, msg string, extra map[string]interface{}) {
	GetGlobalLogger().ErrorContext(ctx, msg, extra)
}

func Fatal(msg string, extra map[string]interface{}) {
	GetGlobalLogger().Fatal(msg, extra)
}

func Fatalf(template string, args ...interface{}) {
	GetGlobalLogger().Fatalf(template, args...)
}

func FatalContext(ctx context.Context, msg string, extra map[string]interface{}) {
	GetGlobalLogger().FatalContext(ctx, msg, extra)
}

func DebugWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	GetGlobalLogger().DebugWithCategory(cat, sub, msg, extra)
}

func InfoWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	GetGlobalLogger().InfoWithCategory(cat, sub, msg, extra)
}

func WarnWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	GetGlobalLogger().WarnWithCategory(cat, sub, msg, extra)
}

func ErrorWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	GetGlobalLogger().ErrorWithCategory(cat, sub, msg, extra)
}

func FatalWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	GetGlobalLogger().FatalWithCategory(cat, sub, msg, extra)
}
