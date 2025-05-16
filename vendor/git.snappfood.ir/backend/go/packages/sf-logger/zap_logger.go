package sflogger

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var once sync.Once
var zapSinLogger *zap.SugaredLogger

// zapLogger implements Logger using zap
type zapLogger struct {
	logger *zap.SugaredLogger
	fields map[string]interface{}
	config Config
}

// Mapping from our log levels to zap log levels
var zapLogLevelMapping = map[string]zapcore.Level{
	levelNames[DebugLevel]:   zapcore.DebugLevel,
	levelNames[InfoLevel]:    zapcore.InfoLevel,
	levelNames[WarningLevel]: zapcore.WarnLevel,
	levelNames[ErrorLevel]:   zapcore.ErrorLevel,
	levelNames[FatalLevel]:   zapcore.FatalLevel,
}

// newZapLogger creates a new zap logger from configuration
func newZapLogger(config Config) *zapLogger {
	logger := &zapLogger{
		fields: make(map[string]interface{}),
		config: config,
	}

	// Setup encoder based on formatter configuration
	formatterConfig := FormatterConfig{
		Type:         config.Formatter,
		TimeFormat:   config.TimeFormat,
		CSVSeparator: config.CSVSeparator,
		NoColors:     config.NoColors,
	}

	encoder := GetEncoder(formatterConfig)

	// Create a multi-writer with appropriate sinks
	var cores []zapcore.Core

	// Add file sink if configured
	if config.FilePath != "" {
		fileName := fmt.Sprintf("%s%s-%s.%s", config.FilePath, time.Now().Format("2006-01-02"), uuid.New(), "log")
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    config.MaxSize,
			MaxAge:     config.MaxAge,
			LocalTime:  true,
			MaxBackups: config.MaxBackups,
			Compress:   config.Compress,
		})

		fileCore := zapcore.NewCore(encoder, fileWriter, logger.getZapLevel())
		cores = append(cores, fileCore)
	}

	// Add console sink if configured
	if config.EnableConsole {
		var consoleEncoder zapcore.Encoder
		if config.Development {
			// In development mode, use colored text by default
			consoleFormatterConfig := formatterConfig
			consoleFormatterConfig.Type = ColoredTextFormatter
			consoleEncoder = GetEncoder(consoleFormatterConfig)
		} else {
			consoleEncoder = encoder
		}

		consoleWriter := zapcore.AddSync(os.Stdout)
		consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, logger.getZapLevel())
		cores = append(cores, consoleCore)
	}

	// If no cores were configured, default to console
	if len(cores) == 0 {
		consoleWriter := zapcore.AddSync(os.Stdout)
		consoleCore := zapcore.NewCore(encoder, consoleWriter, logger.getZapLevel())
		cores = append(cores, consoleCore)
	}

	// Combine cores into a single core
	var core zapcore.Core
	if len(cores) > 1 {
		core = zapcore.NewTee(cores...)
	} else {
		core = cores[0]
	}

	// Create logger with options
	options := []zap.Option{zap.AddCaller(), zap.AddCallerSkip(1)}
	if config.EnableStacktrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	zapLogger := zap.New(core, options...).Sugar()

	// Add default fields
	zapLogger = zapLogger.With("AppName", config.AppName, "LoggerName", "Zaplog")

	logger.logger = zapLogger
	return logger
}

// getZapLevel converts our level to zap level
func (l *zapLogger) getZapLevel() zapcore.Level {
	level, exists := zapLogLevelMapping[l.config.Level.String()]
	if !exists {
		return zapcore.DebugLevel
	}
	return level
}

// Standard logging methods

func (l *zapLogger) Debug(msg string, extra map[string]interface{}) {
	params := prepareLogParams(extra)
	l.logger.Debugw(msg, params...)
}

func (l *zapLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *zapLogger) DebugContext(ctx context.Context, msg string, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}
	extra = addContextFields(ctx, extra)
	l.Debug(msg, extra)
}

func (l *zapLogger) Info(msg string, extra map[string]interface{}) {
	params := prepareLogParams(extra)
	l.logger.Infow(msg, params...)
}

func (l *zapLogger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *zapLogger) InfoContext(ctx context.Context, msg string, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}
	extra = addContextFields(ctx, extra)
	l.Info(msg, extra)
}

func (l *zapLogger) Warn(msg string, extra map[string]interface{}) {
	params := prepareLogParams(extra)
	l.logger.Warnw(msg, params...)
}

func (l *zapLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *zapLogger) WarnContext(ctx context.Context, msg string, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}
	extra = addContextFields(ctx, extra)
	l.Warn(msg, extra)
}

func (l *zapLogger) Error(msg string, extra map[string]interface{}) {
	params := prepareLogParams(extra)
	l.logger.Errorw(msg, params...)
}

func (l *zapLogger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l *zapLogger) ErrorContext(ctx context.Context, msg string, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}
	extra = addContextFields(ctx, extra)
	l.Error(msg, extra)
}

func (l *zapLogger) Fatal(msg string, extra map[string]interface{}) {
	params := prepareLogParams(extra)
	l.logger.Fatalw(msg, params...)
}

func (l *zapLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

func (l *zapLogger) FatalContext(ctx context.Context, msg string, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}
	extra = addContextFields(ctx, extra)
	l.Fatal(msg, extra)
}

// Category-based logging methods

func (l *zapLogger) DebugWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	params := prepareCategoryParams(cat, sub, extra)
	l.logger.Debugw(msg, params...)
}

func (l *zapLogger) InfoWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	params := prepareCategoryParams(cat, sub, extra)
	l.logger.Infow(msg, params...)
}

func (l *zapLogger) WarnWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	params := prepareCategoryParams(cat, sub, extra)
	l.logger.Warnw(msg, params...)
}

func (l *zapLogger) ErrorWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	params := prepareCategoryParams(cat, sub, extra)
	l.logger.Errorw(msg, params...)
}

func (l *zapLogger) FatalWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	params := prepareCategoryParams(cat, sub, extra)
	l.logger.Fatalw(msg, params...)
}

// Helper functions for preparing log parameters

// prepareLogParams converts a map to a list of interface{} for zap's structured logging
func prepareLogParams(extra map[string]interface{}) []interface{} {
	if extra == nil {
		return []interface{}{}
	}

	params := make([]interface{}, 0, len(extra)*2)
	for k, v := range extra {
		params = append(params, k, v)
	}
	return params
}

// prepareCategoryParams prepares parameters for category-based logging
func prepareCategoryParams(cat string, sub string, extra map[string]interface{}) []interface{} {
	params := []interface{}{
		"Category", cat,
		"SubCategory", sub,
	}

	// Add extra fields
	if extra != nil {
		for k, v := range extra {
			params = append(params, k, v)
		}
	}

	return params
}

// Context helpers

type contextKey string

const (
	requestIDKey contextKey = "request_id"
	traceIDKey   contextKey = "trace_id"
	userIDKey    contextKey = "user_id"
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// Helper functions for context values
func addContextFields(ctx context.Context, fields map[string]interface{}) map[string]interface{} {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		fields["request_id"] = requestID
	}

	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		fields["trace_id"] = traceID
	}

	if userID, ok := ctx.Value(userIDKey).(string); ok {
		fields["user_id"] = userID
	}

	return fields
}

// registerZapSinkFactories registers our sinks with zap
func registerZapSinkFactories() {
	// Register sink factories with zap
	_ = zap.RegisterSink("sflogger", func(u *url.URL) (zap.Sink, error) {
		// Create our sink
		sink, err := GetSink(u.String())
		if err != nil {
			return nil, err
		}

		// Adapt our Sink to zap.Sink
		return &zapSinkAdapter{sink: sink}, nil
	})
}

// zapSinkAdapter adapts our Sink interface to zap.Sink
type zapSinkAdapter struct {
	sink Sink
}

func (a *zapSinkAdapter) Write(p []byte) (n int, err error) {
	// Create a basic log entry
	entry := map[string]interface{}{
		"msg": string(p),
	}

	if err := a.sink.Write(entry); err != nil {
		return 0, err
	}

	return len(p), nil
}

func (a *zapSinkAdapter) Sync() error {
	return a.sink.Sync()
}

func (a *zapSinkAdapter) Close() error {
	return a.sink.Close()
}

func init() {
	registerZapSinkFactories()
}
