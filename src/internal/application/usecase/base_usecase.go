package usecase

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
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

func (u *BaseUsecase) CheckAccessModel[T any](existingAddress T) (*resp.Response, error) {
	if *existingAddress.CustomerID > 0 && *existingAddress.CustomerID != *customerID {
		return nil, resp.NewError(resp.Unauthorized, "شما اجازه ویرایش این آدرس را ندارید")
	}

	if *existingAddress.UserID > 0 && *existingAddress.UserID != *userID {
		return nil, resp.NewError(resp.Unauthorized, "شما اجازه ویرایش این آدرس را ندارید")
	}
}
