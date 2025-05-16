package sflogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// NetworkSinkImpl implements Sink for network destinations
type NetworkSinkImpl struct {
	protocol    Protocol
	address     string
	httpClient  *http.Client
	tcpConn     net.Conn
	udpConn     net.Conn
	transformer func(map[string]interface{}) ([]byte, error)
	mutex       sync.Mutex
}

// NetworkSinkOption defines the function signature for network sink options
type NetworkSinkOption func(*networkSinkConfig)

// networkSinkConfig contains options for configuring a network sink
type networkSinkConfig struct {
	Protocol    Protocol
	Address     string
	Timeout     time.Duration
	Username    string
	Password    string
	Headers     map[string]string
	Transformer func(map[string]interface{}) ([]byte, error)
}

// NewNetworkSink creates a new network sink
func NewNetworkSink(opts ...NetworkSinkOption) (Sink, error) {
	config := &networkSinkConfig{
		Protocol:    HTTP,
		Timeout:     time.Second * 5,
		Headers:     make(map[string]string),
		Transformer: jsonTransformer,
	}

	// Apply options
	for _, opt := range opts {
		opt(config)
	}

	sink := &NetworkSinkImpl{
		protocol:    config.Protocol,
		address:     config.Address,
		transformer: config.Transformer,
	}

	// Configure based on protocol
	var err error
	switch config.Protocol {
	case HTTP, HTTPS:
		sink.httpClient = &http.Client{
			Timeout: config.Timeout,
		}
	case TCP:
		sink.tcpConn, err = net.Dial("tcp", config.Address)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to TCP destination: %w", err)
		}
	case UDP:
		sink.udpConn, err = net.Dial("udp", config.Address)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to UDP destination: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported protocol: %s", config.Protocol)
	}

	return sink, nil
}

// Write sends a log entry to the network destination
func (s *NetworkSinkImpl) Write(entry map[string]interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Transform the entry to bytes
	data, err := s.transformer(entry)
	if err != nil {
		return fmt.Errorf("failed to transform log entry: %w", err)
	}

	// Send based on protocol
	switch s.protocol {
	case HTTP, HTTPS:
		req, err := http.NewRequest("POST", s.address, bytes.NewReader(data))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		_, err = s.httpClient.Do(req)
		return err
	case TCP:
		if s.tcpConn == nil {
			return fmt.Errorf("TCP connection is not established")
		}
		_, err := s.tcpConn.Write(append(data, '\n'))
		return err
	case UDP:
		if s.udpConn == nil {
			return fmt.Errorf("UDP connection is not established")
		}
		_, err := s.udpConn.Write(data)
		return err
	default:
		return fmt.Errorf("unsupported protocol: %s", s.protocol)
	}
}

// Close closes the network sink connections
func (s *NetworkSinkImpl) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var err error
	if s.tcpConn != nil {
		err = s.tcpConn.Close()
		s.tcpConn = nil
	}
	if s.udpConn != nil {
		udpErr := s.udpConn.Close()
		if err == nil {
			err = udpErr
		}
		s.udpConn = nil
	}
	return err
}

// Sync is a no-op for network sinks
func (s *NetworkSinkImpl) Sync() error {
	return nil
}

// Default transformer functions
func jsonTransformer(entry map[string]interface{}) ([]byte, error) {
	return json.Marshal(entry)
}

// NetworkWithProtocol sets the protocol for the network sink
func NetworkWithProtocol(protocol Protocol) NetworkSinkOption {
	return func(c *networkSinkConfig) {
		c.Protocol = protocol
	}
}

// NetworkWithAddress sets the address for the network sink
func NetworkWithAddress(address string) NetworkSinkOption {
	return func(c *networkSinkConfig) {
		c.Address = address
	}
}

// NetworkWithTimeout sets the timeout for network requests
func NetworkWithTimeout(timeout time.Duration) NetworkSinkOption {
	return func(c *networkSinkConfig) {
		c.Timeout = timeout
	}
}

// NetworkWithCredentials sets the credentials for the network sink
func NetworkWithCredentials(username, password string) NetworkSinkOption {
	return func(c *networkSinkConfig) {
		c.Username = username
		c.Password = password
	}
}

// NetworkWithHeader adds a header for HTTP requests
func NetworkWithHeader(key, value string) NetworkSinkOption {
	return func(c *networkSinkConfig) {
		c.Headers[key] = value
	}
}

// NetworkWithTransformer sets a custom transformer for log entries
func NetworkWithTransformer(transformer func(map[string]interface{}) ([]byte, error)) NetworkSinkOption {
	return func(c *networkSinkConfig) {
		c.Transformer = transformer
	}
}

// Register network sinks with the registry
func init() {
	// Register common network sink types

	// Graylog sink
	RegisterSink("graylog", func(u *url.URL) (Sink, error) {
		var opts []NetworkSinkOption

		// Default to UDP protocol for Graylog
		protocol := UDP
		if u.Scheme == "graylog+tcp" {
			protocol = TCP
		} else if u.Scheme == "graylog+http" || u.Scheme == "graylog+https" {
			if u.Scheme == "graylog+https" {
				protocol = HTTPS
			} else {
				protocol = HTTP
			}
		}
		opts = append(opts, NetworkWithProtocol(protocol))

		// Set address
		address := u.Host
		if protocol == UDP || protocol == TCP {
			// Make sure port is specified
			if !strings.Contains(address, ":") {
				address = address + ":12201" // Default Graylog port
			}
		} else {
			// For HTTP(S), use the full URL
			address = fmt.Sprintf("%s://%s%s", strings.TrimPrefix(protocol.String(), "graylog+"), u.Host, u.Path)
		}
		opts = append(opts, NetworkWithAddress(address))

		// Parse query parameters
		q := u.Query()
		if timeout := q.Get("timeout"); timeout != "" {
			if d, err := time.ParseDuration(timeout); err == nil {
				opts = append(opts, NetworkWithTimeout(d))
			}
		}

		// Add Graylog transformer
		opts = append(opts, NetworkWithTransformer(func(entry map[string]interface{}) ([]byte, error) {
			// Convert to GELF format
			hostname, _ := entry["host"].(string)
			if hostname == "" {
				hostname = "unknown"
			}

			gelf := map[string]interface{}{
				"version":       "1.1",
				"host":          hostname,
				"short_message": entry["msg"],
				"timestamp":     time.Now().Unix(),
			}

			// Copy all fields to GELF
			for k, v := range entry {
				if k != "msg" {
					if k == "level" {
						// Map log levels to GELF levels
						switch strings.ToUpper(fmt.Sprintf("%v", v)) {
						case "DEBUG":
							gelf["level"] = 7
						case "INFO":
							gelf["level"] = 6
						case "WARN", "WARNING":
							gelf["level"] = 4
						case "ERROR":
							gelf["level"] = 3
						case "FATAL":
							gelf["level"] = 2
						default:
							gelf["level"] = 6
						}
					} else {
						gelf["_"+k] = v
					}
				}
			}

			return json.Marshal(gelf)
		}))

		return NewNetworkSink(opts...)
	})

	// Logstash sink
	RegisterSink("logstash", func(u *url.URL) (Sink, error) {
		var opts []NetworkSinkOption

		// Default to TCP protocol for Logstash
		protocol := TCP
		if u.Scheme == "logstash+udp" {
			protocol = UDP
		} else if u.Scheme == "logstash+http" || u.Scheme == "logstash+https" {
			if u.Scheme == "logstash+https" {
				protocol = HTTPS
			} else {
				protocol = HTTP
			}
		}
		opts = append(opts, NetworkWithProtocol(protocol))

		// Set address
		address := u.Host
		if protocol == UDP || protocol == TCP {
			// Make sure port is specified
			if !strings.Contains(address, ":") {
				address = address + ":5044" // Default Logstash port
			}
		} else {
			// For HTTP(S), use the full URL
			address = fmt.Sprintf("%s://%s%s", strings.TrimPrefix(protocol.String(), "logstash+"), u.Host, u.Path)
		}
		opts = append(opts, NetworkWithAddress(address))

		// Parse query parameters
		q := u.Query()
		if timeout := q.Get("timeout"); timeout != "" {
			if d, err := time.ParseDuration(timeout); err == nil {
				opts = append(opts, NetworkWithTimeout(d))
			}
		}

		return NewNetworkSink(opts...)
	})

	// Loki sink
	RegisterSink("loki", func(u *url.URL) (Sink, error) {
		var opts []NetworkSinkOption

		// Loki only supports HTTP/HTTPS
		protocol := HTTP
		if u.Scheme == "loki+https" {
			protocol = HTTPS
		}
		opts = append(opts, NetworkWithProtocol(protocol))

		// Set address
		address := fmt.Sprintf("%s://%s%s", strings.TrimPrefix(protocol.String(), "loki+"), u.Host, u.Path)
		if !strings.HasSuffix(address, "/loki/api/v1/push") {
			if !strings.HasSuffix(address, "/") {
				address += "/"
			}
			address += "loki/api/v1/push"
		}
		opts = append(opts, NetworkWithAddress(address))

		// Parse query parameters
		q := u.Query()
		if timeout := q.Get("timeout"); timeout != "" {
			if d, err := time.ParseDuration(timeout); err == nil {
				opts = append(opts, NetworkWithTimeout(d))
			}
		}

		// Add authorization header if provided
		if apiKey := q.Get("apiKey"); apiKey != "" {
			opts = append(opts, NetworkWithHeader("X-Scope-OrgID", apiKey))
		}

		// Add Loki transformer
		opts = append(opts, NetworkWithTransformer(func(entry map[string]interface{}) ([]byte, error) {
			// Extract timestamp
			ts := time.Now().UnixNano()
			if tsStr, ok := entry["timestamp"].(string); ok {
				if t, err := time.Parse(time.RFC3339, tsStr); err == nil {
					ts = t.UnixNano()
				}
			}

			// Create labels
			labels := make(map[string]string)
			labels["level"] = fmt.Sprintf("%v", entry["level"])
			if app, ok := entry["AppName"].(string); ok {
				labels["app"] = app
			}

			// Create Loki stream
			jsonData, err := json.Marshal(entry)
			if err != nil {
				return nil, err
			}

			stream := map[string]interface{}{
				"stream": labels,
				"values": [][]string{
					{fmt.Sprintf("%d", ts), string(jsonData)},
				},
			}

			lokiRequest := map[string]interface{}{
				"streams": []interface{}{stream},
			}

			return json.Marshal(lokiRequest)
		}))

		return NewNetworkSink(opts...)
	})
}
