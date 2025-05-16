package sfrabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Client represents a RabbitMQ client with connection and channel
type Client struct {
	ctx     context.Context
	conn    *amqp.Connection
	channel *amqp.Channel
	name    string
}

// MustClient creates a new RabbitMQ client from a registered connection or panics
func MustClient(ctx context.Context, name string) *Client {
	if globalRegistry.logger != nil {
		extras := make(map[string]interface{})
		extras["connection_name"] = name
		globalRegistry.logger.InfoWithCategory(Category.Infrastructure.Network, SubCategory.Operation.Initialization, "Creating client", extras)
	}

	conn, err := getConnection(name)
	if err != nil {
		if globalRegistry.logger != nil {
			extras := make(map[string]interface{})
			extras["connection_name"] = name
			extras[ExtraKey.Error.ErrorMessage] = err.Error()
			globalRegistry.logger.FatalWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to get connection for client", extras)
		}
		panic(fmt.Sprintf("failed to get connection '%s': %v", name, err))
	}

	channel, err := conn.Channel()
	if err != nil {
		if globalRegistry.logger != nil {
			extras := make(map[string]interface{})
			extras["connection_name"] = name
			extras[ExtraKey.Error.ErrorMessage] = err.Error()
			globalRegistry.logger.FatalWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to create channel for client", extras)
		}
		panic(fmt.Sprintf("failed to create channel for connection '%s': %v", name, err))
	}

	if globalRegistry.logger != nil {
		extras := make(map[string]interface{})
		extras["connection_name"] = name
		globalRegistry.logger.InfoWithCategory(Category.Infrastructure.Network, SubCategory.Status.Success, "Successfully created client", extras)
	}

	return &Client{
		ctx:     ctx,
		conn:    conn,
		channel: channel,
		name:    name,
	}
}

// Close closes the client's connection and channel
func (c *Client) Close() error {
	if c.channel == nil || c.conn == nil {
		return nil
	}

	if globalRegistry.logger != nil {
		extras := make(map[string]interface{})
		extras["connection_name"] = c.name
		globalRegistry.logger.InfoWithCategory(Category.Infrastructure.Network, SubCategory.Operation.Shutdown, "Closing client", extras)
	}

	var finalErr error

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			if globalRegistry.logger != nil {
				extras := make(map[string]interface{})
				extras["connection_name"] = c.name
				extras[ExtraKey.Error.ErrorMessage] = err.Error()
				globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to close channel", extras)
			}
			finalErr = fmt.Errorf("failed to close channel: %w", err)
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			if globalRegistry.logger != nil {
				extras := make(map[string]interface{})
				extras["connection_name"] = c.name
				extras[ExtraKey.Error.ErrorMessage] = err.Error()
				globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to close connection", extras)
			}
			if finalErr == nil {
				finalErr = fmt.Errorf("failed to close connection: %w", err)
			} else {
				finalErr = fmt.Errorf("%v; and failed to close connection: %w", finalErr, err)
			}
		}
	}

	if finalErr == nil && globalRegistry.logger != nil {
		extras := make(map[string]interface{})
		extras["connection_name"] = c.name
		globalRegistry.logger.InfoWithCategory(Category.Infrastructure.Network, SubCategory.Status.Success, "Successfully closed client", extras)
	}

	return finalErr
}

// IsClosed checks if the client's connection is closed
func (c *Client) IsClosed() bool {
	return c.conn == nil || c.conn.IsClosed()
}

// Connection returns the underlying AMQP connection
func (c *Client) Connection() *amqp.Connection {
	return c.conn
}

// Channel returns the underlying AMQP channel
func (c *Client) Channel() *amqp.Channel {
	return c.channel
}

// Name returns the client's name
func (c *Client) Name() string {
	return c.name
}
