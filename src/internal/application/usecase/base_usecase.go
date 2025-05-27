package usecase

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"
	"github.com/gin-gonic/gin"
)

type BaseUsecase struct {
	Ctx         *gin.Context
	Logger      sflogger.Logger
	AuthContext func(c *gin.Context) service.IAuthService
}

func (u *BaseUsecase) SetContext(c *gin.Context) *BaseUsecase {
	u.Ctx = c
	return u
}

func (u *BaseUsecase) CheckAccessUserModel(existingModel common.AccessControllable) error {
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil {
		return err
	}

	if existingModel.GetUserID() != nil && *existingModel.GetUserID() > 0 && existingModel.GetUserID() != nil && *existingModel.GetUserID() != *userID {
		return resp.NewError(resp.Unauthorized, "شما اجازه ویرایش را ندارید")
	}
	return nil
}
func (u *BaseUsecase) CheckAccessCustomerModel(existingModel common.AccessControllable) error {
	customerID, err := u.AuthContext(u.Ctx).GetCustomerID()
	if err != nil {
		return err
	}
	if existingModel.GetCustomerID() != nil && *existingModel.GetCustomerID() > 0 && existingModel.GetCustomerID() != nil && *existingModel.GetCustomerID() != *customerID {
		return resp.NewError(resp.Unauthorized, "شما اجازه ویرایش را ندارید")
	}
	return nil
}

func (u *BaseUsecase) CheckAccessSiteModel(siteID *int64) error {
	siteIDs, err := u.AuthContext(u.Ctx).GetSiteIDs()
	if err != nil {
		return resp.NewError(resp.Unauthorized, "خطا در بررسی دسترسی کاربر")
	}
	if siteID == nil || *siteID <= 0 {
		return resp.NewError(resp.BadRequest, "شناسه سایت معتبر نیست")
	}
	for _, id := range siteIDs {
		if id == *siteID {
			return nil
		}
	}
	return resp.NewError(resp.Unauthorized, "شما به این سایت دسترسی ندارید")
}
func (u *BaseUsecase) CheckAccessAdmin() error {
	isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return resp.NewError(resp.Unauthorized, "فقط مدیران به این بخش دسترسی دارند")
	}
	return nil
}
