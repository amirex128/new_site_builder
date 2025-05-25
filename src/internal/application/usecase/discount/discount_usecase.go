package discountusecase

import (
	"errors"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/gin-gonic/gin"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/discount"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type DiscountUsecase struct {
	*usecase.BaseUsecase
	discountRepo repository.IDiscountRepository
	authContext  func(c *gin.Context) service.IAuthService
}

func NewDiscountUsecase(c contract.IContainer) *DiscountUsecase {
	return &DiscountUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		discountRepo: c.GetDiscountRepo(),
		authContext:  c.GetAuthTransientService(),
	}
}

func (u *DiscountUsecase) CreateDiscountCommand(params *discount.CreateDiscountCommand) (*resp.Response, error) {
	existingDiscount, err := u.discountRepo.GetByCode(*params.Code)
	if err != nil {

	}
	if existingDiscount.SiteID == *params.SiteID {
		return nil, resp.NewError(resp.BadRequest, "کد تخفیف تکراری است")
	}

	if params.ExpiryDate.Before(time.Now()) {
		return nil, resp.NewError(resp.BadRequest, "تاریخ انقضا باید در آینده باشد")
	}

	userID, _, _, err := u.authContext(u.Ctx).GetUserOrCustomerID()
	if err != nil || userID == nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}

	newDiscount := domain.Discount{
		Code:            *params.Code,
		Quantity:        *params.Quantity,
		Type:            *params.Type,
		Value:           *params.Value,
		ExpiryDate:      params.ExpiryDate,
		MaxUsagePerUser: params.MaxUsagePerUser,
		SiteID:          *params.SiteID,
		UserID:          *userID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		IsDeleted:       false,
	}

	err = u.discountRepo.Create(&newDiscount)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	createdDiscount, err := u.discountRepo.GetByID(newDiscount.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	return resp.NewResponseData(resp.Created, map[string]interface{}{"discount": createdDiscount}, "تخفیف با موفقیت ایجاد شد"), nil
}

func (u *DiscountUsecase) UpdateDiscountCommand(params *discount.UpdateDiscountCommand) (*resp.Response, error) {
	u.Logger.Info("UpdateDiscountCommand called", map[string]interface{}{
		"id": params.ID,
	})

	existingDiscount, err := u.discountRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "تخفیف یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	userID, _, _, err := u.authContext(u.Ctx).GetUserOrCustomerID()
	if err != nil || userID == nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}
	if existingDiscount.UserID != *userID && !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "شما به این تخفیف دسترسی ندارید")
	}

	if params.Code != nil && *params.Code != existingDiscount.Code {
		codeDiscount, err := u.discountRepo.GetByCode(*params.Code)
		if err == nil && codeDiscount.ID != *params.ID {
			return nil, resp.NewError(resp.BadRequest, "کد تخفیف تکراری است")
		}
	}

	if params.Code != nil {
		existingDiscount.Code = *params.Code
	}
	if params.Quantity != nil {
		existingDiscount.Quantity = *params.Quantity
	}
	if params.Type != nil {
		existingDiscount.Type = *params.Type
	}
	if params.Value != nil {
		existingDiscount.Value = *params.Value
	}
	if params.ExpiryDate != nil {
		if params.ExpiryDate.Before(time.Now()) {
			return nil, resp.NewError(resp.BadRequest, "تاریخ انقضا باید در آینده باشد")
		}
		existingDiscount.ExpiryDate = params.ExpiryDate
	}

	existingDiscount.UpdatedAt = time.Now()

	err = u.discountRepo.Update(existingDiscount)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	updatedDiscount, err := u.discountRepo.GetByID(existingDiscount.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	return resp.NewResponseData(resp.Updated, map[string]interface{}{"discount": updatedDiscount}, "تخفیف با موفقیت بروزرسانی شد"), nil
}

func (u *DiscountUsecase) DeleteDiscountCommand(params *discount.DeleteDiscountCommand) (*resp.Response, error) {
	u.Logger.Info("DeleteDiscountCommand called", map[string]interface{}{
		"id": params.ID,
	})

	existingDiscount, err := u.discountRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "تخفیف یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	userID, _, _, err := u.authContext(u.Ctx).GetUserOrCustomerID()
	if err != nil || userID == nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}
	if existingDiscount.UserID != *userID && !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "شما به این تخفیف دسترسی ندارید")
	}

	err = u.discountRepo.Delete(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	return resp.NewResponseData(resp.Deleted, map[string]interface{}{"success": true}, "تخفیف با موفقیت حذف شد"), nil
}

func (u *DiscountUsecase) GetByIdDiscountQuery(params *discount.GetByIdDiscountQuery) (*resp.Response, error) {
	u.Logger.Info("GetByIdDiscountQuery called", map[string]interface{}{
		"id": params.ID,
	})

	discountObj, err := u.discountRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "تخفیف یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	userID, _, _, _ := u.authContext(u.Ctx).GetUserOrCustomerID()
	if userID != nil {
		u.Logger.Info("Discount accessed by user", map[string]interface{}{
			"discountId": discountObj.ID,
			"userId":     *userID,
		})
	}

	response := map[string]interface{}{
		"id":         discountObj.ID,
		"code":       discountObj.Code,
		"quantity":   discountObj.Quantity,
		"type":       discountObj.Type,
		"value":      discountObj.Value,
		"expiryDate": discountObj.ExpiryDate,
		"siteId":     discountObj.SiteID,
		"createdAt":  discountObj.CreatedAt,
		"updatedAt":  discountObj.UpdatedAt,
	}

	return resp.NewResponseData(resp.Retrieved, response, "تخفیف با موفقیت دریافت شد"), nil
}

func (u *DiscountUsecase) GetAllDiscountQuery(params *discount.GetAllDiscountQuery) (*resp.Response, error) {
	u.Logger.Info("GetAllDiscountQuery called", map[string]interface{}{
		"siteId":   params.SiteID,
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	results, err := u.discountRepo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	items := make([]map[string]interface{}, 0, len(results.Items))
	for _, discountObj := range results.Items {
		item := map[string]interface{}{
			"id":         discountObj.ID,
			"code":       discountObj.Code,
			"quantity":   discountObj.Quantity,
			"type":       discountObj.Type,
			"value":      discountObj.Value,
			"expiryDate": discountObj.ExpiryDate,
			"siteId":     discountObj.SiteID,
			"createdAt":  discountObj.CreatedAt,
			"updatedAt":  discountObj.UpdatedAt,
		}
		items = append(items, item)
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     items,
		"total":     results.TotalCount,
		"page":      results.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": results.TotalPages,
	}, "لیست تخفیف ها با موفقیت دریافت شد"), nil
}

func (u *DiscountUsecase) AdminGetAllDiscountQuery(params *discount.AdminGetAllDiscountQuery) (*resp.Response, error) {
	u.Logger.Info("AdminGetAllDiscountQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}

	results, err := u.discountRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	items := make([]map[string]interface{}, 0, len(results.Items))
	for _, discountObj := range results.Items {
		item := map[string]interface{}{
			"id":         discountObj.ID,
			"code":       discountObj.Code,
			"quantity":   discountObj.Quantity,
			"type":       discountObj.Type,
			"value":      discountObj.Value,
			"expiryDate": discountObj.ExpiryDate,
			"siteId":     discountObj.SiteID,
			"createdAt":  discountObj.CreatedAt,
			"updatedAt":  discountObj.UpdatedAt,
		}
		items = append(items, item)
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     items,
		"total":     results.TotalCount,
		"page":      results.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": results.TotalPages,
	}, "لیست تخفیف ها با موفقیت دریافت شد (ادمین)"), nil
}
