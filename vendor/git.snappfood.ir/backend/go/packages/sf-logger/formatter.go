package sflogger

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// FormatterType represents the type of log formatter
type FormatterType string

const (
	// JSONFormatter formats logs as JSON
	JSONFormatter FormatterType = "json"

	// PrettyJSONFormatter formats logs as pretty-printed JSON
	PrettyJSONFormatter FormatterType = "pretty_json"

	// TextFormatter formats logs as plain text
	TextFormatter FormatterType = "text"

	// ColoredTextFormatter formats logs as colored text
	ColoredTextFormatter FormatterType = "colored"

	// CSVFormatter formats logs as CSV
	CSVFormatter FormatterType = "csv"

	// ZeroLogFormatter uses zerolog for formatting logs
	ZeroLogFormatter FormatterType = "zerolog"
)

// FormatterConfig holds configuration for log formatters
type FormatterConfig struct {
	Type         FormatterType
	TimeFormat   string
	CSVSeparator string
	NoColors     bool
}

// DefaultFormatterConfig returns default formatter configuration
func DefaultFormatterConfig() FormatterConfig {
	return FormatterConfig{
		Type:         JSONFormatter,
		TimeFormat:   time.RFC3339,
		CSVSeparator: ",",
		NoColors:     false,
	}
}

// GetEncoder returns a zapcore.Encoder configured according to the given FormatterConfig
func GetEncoder(config FormatterConfig) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Apply custom encoder settings if provided
	if config.TimeFormat != "" {
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(config.TimeFormat))
		}
	}

	if !config.NoColors {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	switch config.Type {
	case JSONFormatter:
		return zapcore.NewJSONEncoder(encoderConfig)

	case TextFormatter:
		return zapcore.NewConsoleEncoder(encoderConfig)

	case ColoredTextFormatter:
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(encoderConfig)

	case CSVFormatter:
		return newCSVEncoder(encoderConfig, config.CSVSeparator)

	case ZeroLogFormatter:
		return zapcore.NewJSONEncoder(encoderConfig)

	default:
		return zapcore.NewJSONEncoder(encoderConfig)
	}
}

// --- Custom encoder implementations ---

// PrettyJSON encoder - wraps the JSON encoder and pretty-prints the output
type prettyJSONEncoder struct {
	zapcore.Encoder
}

func newPrettyJSONEncoder(config zapcore.EncoderConfig) zapcore.Encoder {
	return &prettyJSONEncoder{Encoder: zapcore.NewJSONEncoder(config)}
}

func (e *prettyJSONEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return buf, err
	}

	// Pretty-print the JSON
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &jsonMap); err != nil {
		return buf, err
	}

	prettyJSON, err := json.MarshalIndent(jsonMap, "", "  ")
	if err != nil {
		return buf, err
	}

	buf.Reset()
	_, err = buf.Write(prettyJSON)
	return buf, err
}

// Text encoder - produces simple text-based logs
type textEncoder struct {
	zapcore.Encoder
	config  zapcore.EncoderConfig
	colored bool
}

func newTextEncoder(config zapcore.EncoderConfig, colored bool) zapcore.Encoder {
	return &textEncoder{
		Encoder: zapcore.NewJSONEncoder(config),
		config:  config,
		colored: colored,
	}
}

func (e *textEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// First encode to JSON to get the fields
	buf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return buf, err
	}

	// Decode the JSON to extract the fields
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &jsonMap); err != nil {
		return buf, err
	}

	// Format as text
	buf.Reset()

	// Time
	if timeStr, ok := jsonMap[e.config.TimeKey].(string); ok {
		buf.AppendString(timeStr)
		buf.AppendByte(' ')
	}

	// Level
	levelStr := entry.Level.CapitalString()
	if e.colored {
		levelColor := e.getLevelColor(entry.Level)
		buf.AppendString(fmt.Sprintf("\x1b[%sm%s\x1b[0m", levelColor, levelStr))
	} else {
		buf.AppendString(levelStr)
	}
	buf.AppendByte(' ')

	// Message
	if message, ok := jsonMap[e.config.MessageKey].(string); ok {
		buf.AppendString(message)
	}

	// Caller
	if caller, ok := jsonMap[e.config.CallerKey].(string); ok {
		buf.AppendString(" (")
		buf.AppendString(caller)
		buf.AppendString(")")
	}

	// Additional fields
	buf.AppendString(" ")
	for k, v := range jsonMap {
		if k != e.config.TimeKey && k != e.config.LevelKey &&
			k != e.config.MessageKey && k != e.config.CallerKey &&
			k != e.config.StacktraceKey {
			buf.AppendString(k)
			buf.AppendString("=")
			buf.AppendString(fmt.Sprintf("%v", v))
			buf.AppendString(" ")
		}
	}

	buf.AppendByte('\n')

	// Stack trace (if present)
	if stack, ok := jsonMap[e.config.StacktraceKey].(string); ok && stack != "" {
		buf.AppendString("Stacktrace:\n")
		buf.AppendString(stack)
		buf.AppendByte('\n')
	}

	return buf, nil
}

func (e *textEncoder) getLevelColor(level zapcore.Level) string {
	switch level {
	case zapcore.DebugLevel:
		return "36" // Cyan
	case zapcore.InfoLevel:
		return "32" // Green
	case zapcore.WarnLevel:
		return "33" // Yellow
	case zapcore.ErrorLevel:
		return "31" // Red
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		return "35" // Magenta
	default:
		return "37" // White
	}
}

// CSV encoder - produces comma-separated value logs
type csvEncoder struct {
	zapcore.Encoder
	config    zapcore.EncoderConfig
	separator string
}

func newCSVEncoder(config zapcore.EncoderConfig, separator string) zapcore.Encoder {
	return &csvEncoder{
		Encoder:   zapcore.NewJSONEncoder(config),
		config:    config,
		separator: separator,
	}
}

func (e *csvEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// First encode to JSON to get the fields
	buf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return buf, err
	}

	// Decode the JSON to extract the fields
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &jsonMap); err != nil {
		return buf, err
	}

	// Format as CSV
	buf.Reset()

	// Build CSV row
	values := []string{}

	// Always include these fields in this order
	// Time
	if timeStr, ok := jsonMap[e.config.TimeKey].(string); ok {
		values = append(values, e.escapeCSV(timeStr))
	} else {
		values = append(values, "")
	}

	// Level
	values = append(values, entry.Level.String())

	// Message
	if message, ok := jsonMap[e.config.MessageKey].(string); ok {
		values = append(values, e.escapeCSV(message))
	} else {
		values = append(values, "")
	}

	// Caller
	if caller, ok := jsonMap[e.config.CallerKey].(string); ok {
		values = append(values, e.escapeCSV(caller))
	} else {
		values = append(values, "")
	}

	// Join CSV fields
	buf.AppendString(strings.Join(values, e.separator))
	buf.AppendByte('\n')

	return buf, nil
}

// Escape CSV value by quoting if it contains the separator
func (e *csvEncoder) escapeCSV(value string) string {
	if strings.Contains(value, e.separator) || strings.Contains(value, "\"") ||
		strings.Contains(value, "\n") {
		return "\"" + strings.ReplaceAll(value, "\"", "\"\"") + "\""
	}
	return value
}
