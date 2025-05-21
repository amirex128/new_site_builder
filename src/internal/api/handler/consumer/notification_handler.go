package consumer

import (
	"context"
	sfrabbitmq "git.snappfood.ir/backend/go/packages/sf-rabbitmq"
	"github.com/amirex128/new_site_builder/src/internal/application/usecase/notification"
)

type NotificationConsumerHandler struct {
	usecase *notificationusecase.NotificationUsecase
}

func NewNotificationConsumerHandler(usc *notificationusecase.NotificationUsecase) *NotificationConsumerHandler {
	return &NotificationConsumerHandler{
		usecase: usc,
	}
}

func (h *NotificationConsumerHandler) SmsHandler(ctx context.Context, msg *sfrabbitmq.Message) error {
	h.usecase.SmsUsecase()
	return nil
}
func (h *NotificationConsumerHandler) EmailHandler(ctx context.Context, msg *sfrabbitmq.Message) error {

	h.usecase.EmailUsecase()
	return nil
}
