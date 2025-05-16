package sfrabbitmq

import (
	"net"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Config represents RabbitMQ connection configuration
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

// Message represents a RabbitMQ message
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

// ExchangeConfig represents exchange configuration
type ExchangeConfig struct {
	Name       string
	Type       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

// QueueConfig represents queue configuration
type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

// BindConfig represents binding configuration
type BindConfig struct {
	QueueName    string
	RoutingKey   string
	ExchangeName string
	NoWait       bool
	Args         amqp.Table
}

// OutboxMessage represents a message that failed to publish
type OutboxMessage struct {
	ConnectionName string
	Exchange       string
	RoutingKey     string
	Mandatory      bool
	Immediate      bool
	Message        *Message
	RetryCount     int
	LastError      error
	NextRetry      time.Time
}

// ConsumeConfig represents configuration for single message consumption
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

// BatchConsumeConfig represents configuration for batch message consumption
type BatchConsumeConfig struct {
	ConsumeConfig
	ChunkSize int // Number of messages to process in a chunk
}

// ExchangeQueueConfig represents configuration for exchange and queue declaration
type ExchangeQueueConfig struct {
	ExchangeName string
	ExchangeType string
	QueueName    string
	RoutingKey   string
	Durable      bool
	AutoDelete   bool
	Internal     bool
	Exclusive    bool
	NoWait       bool
	Args         amqp.Table
}

// Declaration represents a RabbitMQ declaration (exchange, queue, or binding)
type Declaration struct {
	Type         string
	ConnName     string
	Exchange     string
	Queue        string
	RoutingKey   string
	ExchangeType string
	Durable      bool
	AutoDelete   bool
	Internal     bool
	Exclusive    bool
	NoWait       bool
	Args         amqp.Table
}

// ToAMQP converts Config to amqp.Config
func (c *Config) ToAMQP() amqp.Config {
	config := amqp.Config{
		Vhost:      c.Vhost,
		Heartbeat:  c.Heartbeat,
		Dial:       c.Dial,
		FrameSize:  c.FrameSize,
		Locale:     c.Locale,
		Properties: c.Properties,
		ChannelMax: c.ChannelMax,
	}
	for _, opt := range c.Options {
		opt(&config)
	}
	return config
}

// ToAMQP converts Message to amqp.Publishing
func (m *Message) ToAMQP() amqp.Publishing {
	return amqp.Publishing{
		ContentType:     m.ContentType,
		ContentEncoding: m.ContentEncoding,
		DeliveryMode:    m.DeliveryMode,
		Priority:        m.Priority,
		CorrelationId:   m.CorrelationId,
		ReplyTo:         m.ReplyTo,
		Expiration:      m.Expiration,
		MessageId:       m.MessageId,
		Timestamp:       time.Unix(m.Timestamp, 0),
		Type:            m.Type,
		UserId:          m.UserId,
		AppId:           m.AppId,
		Body:            m.Body,
		Headers:         m.Headers,
	}
}

// FromAMQP creates Message from amqp.Delivery
func FromAMQP(delivery amqp.Delivery) *Message {
	return &Message{
		ContentType:     delivery.ContentType,
		ContentEncoding: delivery.ContentEncoding,
		DeliveryMode:    delivery.DeliveryMode,
		Priority:        delivery.Priority,
		CorrelationId:   delivery.CorrelationId,
		ReplyTo:         delivery.ReplyTo,
		Expiration:      delivery.Expiration,
		MessageId:       delivery.MessageId,
		Timestamp:       delivery.Timestamp.Unix(),
		Type:            delivery.Type,
		UserId:          delivery.UserId,
		AppId:           delivery.AppId,
		Body:            delivery.Body,
		Headers:         delivery.Headers,
	}
}
