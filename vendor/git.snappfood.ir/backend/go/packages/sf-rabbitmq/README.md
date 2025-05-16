# SF RabbitMQ

A Go library for RabbitMQ that provides a simple and flexible way to work with RabbitMQ in your Go applications.

## Features

- Connection management with automatic reconnection
- Support for multiple connections
- Outbox pattern for message retry
- Single and batch message consumption
- JSON message handling
- Exchange and queue declaration
- Configurable retry mechanisms
- Context support for graceful shutdown

## Installation

```bash
go get git.snappfood.ir/backend/go/packages/sf-rabbitmq
```

## Usage

### Basic Setup

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	sfrabbitmq "git.snappfood.ir/backend/go/packages/sf-rabbitmq"
)

// OrderMessage represents a sample message structure
type OrderMessage struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	// Register connections with global and connection-specific options
	err := sfrabbitmq.GetConnector().RegisterConnection(
		// Global options applied to all connections
		sfrabbitmq.WithGlobalOptions(func(c *sfrabbitmq.Config) {
			c.Heartbeat = 30 * time.Second
			c.Dial = sfrabbitmq.DefaultDial(10 * time.Second)
			c.Vhost = "/"
		}),

		// First connection with specific options
		sfrabbitmq.WithConnectionDetails(
			"order-connection",
			"localhost",
			5672,
			"guest",
			"guest",
			sfrabbitmq.WithOptions(func(c *sfrabbitmq.Config) {
				c.Vhost = "/"
			}),
		),

		// Configure outbox retry behavior
		sfrabbitmq.WithOutboxConfig(5*time.Second, 3),

		
		// Declare exchange
		sfrabbitmq.WithDeclareExchange("order-connection", "order_exchange", "direct"),
		// Declare queue
		sfrabbitmq.WithDeclareQueue("order-connection", "order_queue"),
		// Bind queue to exchange
		sfrabbitmq.WithBind("order-connection", "order_queue", "order_routing_key", "order_exchange"),


		// Or Instead of the above three functions, you can only use the following function
		//sfrabbitmq.WithDeclareExchangeAndQueue("order-connection","order_exchange",
		//	"direct","order_queue","order_routing_key"),
	)

	if err != nil {
		log.Fatalf("Failed to register connection: %v", err)
	}

	// Create a context
	ctx := context.Background()

	// Example: Publish messages
	msg := OrderMessage{
		ID:        "123",
		Content:   "Hello, RabbitMQ!",
		Timestamp: time.Now(),
	}

	// Publish using Publish
	body, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
	}

	err = sfrabbitmq.Publish(ctx, "order-connection", "order_exchange", "order_routing_key",
		&sfrabbitmq.Message{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	}

	// Example 1: Single message consumption
	err = sfrabbitmq.Consume(
		ctx,
		"order-connection",
		&sfrabbitmq.ConsumeConfig{
			QueueName:     "order_queue",
			ConsumerName:  "order_consumer",
			AutoAck:       false,
			MaxRetries:    3,
			RetryInterval: 5 * time.Second,
		},
		func(msg *sfrabbitmq.Message) error {
			var receivedMsg OrderMessage
			if err := json.Unmarshal(msg.Body, &receivedMsg); err != nil {
				return fmt.Errorf("failed to unmarshal message: %w", err)
			}
			fmt.Printf("Received single message: %+v\n", receivedMsg)
			return nil
		},
	)

	if err != nil {
		log.Printf("Failed to start single message consumption: %v", err)
	}

	// Example 2: Batch message consumption
	err = sfrabbitmq.BatchConsume(
		ctx,
		"order-connection",
		&sfrabbitmq.BatchConsumeConfig{
			ConsumeConfig: sfrabbitmq.ConsumeConfig{
				QueueName:     "order_queue",
				ConsumerName:  "order_batch_consumer",
				AutoAck:       false,
				MaxRetries:    3,
				RetryInterval: 5 * time.Second,
			},
			ChunkSize: 10,
		},
		func(messages []*sfrabbitmq.Message) error {
			for _, msg := range messages {
				var receivedMsg OrderMessage
				if err := json.Unmarshal(msg.Body, &receivedMsg); err != nil {
					return fmt.Errorf("failed to unmarshal message: %w", err)
				}
				fmt.Printf("Received batch message: %+v\n", receivedMsg)
			}
			return nil
		},
	)

	if err != nil {
		log.Printf("Failed to start batch consumption: %v", err)
	}

	// Keep the program running
	select {}
}
```

### Configuration Options

#### Connection Configuration

```go
type Config struct {
    Host       string
    Port       int
    Username   string
    Password   string
    Vhost      string
    Options    []func(*amqp.Config)
    Heartbeat  time.Duration
    Dial       func(network, addr string) (net.Conn, error)
    FrameSize  int
    Locale     string
    Properties amqp.Table
    ChannelMax uint16
}
```

#### Consume Configuration

```go
type ConsumeConfig struct {
    QueueName     string
    ConsumerName  string
    AutoAck       bool
    Exclusive     bool
    NoLocal       bool
    NoWait        bool
    Args          amqp.Table
    MaxRetries    int           // Maximum number of retries for failed messages
    RetryInterval time.Duration // Time between retries
}

type BatchConsumeConfig struct {
    ConsumeConfig
    ChunkSize int // Number of messages to process in a chunk
}
```

### Message Structure

```go
type Message struct {
    ContentType     string
    ContentEncoding string
    DeliveryMode    uint8
    Priority        uint8
    CorrelationId   string
    ReplyTo         string
    Expiration      string
    MessageId       string
    Timestamp       int64
    Type            string
    UserId          string
    AppId           string
    Body            []byte
    Headers         map[string]interface{}
    DeliveryTag     uint64
}
```

## API Reference

### Connection Management

- `GetConnector()` - Returns a new ServiceConnector instance
- `RegisterConnection(opts ...RegistryOption)` - Registers RabbitMQ connections with provided options
- `WithGlobalOptions(option func(*Config))` - Adds global options for all connections
- `WithConnectionDetails(name, host string, port int, username, password string, opts ...func(*amqp.Config))` - Sets connection details
- `WithOptions(option func(*Config))` - Adds connection-specific options
- `WithOutboxConfig(retryInterval time.Duration, maxRetries int)` - Configures outbox retry behavior
- `DefaultDial(timeout time.Duration)` - Returns a dial function with the specified timeout

### Publishing

- `Publish(ctx context.Context, connName, exchange, routingKey string, msg *Message)` - Publishes a message to an exchange with outbox pattern support

### Consuming

- `Consume(ctx context.Context, connName string, config *ConsumeConfig, handler func(*Message) error)` - Consumes messages one at a time
- `BatchConsume(ctx context.Context, connName string, config *BatchConsumeConfig, handler func([]*Message) error)` - Consumes messages in batches

### Exchange and Queue Management

- `WithDeclareExchange(connName, name, exchangeType string)` - Declares an exchange
- `WithDeclareQueue(connName, name string)` - Declares a queue
- `WithBind(connName, queueName, routingKey, exchangeName string)` - Binds a queue to an exchange

## Best Practices

1. Always use context for graceful shutdown
2. Configure appropriate retry intervals and max retries
3. Use batch consumption for high-volume scenarios
4. Implement proper error handling in message handlers
5. Use the outbox pattern for reliable message delivery
6. Configure appropriate heartbeat and frame size for your use case
7. Use the library's configuration options instead of direct AMQP configuration
8. Use `sfrabbitmq.DefaultDial` for connection timeouts instead of direct AMQP dial functions