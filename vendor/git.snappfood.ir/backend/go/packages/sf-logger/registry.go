package sflogger

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
	"strings"
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/natefinch/lumberjack"

	"bytes"
)

// LoggerType enumeration
type LoggerType int

const (
	ZapLoggerType LoggerType = iota
)

type Level zapcore.Level

const (
	DebugLevel Level = Level(zap.DebugLevel)
	InfoLevel  Level = Level(zap.InfoLevel)
	WarnLevel  Level = Level(zap.WarnLevel)
	ErrorLevel Level = Level(zap.ErrorLevel)
	FatalLevel Level = Level(zap.FatalLevel)
)

// FormatterType enumeration
type FormatterType int

const (
	ColoredTextFormatter FormatterType = iota
	JSONFormatter
)

// Option defines a functional option for logger configuration
type Option func(*Registry)

// Registry holds logger config and zap instance
type Registry struct {
	mu           sync.RWMutex
	loggerEngine *zap.SugaredLogger
	loggerType   LoggerType
	level        Level
	appName      string
	formatter    FormatterType
	stacktrace   bool
	fileSinkCfg  *fileSinkConfig
	mongoSinkCfg *mongoSinkConfig
	esSinkCfg    *elasticSinkConfig
}

type fileSinkConfig struct {
	path       string
	maxSizeMB  int
	maxAgeDays int
	maxBackups int
	compress   bool
}

type mongoSinkConfig struct {
	host       string
	port       int
	database   string
	collection string
	username   string
	password   string
	flushSec   int
	compress   bool
}

type elasticSinkConfig struct {
	url      string
	username string
	password string
	index    string
	flushSec int
}

// Global singleton registry
var globalRegistry = &Registry{}

// RegisterLogger initializes logger based on options, returns Logger interface
func RegisterLogger(opts ...Option) Logger {
	for _, opt := range opts {
		opt(globalRegistry)
	}
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.Level(globalRegistry.level)),
		Development: false,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		InitialFields:     map[string]interface{}{"appName": globalRegistry.appName},
		DisableStacktrace: !globalRegistry.stacktrace,
	}

	if globalRegistry.formatter == JSONFormatter {
		cfg.Encoding = "json"
		// adjust encoder config for json if needed
		cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	}

	var cores []zapcore.Core

	// Base core writes to stdout
	consoleSyncer := zapcore.Lock(os.Stdout)
	encoder := getEncoder(cfg)
	cores = append(cores, zapcore.NewCore(encoder, consoleSyncer, zap.NewAtomicLevelAt(zapcore.Level(globalRegistry.level))))

	// Add file sink if configured
	if globalRegistry.fileSinkCfg != nil {
		fs, err := newFileSink(globalRegistry.fileSinkCfg)
		if err == nil {
			fileCore := zapcore.NewCore(encoder, fs, zap.NewAtomicLevelAt(zapcore.Level(globalRegistry.level)))
			cores = append(cores, fileCore)
		} else {
			fmt.Fprintf(os.Stderr, "file sink error: %v\n", err)
		}
	}

	// Add MongoDB sink if configured
	if globalRegistry.mongoSinkCfg != nil {
		ms, err := newMongoSink(globalRegistry.mongoSinkCfg)
		if err == nil {
			mongoCore := zapcore.NewCore(encoder, ms, zap.NewAtomicLevelAt(zapcore.Level(globalRegistry.level)))
			cores = append(cores, mongoCore)
		} else {
			fmt.Fprintf(os.Stderr, "mongo sink error: %v\n", err)
		}
	}

	// Add Elasticsearch sink if configured
	if globalRegistry.esSinkCfg != nil {
		es, err := newElasticSink(globalRegistry.esSinkCfg)
		if err == nil {
			esCore := zapcore.NewCore(encoder, es, zap.NewAtomicLevelAt(zapcore.Level(globalRegistry.level)))
			cores = append(cores, esCore)
		} else {
			fmt.Fprintf(os.Stderr, "elasticsearch sink error: %v\n", err)
		}
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	globalRegistry.loggerEngine = logger.Sugar()

	return newZapLogger(globalRegistry.loggerEngine)
}

// Functional options

func WithLoggerType(t LoggerType) Option {
	return func(r *Registry) {
		r.loggerType = t
	}
}

func WithLevel(l Level) Option {
	return func(r *Registry) {
		r.level = l
	}
}

func WithAppName(name string) Option {
	return func(r *Registry) {
		r.appName = name
	}
}

func WithFormatter(f FormatterType) Option {
	return func(r *Registry) {
		r.formatter = f
	}
}

func WithStacktrace(enabled bool) Option {
	return func(r *Registry) {
		r.stacktrace = enabled
	}
}

func WithFileSink(path string, maxSizeMB, maxAgeDays, maxBackups int, compress bool) Option {
	return func(r *Registry) {
		r.fileSinkCfg = &fileSinkConfig{
			path:       path,
			maxSizeMB:  maxSizeMB,
			maxAgeDays: maxAgeDays,
			maxBackups: maxBackups,
			compress:   compress,
		}
	}
}

func WithMongoDBSink(host string, port int, database, collection, username, password string, flushSec int, compress bool) Option {
	return func(r *Registry) {
		r.mongoSinkCfg = &mongoSinkConfig{
			host:       host,
			port:       port,
			database:   database,
			collection: collection,
			username:   username,
			password:   password,
			flushSec:   flushSec,
			compress:   compress,
		}
	}
}

func WithElasticSearchSink(url, username, password, index string, flushSec int) Option {
	return func(r *Registry) {
		r.esSinkCfg = &elasticSinkConfig{
			url:      url,
			username: username,
			password: password,
			index:    index,
			flushSec: flushSec,
		}
	}
}

// --- Encoders

func getEncoder(cfg zap.Config) zapcore.Encoder {
	if cfg.Encoding == "json" {
		return zapcore.NewJSONEncoder(cfg.EncoderConfig)
	}
	return zapcore.NewConsoleEncoder(cfg.EncoderConfig)
}

func newFileSink(cfg *fileSinkConfig) (zapcore.WriteSyncer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("file sink config is nil")
	}
	writer := &lumberjack.Logger{
		Filename:   cfg.path,
		MaxSize:    cfg.maxSizeMB,
		MaxAge:     cfg.maxAgeDays,
		MaxBackups: cfg.maxBackups,
		Compress:   cfg.compress,
	}
	return zapcore.AddSync(writer), nil
}

// --- MongoDB sink implementation

type mongoSink struct {
	client     *mongo.Client
	collection *mongo.Collection
	buffer     []map[string]interface{}
	mu         sync.Mutex
	flushSec   int
	closed     chan struct{}
}

func newMongoSink(cfg *mongoSinkConfig) (zapcore.WriteSyncer, error) {
	ctx := context.Background()
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", cfg.username, cfg.password, cfg.host, cfg.port, cfg.database)
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	coll := client.Database(cfg.database).Collection(cfg.collection)

	ms := &mongoSink{
		client:     client,
		collection: coll,
		buffer:     make([]map[string]interface{}, 0, 100),
		flushSec:   cfg.flushSec,
		closed:     make(chan struct{}),
	}

	go ms.flusher(ctx)

	return ms, nil
}

func (m *mongoSink) Write(p []byte) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	logStr := string(p)
	doc := map[string]interface{}{
		"timestamp": time.Now(),
		"raw_log":   logStr,
	}

	jsonStart := strings.LastIndex(logStr, "{")
	jsonEnd := strings.LastIndex(logStr, "}")
	if jsonStart != -1 && jsonEnd != -1 && jsonEnd > jsonStart {
		jsonPart := logStr[jsonStart : jsonEnd+1]
		var extraFields map[string]interface{}
		if err := json.Unmarshal([]byte(jsonPart), &extraFields); err == nil {
			for k, v := range extraFields {
				doc[k] = v
			}
		}
	}

	m.buffer = append(m.buffer, doc)
	return len(p), nil
}

func (m *mongoSink) Sync() error {
	return nil
}

func (m *mongoSink) flusher(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(m.flushSec) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m.mu.Lock()
			if len(m.buffer) > 0 {
				docs := m.buffer
				m.buffer = make([]map[string]interface{}, 0, 100)
				m.mu.Unlock()

				// Convert []map[string]interface{} to []interface{}
				ifaceDocs := make([]interface{}, len(docs))
				for i, doc := range docs {
					ifaceDocs[i] = doc
				}

				_, err := m.collection.InsertMany(ctx, ifaceDocs)
				if err != nil {
					fmt.Fprintf(os.Stderr, "mongo sink insert error: %v\n", err)
				}
			} else {
				m.mu.Unlock()
			}
		case <-m.closed:
			return
		}
	}
}

func (m *mongoSink) Close() error {
	close(m.closed)
	return m.client.Disconnect(context.Background())
}

// --- Elasticsearch sink implementation

type elasticSink struct {
	client   *elasticsearch.Client
	buffer   [][]byte
	mu       sync.Mutex
	flushSec int
	closed   chan struct{}
}

// Move this type to top-level
// Wraps elasticSink and adds index name

type elasticSinkWithIndex struct {
	*elasticSink
	index string
}

func newElasticSink(cfg *elasticSinkConfig) (zapcore.WriteSyncer, error) {
	cfgES := elasticsearch.Config{
		Addresses: []string{cfg.url},
	}
	if cfg.username != "" && cfg.password != "" {
		cfgES.Username = cfg.username
		cfgES.Password = cfg.password
	}
	client, err := elasticsearch.NewClient(cfgES)
	if err != nil {
		return nil, err
	}

	es := &elasticSink{
		client:   client,
		buffer:   make([][]byte, 0, 100),
		flushSec: cfg.flushSec,
		closed:   make(chan struct{}),
	}

	esWithIndex := &elasticSinkWithIndex{
		elasticSink: es,
		index:       cfg.index,
	}

	go esWithIndex.flusher()

	return esWithIndex, nil
}

func (e *elasticSinkWithIndex) Write(p []byte) (n int, err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	logStr := string(p)
	doc := map[string]interface{}{
		"timestamp": time.Now(),
		"raw_log":   logStr,
	}

	jsonStart := strings.LastIndex(logStr, "{")
	jsonEnd := strings.LastIndex(logStr, "}")
	if jsonStart != -1 && jsonEnd != -1 && jsonEnd > jsonStart {
		jsonPart := logStr[jsonStart : jsonEnd+1]
		var extraFields map[string]interface{}
		if err := json.Unmarshal([]byte(jsonPart), &extraFields); err == nil {
			for k, v := range extraFields {
				doc[k] = v
			}
		}
	}

	jsonDoc, err := json.Marshal(doc)
	if err != nil {
		return 0, err
	}
	e.buffer = append(e.buffer, jsonDoc)
	return len(p), nil
}

func (e *elasticSinkWithIndex) Sync() error { return nil }

func (e *elasticSinkWithIndex) flusher() {
	ticker := time.NewTicker(time.Duration(e.flushSec) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			e.mu.Lock()
			if len(e.buffer) > 0 {
				batch := e.buffer
				e.buffer = make([][]byte, 0, 100)
				e.mu.Unlock()
				for _, doc := range batch {
					_, err := e.client.Index(
						e.index,
						bytes.NewReader(doc),
					)
					if err != nil {
						fmt.Fprintf(os.Stderr, "elasticsearch sink error: %v\n", err)
					}
				}
			} else {
				e.mu.Unlock()
			}
		case <-e.closed:
			return
		}
	}
}

func (e *elasticSinkWithIndex) Close() error {
	close(e.closed)
	return nil
}

// --- zapLogger implementation

type zapLogger struct {
	logger *zap.SugaredLogger
}

func newZapLogger(l *zap.SugaredLogger) Logger {
	return &zapLogger{logger: l}
}

// mergeCategory adds category and subcategory to extra fields
func mergeCategory(cat, sub string, extra map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range extra {
		m[k] = v
	}
	if cat != "" {
		m["category"] = cat
	}
	if sub != "" {
		m["subcategory"] = sub
	}
	return m
}

// flattenMapToFields converts a map[string]interface{} to a flat list of key-value pairs for zap.SugaredLogger.Infow, etc.
func flattenMapToFields(m map[string]interface{}) []interface{} {
	fields := make([]interface{}, 0, len(m)*2)
	for k, v := range m {
		fields = append(fields, k, v)
	}
	return fields
}

// Core logging methods
func (z *zapLogger) Debug(msg string, extra map[string]interface{}) {
	z.logger.Debugw(msg, flattenMapToFields(extra)...)
}
func (z *zapLogger) Info(msg string, extra map[string]interface{}) {
	z.logger.Infow(msg, flattenMapToFields(extra)...)
}
func (z *zapLogger) Warn(msg string, extra map[string]interface{}) {
	z.logger.Warnw(msg, flattenMapToFields(extra)...)
}
func (z *zapLogger) Error(msg string, extra map[string]interface{}) {
	z.logger.Errorw(msg, flattenMapToFields(extra)...)
}
func (z *zapLogger) Fatal(msg string, extra map[string]interface{}) {
	z.logger.Fatalw(msg, flattenMapToFields(extra)...)
}

// Formatted logging methods
func (z *zapLogger) Debugf(template string, args ...interface{}) {
	z.logger.Debugf(template, args...)
}
func (z *zapLogger) Infof(template string, args ...interface{}) {
	z.logger.Infof(template, args...)
}
func (z *zapLogger) Warnf(template string, args ...interface{}) {
	z.logger.Warnf(template, args...)
}
func (z *zapLogger) Errorf(template string, args ...interface{}) {
	z.logger.Errorf(template, args...)
}
func (z *zapLogger) Fatalf(template string, args ...interface{}) {
	z.logger.Fatalf(template, args...)
}

// Context-aware logging methods
func (z *zapLogger) DebugContext(ctx context.Context, msg string, extra map[string]interface{}) {
	z.logger.Debugw(msg, flattenMapToFields(extra)...)
}
func (z *zapLogger) InfoContext(ctx context.Context, msg string, extra map[string]interface{}) {
	z.logger.Infow(msg, flattenMapToFields(extra)...)
}
func (z *zapLogger) WarnContext(ctx context.Context, msg string, extra map[string]interface{}) {
	z.logger.Warnw(msg, flattenMapToFields(extra)...)
}
func (z *zapLogger) ErrorContext(ctx context.Context, msg string, extra map[string]interface{}) {
	z.logger.Errorw(msg, flattenMapToFields(extra)...)
}
func (z *zapLogger) FatalContext(ctx context.Context, msg string, extra map[string]interface{}) {
	z.logger.Fatalw(msg, flattenMapToFields(extra)...)
}

// Category-based logging methods
func (z *zapLogger) DebugWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	z.logger.Debugw(msg, flattenMapToFields(mergeCategory(cat, sub, extra))...)
}
func (z *zapLogger) InfoWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	z.logger.Infow(msg, flattenMapToFields(mergeCategory(cat, sub, extra))...)
}
func (z *zapLogger) WarnWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	z.logger.Warnw(msg, flattenMapToFields(mergeCategory(cat, sub, extra))...)
}
func (z *zapLogger) ErrorWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	z.logger.Errorw(msg, flattenMapToFields(mergeCategory(cat, sub, extra))...)
}
func (z *zapLogger) FatalWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	z.logger.Fatalw(msg, flattenMapToFields(mergeCategory(cat, sub, extra))...)
}
