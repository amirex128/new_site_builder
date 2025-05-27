package sflogger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/natefinch/lumberjack"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

type FormatterType int

const (
	ColoredTextFormatter FormatterType = iota
	JSONFormatter
)

type Option func(*Registry)

type Registry struct {
	mu           sync.RWMutex
	loggerEngine *zap.Logger
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

var globalRegistry = &Registry{}

// --------------------------------------------
// Logger registration (no big changes here)
// --------------------------------------------
func RegisterLogger(opts ...Option) Logger {
	for _, opt := range opts {
		opt(globalRegistry)
	}
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.Level(globalRegistry.level)),
		Development:       false,
		Encoding:          "console",
		EncoderConfig:     defaultEncoderConfig(),
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		InitialFields:     map[string]interface{}{"appName": globalRegistry.appName},
		DisableStacktrace: !globalRegistry.stacktrace,
	}

	if globalRegistry.formatter == JSONFormatter {
		cfg.Encoding = "json"
		cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	}

	var cores []zapcore.Core

	consoleSyncer := zapcore.Lock(os.Stdout)
	encoder := getEncoder(cfg)
	cores = append(cores, zapcore.NewCore(encoder, consoleSyncer, zap.NewAtomicLevelAt(zapcore.Level(globalRegistry.level))))

	if globalRegistry.fileSinkCfg != nil {
		fs, err := newFileSink(globalRegistry.fileSinkCfg)
		if err == nil {
			fileCore := zapcore.NewCore(encoder, fs, zap.NewAtomicLevelAt(zapcore.Level(globalRegistry.level)))
			cores = append(cores, fileCore)
		} else {
			fmt.Fprintf(os.Stderr, "file sink error: %v\n", err)
		}
	}

	if globalRegistry.mongoSinkCfg != nil {
		mongoCfg := cfg
		mongoCfg.Encoding = "json"
		mongoCfg.EncoderConfig = defaultEncoderConfig()
		mongoCfg.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		mongoEncoder := getEncoder(mongoCfg)

		ms, err := newMongoSink(globalRegistry.mongoSinkCfg)
		if err == nil {
			mongoCore := zapcore.NewCore(mongoEncoder, ms, zap.NewAtomicLevelAt(zapcore.Level(globalRegistry.level)))
			cores = append(cores, mongoCore)
		} else {
			fmt.Fprintf(os.Stderr, "mongo sink error: %v\n", err)
		}
	}

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
	logger := zap.New(core, zap.AddCaller(),zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	globalRegistry.loggerEngine = logger

	return newZapLogger(globalRegistry.loggerEngine)
}

func defaultEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
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
	}
}

func getEncoder(cfg zap.Config) zapcore.Encoder {
	if cfg.Encoding == "json" {
		return zapcore.NewJSONEncoder(cfg.EncoderConfig)
	}
	return zapcore.NewConsoleEncoder(cfg.EncoderConfig)
}

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

// ---------------------------------
// File Sink with background flush
// ---------------------------------
type fileSink struct {
	writer *lumberjack.Logger
	buffer chan []byte
	closed chan struct{}
	wg     sync.WaitGroup
}

func newFileSink(cfg *fileSinkConfig) (zapcore.WriteSyncer, error) {
	writer := &lumberjack.Logger{
		Filename:   cfg.path,
		MaxSize:    cfg.maxSizeMB,
		MaxAge:     cfg.maxAgeDays,
		MaxBackups: cfg.maxBackups,
		Compress:   cfg.compress,
	}

	fs := &fileSink{
		writer: writer,
		buffer: make(chan []byte, 1000),
		closed: make(chan struct{}),
	}
	fs.wg.Add(1)
	go fs.backgroundWriter()

	return zapcore.AddSync(fs), nil
}

func (f *fileSink) Write(p []byte) (n int, err error) {
	select {
	case f.buffer <- append([]byte{}, p...): // copy p before enqueue
		return len(p), nil
	default:
		// buffer full, drop log or handle differently
		return 0, fmt.Errorf("file sink buffer full")
	}
}

func (f *fileSink) Sync() error {
	// Flush remaining logs synchronously
	close(f.closed)
	f.wg.Wait()
	return nil
}

func (f *fileSink) backgroundWriter() {
	defer f.wg.Done()
	for {
		select {
		case p := <-f.buffer:
			_, _ = f.writer.Write(p)
		case <-f.closed:
			// Flush remaining items
			for {
				select {
				case p := <-f.buffer:
					_, _ = f.writer.Write(p)
				default:
					return
				}
			}
		}
	}
}

// ---------------------------------
// Mongo Sink with background flush
// ---------------------------------
type mongoSink struct {
	client     *mongo.Client
	collection *mongo.Collection
	buffer     chan map[string]interface{}
	closed     chan struct{}
	wg         sync.WaitGroup
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
		buffer:     make(chan map[string]interface{}, 1000),
		closed:     make(chan struct{}),
	}
	ms.wg.Add(1)
	go ms.backgroundInserter()

	return ms, nil
}

func (m *mongoSink) Write(p []byte) (n int, err error) {
	var doc map[string]interface{}
	if err := json.Unmarshal(p, &doc); err != nil {
		doc = map[string]interface{}{"raw_log": string(p)}
	}
	if _, ok := doc["time"]; !ok {
		doc["timestamp"] = time.Now()
	}
	select {
	case m.buffer <- doc:
		return len(p), nil
	default:
		// buffer full: drop or handle error
		return 0, fmt.Errorf("mongo sink buffer full")
	}
}

func (m *mongoSink) Sync() error {
	close(m.closed)
	m.wg.Wait()
	return m.client.Disconnect(context.Background())
}

func (m *mongoSink) backgroundInserter() {
	defer m.wg.Done()

	batch := make([]map[string]interface{}, 0, 100)
	flushTicker := time.NewTicker(time.Second * time.Duration(5))
	defer flushTicker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if _, err := m.collection.InsertMany(ctx, toInterfaceSlice(batch)); err != nil {
			fmt.Fprintf(os.Stderr, "mongo sink insert error: %v\n", err)
		}
		batch = batch[:0]
	}

	for {
		select {
		case doc := <-m.buffer:
			batch = append(batch, doc)
			if len(batch) >= 100 {
				flush()
			}
		case <-flushTicker.C:
			flush()
		case <-m.closed:
			flush()
			return
		}
	}
}

func toInterfaceSlice(maps []map[string]interface{}) []interface{} {
	out := make([]interface{}, len(maps))
	for i := range maps {
		out[i] = maps[i]
	}
	return out
}

// ---------------------------------
// Elasticsearch Sink with background flush
// ---------------------------------

type elasticSink struct {
	client   *elasticsearch.Client
	buffer   chan []byte
	closed   chan struct{}
	wg       sync.WaitGroup
	index    string
	flushSec int
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
		buffer:   make(chan []byte, 1000),
		closed:   make(chan struct{}),
		index:    cfg.index,
		flushSec: cfg.flushSec,
	}
	es.wg.Add(1)
	go es.backgroundIndexer()
	return es, nil
}

func (e *elasticSink) Write(p []byte) (n int, err error) {
	select {
	case e.buffer <- append([]byte{}, p...):
		return len(p), nil
	default:
		return 0, fmt.Errorf("elasticsearch sink buffer full")
	}
}

func (e *elasticSink) Sync() error {
	close(e.closed)
	e.wg.Wait()
	return nil
}

func (e *elasticSink) backgroundIndexer() {
	defer e.wg.Done()
	batch := make([][]byte, 0, 100)
	flushTicker := time.NewTicker(time.Second * time.Duration(e.flushSec))
	defer flushTicker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		for _, doc := range batch {
			_, err := e.client.Index(e.index, bytes.NewReader(doc))
			if err != nil {
				fmt.Fprintf(os.Stderr, "elasticsearch sink error: %v\n", err)
			}
		}
		batch = batch[:0]
	}

	for {
		select {
		case doc := <-e.buffer:
			batch = append(batch, doc)
			if len(batch) >= 100 {
				flush()
			}
		case <-flushTicker.C:
			flush()
		case <-e.closed:
			flush()
			return
		}
	}
}

// ---------------------------------
// zapLogger implementation
// ---------------------------------

type zapLogger struct {
	logger *zap.Logger
}

func newZapLogger(l *zap.Logger) Logger {
	return &zapLogger{logger: l}
}

func mapToZapFields(m map[string]interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(m))
	for k, v := range m {
		switch val := v.(type) {
		case string:
			fields = append(fields, zap.String(k, val))
		case int:
			fields = append(fields, zap.Int(k, val))
		case int64:
			fields = append(fields, zap.Int64(k, val))
		case float64:
			fields = append(fields, zap.Float64(k, val))
		case bool:
			fields = append(fields, zap.Bool(k, val))
		case time.Time:
			fields = append(fields, zap.Time(k, val))
		case []byte:
			fields = append(fields, zap.ByteString(k, val))
		case []interface{}:
			fields = append(fields, zap.Array(k, zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
				for _, v := range val {
					enc.AppendReflected(v)
				}
				return nil
			})))
		case map[string]interface{}:
			fields = append(fields, zap.Object(k, zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
				for kk, vv := range val {
					enc.AddReflected(kk, vv)
				}
				return nil
			})))
		case error:
			fields = append(fields, zap.Error(val))
		default:
			fields = append(fields, zap.Any(k, val))
		}
	}
	return fields
}

func (z *zapLogger) Debug(msg string, extra map[string]interface{}) {
	z.logger.Debug(msg, mapToZapFields(extra)...)
}
func (z *zapLogger) Info(msg string, extra map[string]interface{}) {
	z.logger.Info(msg, mapToZapFields(extra)...)
}
func (z *zapLogger) Warn(msg string, extra map[string]interface{}) {
	z.logger.Warn(msg, mapToZapFields(extra)...)
}
func (z *zapLogger) Error(msg string, extra map[string]interface{}) {
	z.logger.Error(msg, mapToZapFields(extra)...)
}
func (z *zapLogger) Fatal(msg string, extra map[string]interface{}) {
	z.logger.Fatal(msg, mapToZapFields(extra)...)
}

func (z *zapLogger) Debugf(template string, args ...interface{}) {
	z.logger.Sugar().Debugf(template, args...)
}
func (z *zapLogger) Infof(template string, args ...interface{}) {
	z.logger.Sugar().Infof(template, args...)
}
func (z *zapLogger) Warnf(template string, args ...interface{}) {
	z.logger.Sugar().Warnf(template, args...)
}
func (z *zapLogger) Errorf(template string, args ...interface{}) {
	z.logger.Sugar().Errorf(template, args...)
}
func (z *zapLogger) Fatalf(template string, args ...interface{}) {
	z.logger.Sugar().Fatalf(template, args...)
}

func (z *zapLogger) DebugContext(ctx context.Context, msg string, extra map[string]interface{}) {
	z.logger.Debug(msg, mapToZapFields(extra)...)
}
func (z *zapLogger) InfoContext(ctx context.Context, msg string, extra map[string]interface{}) {
	z.logger.Info(msg, mapToZapFields(extra)...)
}
func (z *zapLogger) WarnContext(ctx context.Context, msg string, extra map[string]interface{}) {
	z.logger.Warn(msg, mapToZapFields(extra)...)
}
func (z *zapLogger) ErrorContext(ctx context.Context, msg string, extra map[string]interface{}) {
	z.logger.Error(msg, mapToZapFields(extra)...)
}
func (z *zapLogger) FatalContext(ctx context.Context, msg string, extra map[string]interface{}) {
	z.logger.Fatal(msg, mapToZapFields(extra)...)
}

func mergeCategory(cat, sub string, extra map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(extra)+2)
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

func (z *zapLogger) DebugWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	z.logger.Debug(msg, mapToZapFields(mergeCategory(cat, sub, extra))...)
}
func (z *zapLogger) InfoWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	z.logger.Info(msg, mapToZapFields(mergeCategory(cat, sub, extra))...)
}
func (z *zapLogger) WarnWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	z.logger.Warn(msg, mapToZapFields(mergeCategory(cat, sub, extra))...)
}
func (z *zapLogger) ErrorWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	z.logger.Error(msg, mapToZapFields(mergeCategory(cat, sub, extra))...)
}
func (z *zapLogger) FatalWithCategory(cat string, sub string, msg string, extra map[string]interface{}) {
	z.logger.Fatal(msg, mapToZapFields(mergeCategory(cat, sub, extra))...)
}
