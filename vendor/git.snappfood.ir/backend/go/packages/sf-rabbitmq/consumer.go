package sfrabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageHandler is a function for processing messages
type MessageHandler func(ctx context.Context, msg *Message) error

// Consume consumes messages from a queue with retry support
func Consume(ctx context.Context, connName, queueName string, handler MessageHandler, config *ConsumeConfig) error {
	if globalRegistry.logger != nil {
		extras := make(map[string]interface{})
		extras["connection_name"] = connName
		extras["queue_name"] = queueName
		extras["consumer_name"] = config.ConsumerName
		globalRegistry.logger.InfoWithCategory(Category.API.Messaging, SubCategory.Operation.Initialization, "Starting to consume messages", extras)
	}

	conn, err := getConnection(connName)
	if err != nil {
		if globalRegistry.logger != nil {
			extras := make(map[string]interface{})
			extras["connection_name"] = connName
			extras[ExtraKey.Error.ErrorMessage] = err.Error()
			globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to get connection for consuming", extras)
		}
		return fmt.Errorf("failed to get connection: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		if globalRegistry.logger != nil {
			extras := make(map[string]interface{})
			extras["connection_name"] = connName
			extras[ExtraKey.Error.ErrorMessage] = err.Error()
			globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to create channel for consuming", extras)
		}
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer ch.Close()

	deliveries, err := ch.Consume(
		queueName,
		config.ConsumerName,
		config.AutoAck,
		config.Exclusive,
		config.NoLocal,
		config.NoWait,
		config.Args,
	)

	if err != nil {
		if globalRegistry.logger != nil {
			extras := make(map[string]interface{})
			extras["connection_name"] = connName
			extras["queue_name"] = queueName
			extras["consumer_name"] = config.ConsumerName
			extras[ExtraKey.Error.ErrorMessage] = err.Error()
			globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to start consuming", extras)
		}
		return fmt.Errorf("failed to consume: %w", err)
	}

	for {
		if globalRegistry.logger != nil {
			extras := make(map[string]interface{})
			extras["connection_name"] = connName
			extras["queue_name"] = queueName
			globalRegistry.logger.DebugWithCategory(Category.API.Messaging, SubCategory.Status.Debug, "Waiting for message", extras)
		}

		select {
		case <-ctx.Done():
			if globalRegistry.logger != nil {
				extras := make(map[string]interface{})
				extras["connection_name"] = connName
				extras["queue_name"] = queueName
				globalRegistry.logger.InfoWithCategory(Category.API.Messaging, SubCategory.Operation.Shutdown, "Consumer context cancelled", extras)
			}
			return nil
		case d, ok := <-deliveries:
			if !ok {
				if globalRegistry.logger != nil {
					extras := make(map[string]interface{})
					extras["connection_name"] = connName
					extras["queue_name"] = queueName
					globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Delivery channel closed", extras)
				}
				return fmt.Errorf("delivery channel closed")
			}

			if globalRegistry.logger != nil {
				extras := make(map[string]interface{})
				extras["connection_name"] = connName
				extras["queue_name"] = queueName
				extras["message_id"] = d.MessageId
				globalRegistry.logger.DebugWithCategory(Category.API.Messaging, SubCategory.Status.Debug, "Received message", extras)
			}

			// Process with retry support
			var processErr error
			retries := 0

			for retries <= config.MaxRetries {
				if retries > 0 && globalRegistry.logger != nil {
					extras := make(map[string]interface{})
					extras["connection_name"] = connName
					extras["queue_name"] = queueName
					extras["message_id"] = d.MessageId
					extras["retry_count"] = retries
					extras["max_retries"] = config.MaxRetries
					globalRegistry.logger.WarnWithCategory(Category.Error.Retry, SubCategory.Status.Retry, "Retrying message processing", extras)
				}

				processErr = handler(ctx, FromAMQP(d))
				if processErr == nil {
					break
				}

				retries++
				if retries <= config.MaxRetries {
					time.Sleep(config.RetryInterval)
				}
			}

			if !config.AutoAck {
				if processErr == nil {
					if err := d.Ack(false); err != nil {
						if globalRegistry.logger != nil {
							extras := make(map[string]interface{})
							extras["connection_name"] = connName
							extras["queue_name"] = queueName
							extras["message_id"] = d.MessageId
							extras[ExtraKey.Error.ErrorMessage] = err.Error()
							globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to ack message", extras)
						}
						return fmt.Errorf("failed to ack message: %w", err)
					}

					if globalRegistry.logger != nil {
						extras := make(map[string]interface{})
						extras["connection_name"] = connName
						extras["queue_name"] = queueName
						extras["message_id"] = d.MessageId
						globalRegistry.logger.InfoWithCategory(Category.API.Messaging, SubCategory.Status.Success, "Successfully processed message", extras)
					}
				} else {
					// Max retries reached, reject the message
					if globalRegistry.logger != nil {
						extras := make(map[string]interface{})
						extras["connection_name"] = connName
						extras["queue_name"] = queueName
						extras["message_id"] = d.MessageId
						extras["retry_count"] = retries
						extras["max_retries"] = config.MaxRetries
						extras[ExtraKey.Error.ErrorMessage] = processErr.Error()
						globalRegistry.logger.ErrorWithCategory(Category.Error.Error, SubCategory.Status.Error, "Failed to process message after max retries", extras)
					}

					return processErr
				}
			}
		}
	}
}

// BatchConsume starts consuming messages in batches from a queue
func BatchConsume(ctx context.Context, connName string, config *BatchConsumeConfig, handler func([]*Message) error) error {
	if globalRegistry.logger != nil {
		extras := make(map[string]interface{})
		extras["connection_name"] = connName
		extras["queue_name"] = config.QueueName
		extras["consumer_name"] = config.ConsumerName
		extras["chunk_size"] = config.ChunkSize
		globalRegistry.logger.InfoWithCategory(Category.API.Messaging, SubCategory.Operation.Initialization, "Starting batch consumption", extras)
	}

	conn, err := getConnection(connName)
	if err != nil {
		if globalRegistry.logger != nil {
			extras := make(map[string]interface{})
			extras["connection_name"] = connName
			extras[ExtraKey.Error.ErrorMessage] = err.Error()
			globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to get connection for batch consuming", extras)
		}
		return fmt.Errorf("failed to get connection: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		if globalRegistry.logger != nil {
			extras := make(map[string]interface{})
			extras["connection_name"] = connName
			extras[ExtraKey.Error.ErrorMessage] = err.Error()
			globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to create channel for batch consuming", extras)
		}
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer ch.Close()

	// Configure QoS for batch processing
	err = ch.Qos(config.ChunkSize, 0, false)
	if err != nil {
		if globalRegistry.logger != nil {
			extras := make(map[string]interface{})
			extras["connection_name"] = connName
			extras["queue_name"] = config.QueueName
			extras["chunk_size"] = config.ChunkSize
			extras[ExtraKey.Error.ErrorMessage] = err.Error()
			globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to set QoS for batch consuming", extras)
		}
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	deliveries, err := ch.Consume(
		config.QueueName,
		config.ConsumerName,
		config.AutoAck,
		config.Exclusive,
		config.NoLocal,
		config.NoWait,
		config.Args,
	)

	if err != nil {
		if globalRegistry.logger != nil {
			extras := make(map[string]interface{})
			extras["connection_name"] = connName
			extras["queue_name"] = config.QueueName
			extras["consumer_name"] = config.ConsumerName
			extras[ExtraKey.Error.ErrorMessage] = err.Error()
			globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to start batch consuming", extras)
		}
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	// Process messages in batches
	batch := make([]*Message, 0, config.ChunkSize)
	deliveryTags := make([]uint64, 0, config.ChunkSize)

	for {
		select {
		case <-ctx.Done():
			if globalRegistry.logger != nil {
				extras := make(map[string]interface{})
				extras["connection_name"] = connName
				extras["queue_name"] = config.QueueName
				globalRegistry.logger.InfoWithCategory(Category.API.Messaging, SubCategory.Operation.Shutdown, "Batch consumer context cancelled", extras)
			}
			return nil
		case d, ok := <-deliveries:
			if !ok {
				if globalRegistry.logger != nil {
					extras := make(map[string]interface{})
					extras["connection_name"] = connName
					extras["queue_name"] = config.QueueName
					globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Delivery channel closed for batch consumer", extras)
				}
				return fmt.Errorf("delivery channel closed")
			}

			// Add message to batch
			msg := FromAMQP(d)
			batch = append(batch, msg)
			deliveryTags = append(deliveryTags, d.DeliveryTag)

			// If we have a full batch or it's the end of the channel, process the batch
			if len(batch) >= config.ChunkSize {
				if globalRegistry.logger != nil {
					extras := make(map[string]interface{})
					extras["connection_name"] = connName
					extras["queue_name"] = config.QueueName
					extras["batch_size"] = len(batch)
					globalRegistry.logger.DebugWithCategory(Category.API.Messaging, SubCategory.Status.Debug, "Processing message batch", extras)
				}

				// Process batch with retry
				var processErr error
				retries := 0

				for retries <= config.MaxRetries {
					if retries > 0 && globalRegistry.logger != nil {
						extras := make(map[string]interface{})
						extras["connection_name"] = connName
						extras["queue_name"] = config.QueueName
						extras["batch_size"] = len(batch)
						extras["retry_count"] = retries
						extras["max_retries"] = config.MaxRetries
						globalRegistry.logger.WarnWithCategory(Category.Error.Retry, SubCategory.Status.Retry, "Retrying batch processing", extras)
					}

					processErr = handler(batch)
					if processErr == nil {
						break
					}

					retries++
					if retries <= config.MaxRetries {
						time.Sleep(config.RetryInterval)
					}
				}

				if !config.AutoAck {
					if processErr == nil {
						// Acknowledge all messages in the batch
						if err := ch.Ack(deliveryTags[len(deliveryTags)-1], true); err != nil {
							if globalRegistry.logger != nil {
								extras := make(map[string]interface{})
								extras["connection_name"] = connName
								extras["queue_name"] = config.QueueName
								extras[ExtraKey.Error.ErrorMessage] = err.Error()
								globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to ack batch", extras)
							}
							return fmt.Errorf("failed to ack batch: %w", err)
						}

						if globalRegistry.logger != nil {
							extras := make(map[string]interface{})
							extras["connection_name"] = connName
							extras["queue_name"] = config.QueueName
							extras["batch_size"] = len(batch)
							globalRegistry.logger.InfoWithCategory(Category.API.Messaging, SubCategory.Status.Success, "Successfully processed batch", extras)
						}
					} else {
						// Max retries reached, nack all messages
						if globalRegistry.logger != nil {
							extras := make(map[string]interface{})
							extras["connection_name"] = connName
							extras["queue_name"] = config.QueueName
							extras["batch_size"] = len(batch)
							extras["retry_count"] = retries
							extras["max_retries"] = config.MaxRetries
							extras[ExtraKey.Error.ErrorMessage] = processErr.Error()
							globalRegistry.logger.ErrorWithCategory(Category.Error.Error, SubCategory.Status.Error, "Failed to process batch after max retries", extras)
						}
						return processErr
					}
				}

				// Reset batch and delivery tags for the next batch
				batch = make([]*Message, 0, config.ChunkSize)
				deliveryTags = make([]uint64, 0, config.ChunkSize)
			}
		}
	}
}

// processBatch processes a batch of messages with retry support
func processBatch(ctx context.Context, ch *amqp.Channel, batch []*Message, handler func([]*Message) error, config *ConsumeConfig) error {
	retries := 0
	var lastErr error

	for retries <= config.MaxRetries {
		if retries > 0 && globalRegistry.logger != nil {
			extras := make(map[string]interface{})
			extras["retry_count"] = retries
			extras["max_retries"] = config.MaxRetries
			if lastErr != nil {
				extras[ExtraKey.Error.ErrorMessage] = lastErr.Error()
			}
			globalRegistry.logger.WarnWithCategory(Category.Error.Retry, SubCategory.Status.Retry, "Retrying batch processing", extras)
		}

		if err := handler(batch); err != nil {
			lastErr = err
			retries++
			if retries <= config.MaxRetries {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(config.RetryInterval):
				}
			}
			continue
		}

		// If successful, ack all messages in the batch
		if !config.AutoAck {
			// Get the last delivery tag
			lastTag := batch[len(batch)-1].DeliveryTag
			if err := ch.Ack(lastTag, true); err != nil {
				if globalRegistry.logger != nil {
					extras := make(map[string]interface{})
					extras["batch_size"] = len(batch)
					extras[ExtraKey.Error.ErrorMessage] = err.Error()
					globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to ack batch", extras)
				}
				return fmt.Errorf("failed to ack batch: %w", err)
			}
		}

		if globalRegistry.logger != nil {
			extras := make(map[string]interface{})
			extras["batch_size"] = len(batch)
			globalRegistry.logger.InfoWithCategory(Category.API.Messaging, SubCategory.Status.Success, "Successfully processed batch", extras)
		}
		return nil
	}

	// Max retries reached
	if globalRegistry.logger != nil {
		extras := make(map[string]interface{})
		extras["batch_size"] = len(batch)
		extras["retry_count"] = retries
		extras["max_retries"] = config.MaxRetries
		if lastErr != nil {
			extras[ExtraKey.Error.ErrorMessage] = lastErr.Error()
		}
		globalRegistry.logger.ErrorWithCategory(Category.Error.Error, SubCategory.Status.Error, "Failed to process batch after max retries", extras)
	}

	return fmt.Errorf("batch processing failed after %d retries: %w", config.MaxRetries, lastErr)
}
