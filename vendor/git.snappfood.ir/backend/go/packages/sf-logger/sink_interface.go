package sflogger

import "io"

// Sink represents a destination where logs can be sent
type Sink interface {
	// Write sends a log entry to the sink
	Write(entry map[string]interface{}) error

	// Close cleans up resources used by the sink
	Close() error

	// Sync flushes any buffered log entries
	Sync() error
}

// WriteSyncer is an interface that implements io.Writer, Sync
type WriteSyncer interface {
	io.Writer
	Sync() error
}

// Converting a Sink to WriteSyncer for compatibility with zap
func AsSyncWriter(sink Sink) WriteSyncer {
	return &sinkWriter{sink: sink}
}

// sinkWriter adapts a Sink to a WriteSyncer
type sinkWriter struct {
	sink Sink
}

func (w *sinkWriter) Write(p []byte) (n int, err error) {
	// For compatibility with zap, we need to implement Write
	// This is a simplified implementation
	entry := map[string]interface{}{
		"msg": string(p),
	}

	if err := w.sink.Write(entry); err != nil {
		return 0, err
	}

	return len(p), nil
}

func (w *sinkWriter) Sync() error {
	return w.sink.Sync()
}
