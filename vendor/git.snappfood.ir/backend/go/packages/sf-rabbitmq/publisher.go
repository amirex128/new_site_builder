package sfrabbitmq

import (
	"context"
	"fmt"
)

// Publish publishes a message to an exchange with outbox pattern support
func Publish(ctx context.Context, connName, exchange, routingKey string, msg *Message) error {
	if globalRegistry.logger != nil {
		extras := make(map[string]interface{})
		extras["connection_name"] = connName
		extras["exchange"] = exchange
		extras["routing_key"] = routingKey
		extras["message_id"] = msg.MessageId
		extras["content_type"] = msg.ContentType
		globalRegistry.logger.InfoWithCategory(Category.API.Messaging, SubCategory.Operation.Initialization, "Publishing message", extras)
	}

	conn, err := getConnection(connName)
	if err != nil {
		if globalRegistry.logger != nil {
			extras := make(map[string]interface{})
			extras["connection_name"] = connName
			extras["exchange"] = exchange
			extras["routing_key"] = routingKey
			extras[ExtraKey.Error.ErrorMessage] = err.Error()
			globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to get connection for publishing", extras)
		}
		return fmt.Errorf("failed to get connection: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		if globalRegistry.logger != nil {
			extras := make(map[string]interface{})
			extras["connection_name"] = connName
			extras["exchange"] = exchange
			extras["routing_key"] = routingKey
			extras[ExtraKey.Error.ErrorMessage] = err.Error()
			globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to create channel for publishing", extras)
		}
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer ch.Close()

	err = ch.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		msg.ToAMQP(),
	)

	if err != nil {
		if globalRegistry.logger != nil {
			extras := make(map[string]interface{})
			extras["connection_name"] = connName
			extras["exchange"] = exchange
			extras["routing_key"] = routingKey
			extras["message_id"] = msg.MessageId
			extras[ExtraKey.Error.ErrorMessage] = err.Error()
			globalRegistry.logger.WarnWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to publish message, adding to outbox", extras)
		}

		// Add to outbox for retry
		globalRegistry.AddToOutbox(&OutboxMessage{
			ConnectionName: connName,
			Exchange:       exchange,
			RoutingKey:     routingKey,
			Mandatory:      false,
			Immediate:      false,
			Message:        msg,
			RetryCount:     0,
			LastError:      err,
		})
		return fmt.Errorf("failed to publish message, added to outbox: %w", err)
	}

	if globalRegistry.logger != nil {
		extras := make(map[string]interface{})
		extras["connection_name"] = connName
		extras["exchange"] = exchange
		extras["routing_key"] = routingKey
		extras["message_id"] = msg.MessageId
		globalRegistry.logger.InfoWithCategory(Category.API.Messaging, SubCategory.Status.Success, "Successfully published message", extras)
	}

	return nil
}
