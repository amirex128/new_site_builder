package usecase

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/gin-gonic/gin"
)

type BaseUsecase struct {
	Ctx    *gin.Context
	Logger sflogger.Logger
}

func (u *BaseUsecase) SetContext(c *gin.Context) *BaseUsecase {
	u.Ctx = c
	return u
}

func (u *BaseUsecase) CheckAccessUserModel(existingModel common.AccessControllable, userID *int64) error {
	if existingModel.GetUserID() != nil && *existingModel.GetUserID() > 0 && existingModel.GetUserID() != nil && *existingModel.GetUserID() != *userID {
		return resp.NewError(resp.Unauthorized, "شما اجازه ویرایش را ندارید")
	}
	return nil
}
func (u *BaseUsecase) CheckAccessCustomerModel(existingModel common.AccessControllable, customerID *int64) error {
	if existingModel.GetCustomerID() != nil && *existingModel.GetCustomerID() > 0 && existingModel.GetCustomerID != nil && *existingModel.GetCustomerID() != *customerID {
		return resp.NewError(resp.Unauthorized, "شما اجازه ویرایش را ندارید")
	}
	return nil
}
