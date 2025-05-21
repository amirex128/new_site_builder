package service

import (
	"context"
	"encoding/json"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sfrabbitmq "git.snappfood.ir/backend/go/packages/sf-rabbitmq"
)

type RabbitMqService struct {
	ctx    *context.Context
	logger sflogger.Logger
}

func NewRabbitMqService(ctx *context.Context, logger sflogger.Logger) *RabbitMqService {
	return &RabbitMqService{
		ctx:    ctx,
		logger: logger,
	}
}

func (s *RabbitMqService) SendSms(msg any) error {
	body, err := json.Marshal(msg)
	if err != nil {
		s.logger.Errorf("Failed to marshal message: %v", err)
		return err
	}

	err = sfrabbitmq.Publish(s.ctx, "order-connection", "order_exchange", "order_routing_key",
		&sfrabbitmq.Message{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		s.logger.Errorf("Failed to publish message: %v", err)
		return err
	}
	return nil

}

func (s *RabbitMqService) SendEmail(msg any) error {
	body, err := json.Marshal(msg)
	if err != nil {
		s.logger.Errorf("Failed to marshal message: %v", err)
		return err
	}

	err = sfrabbitmq.Publish(s.ctx, "order-connection", "order_exchange", "order_routing_key",
		&sfrabbitmq.Message{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		s.logger.Errorf("Failed to publish message: %v", err)
		return err
	}
	return nil
}
