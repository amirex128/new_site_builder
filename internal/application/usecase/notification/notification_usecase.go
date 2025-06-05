package notificationusecase

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/internal/application/usecase"
	"github.com/amirex128/new_site_builder/internal/contract"
)

type NotificationUsecase struct {
	*usecase.BaseUsecase
	logger sflogger.Logger
}

func NewNotificationUsecase(c contract.IContainer) *NotificationUsecase {
	return &NotificationUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger:      c.GetLogger(),
			AuthContext: c.GetAuthTransientService(),
		},
	}
}

func (u *NotificationUsecase) SmsUsecase() bool {
	return true
}

func (u *NotificationUsecase) EmailUsecase() bool {
	return true

}
