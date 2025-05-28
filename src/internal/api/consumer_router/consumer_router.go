package consumerrouter

import (
	"context"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sfrabbitmq "git.snappfood.ir/backend/go/packages/sf-rabbitmq"
	"github.com/amirex128/new_site_builder/src/bootstrap"
)

func BindConsumers(ctx *context.Context, handlers *bootstrap.ConsumerHandlerManager, logger sflogger.Logger) {
	logger.Infof("Binding consumers successfully")

	err := sfrabbitmq.Consume(
		*ctx,
		"main",
		"sms",
		handlers.NotificationHandler.SmsHandler,
		&sfrabbitmq.ConsumeConfig{
			QueueName:     "",
			ConsumerName:  "",
			AutoAck:       false,
			Exclusive:     false,
			NoLocal:       false,
			NoWait:        false,
			Args:          nil,
			MaxRetries:    0,
			RetryInterval: 0,
		},
	)
	if err != nil {
		logger.Errorf("Failed to start single message consumption: %v", err)
	}

}
