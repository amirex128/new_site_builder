package consumer

import (
	"context"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sfrabbitmq "git.snappfood.ir/backend/go/packages/sf-rabbitmq"
	"github.com/amirex128/new_site_builder/src/internal/application/usecase/notification"
)

type NotificationConsumerHandler struct {
	usecase *notificationusecase.NotificationUsecase
	logger  sflogger.Logger
}

func NewNotificationConsumerHandler(usc *notificationusecase.NotificationUsecase, logger sflogger.Logger) *NotificationConsumerHandler {
	return &NotificationConsumerHandler{
		usecase: usc,
		logger:  logger,
	}
}

func (h *NotificationConsumerHandler) SmsHandler(ctx context.Context, msg *sfrabbitmq.Message) error {
	h.logger.Infof("Start SmsHandler consumers successfully")

	h.usecase.SmsUsecase()
	return nil
}
func (h *NotificationConsumerHandler) EmailHandler(ctx context.Context, msg *sfrabbitmq.Message) error {
	h.logger.Infof("Start EmailHandler consumers successfully")

	h.usecase.EmailUsecase()
	return nil
}
