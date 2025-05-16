package sflogger

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

var zeroSinLogger *zerolog.Logger

// zeroLogger implements Logger using zerolog
type zeroLogger struct {
	logger *zerolog.Logger
	fields map[string]interface{}
	config Config
}

// Mapping from our log levels to zerolog levels
var zeroLogLevelMapping = map[string]zerolog.Level{
	levelNames[DebugLevel]:   zerolog.DebugLevel,
	levelNames[InfoLevel]:    zerolog.InfoLevel,
	levelNames[WarningLevel]: zerolog.WarnLevel,
	levelNames[ErrorLevel]:   zerolog.ErrorLevel,
	levelNames[FatalLevel]:   zerolog.FatalLevel,
}

// zerologSinkWriter implements io.Writer by writing to a Sink
type zerologSinkWriter struct {
	sink Sink
}

func (w *zerologSinkWriter) Write(p []byte) (n int, err error) {
	// Create a basic log entry
	entry := map[string]interface{}{
		"msg": string(p),
	}

	if err := w.sink.Write(entry); err != nil {
		return 0, err
	}

	return len(p), nil
}

// NewZerologSinkWriter creates a new io.Writer that writes to a Sink
func NewZerologSinkWriter(sink Sink) io.Writer {
	return &zerologSinkWriter{sink: sink}
}

// newZeroLogger creates a new zerolog logger from configuration
func newZeroLogger(config Config) *zeroLogger {
	logger := &zeroLogger{
		fields: make(map[string]interface{}),
		config: config,
	}

	// Configure error stack marshaler
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// Set global log level
	zerolog.SetGlobalLevel(logger.getZeroLevel())

	// Setup writers
	var writers []io.Writer

	// Setup file output if path is provided
	if config.FilePath != "" {
		fileName := fmt.Sprintf("%s%s-%s.%s", config.FilePath, time.Now().Format("2006-01-02"), uuid.New(), "log")

		// Configure file rotation
		rotatingLogger := &lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    config.MaxSize,
			MaxAge:     config.MaxAge,
			MaxBackups: config.MaxBackups,
			Compress:   config.Compress,
		}
		writers = append(writers, rotatingLogger)
	}

	// Setup console output if enabled
	if config.EnableConsole {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: config.TimeFormat,
			NoColor:    config.NoColors,
		}
		writers = append(writers, consoleWriter)
	}

	// If no writers configured, default to stdout
	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	// Create multi-writer if we have multiple outputs
	var output io.Writer
	if len(writers) > 1 {
		output = zerolog.MultiLevelWriter(writers...)
	} else {
		output = writers[0]
	}

	// Create and configure the logger
	zlogger := zerolog.New(output).
		With().
		Timestamp().
		Str("AppName", config.AppName).
		Str("LoggerName", "Zerolog").
		Logger()

	// Enable stack traces if configured
	if config.EnableStacktrace {
		zlogger = zlogger.With().Stack().Logger()
	}

	logger.logger = &zlogger
	return logger
}

// getZeroLevel converts our level to zerolog level
func (l *zeroLogger) getZeroLevel() zerolog.Level {
	level, exists := zeroLogLevelMapping[l.config.Level.String()]
	if !exists {
		return zerolog.DebugLevel
	}
	return level
}

// Standard logging methods

func (l *zeroLogger) Debug(msg string, extra map[string]interface{}) {
	event := l.logger.Debug()
	if extra != nil {
		event = event.Fields(extra)
	}
	event.Msg(msg)
}

func (l *zeroLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debug().Msgf(template, args...)
}

func (l *zeroLogger) DebugContext(ctx context.Context, msg string, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}
	extra = addContextFields(ctx, extra)
	l.Debug(msg, extra)
}

func (l *zeroLogger) Info(msg string, extra map[string]interface{}) {
	event := l.logger.Info()
	if extra != nil {
		event = event.Fields(extra)
	}
	event.Msg(msg)
}

func (l *zeroLogger) Infof(template string, args ...interface{}) {
	l.logger.Info().Msgf(template, args...)
}

func (l *zeroLogger) InfoContext(ctx context.Context, msg string, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}
	extra = addContextFields(ctx, extra)
	l.Info(msg, extra)
}

func (l *zeroLogger) Warn(msg string, extra map[string]interface{}) {
	event := l.logger.Warn()
	if extra != nil {
		event = event.Fields(extra)
	}
	event.Msg(msg)
}

func (l *zeroLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warn().Msgf(template, args...)
}

func (l *zeroLogger) WarnContext(ctx context.Context, msg string, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}
	extra = addContextFields(ctx, extra)
	l.Warn(msg, extra)
}

func (l *zeroLogger) Error(msg string, extra map[string]interface{}) {
	event := l.logger.Error()
	if extra != nil {
		event = event.Fields(extra)
	}
	event.Msg(msg)
}

func (l *zeroLogger) Errorf(template string, args ...interface{}) {
	l.logger.Error().Msgf(template, args...)
}

func (l *zeroLogger) ErrorContext(ctx context.Context, msg string, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}
	extra = addContextFields(ctx, extra)
	l.Error(msg, extra)
}

func (l *zeroLogger) Fatal(msg string, extra map[string]interface{}) {
	event := l.logger.Fatal()
	if extra != nil {
		event = event.Fields(extra)
	}
	event.Msg(msg)
}

func (l *zeroLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatal().Msgf(template, args...)
}

func (l *zeroLogger) FatalContext(ctx context.Context, msg string, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}
	extra = addContextFields(ctx, extra)
	l.Fatal(msg, extra)
}

// Category-based logging methods

func (l *zeroLogger) DebugWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	l.logger.Debug().
		Str("Category", cat).
		Str("SubCategory", sub).
		Fields(extra).
		Msg(msg)
}

func (l *zeroLogger) InfoWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	l.logger.Info().
		Str("Category", cat).
		Str("SubCategory", sub).
		Fields(extra).
		Msg(msg)
}

func (l *zeroLogger) WarnWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	l.logger.Warn().
		Str("Category", cat).
		Str("SubCategory", sub).
		Fields(extra).
		Msg(msg)
}

func (l *zeroLogger) ErrorWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	l.logger.Error().
		Str("Category", cat).
		Str("SubCategory", sub).
		Fields(extra).
		Msg(msg)
}

func (l *zeroLogger) FatalWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	l.logger.Fatal().
		Str("Category", cat).
		Str("SubCategory", sub).
		Fields(extra).
		Msg(msg)
}
