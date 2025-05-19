package discountusecase

import (
	"errors"
	"strconv"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/discount"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type DiscountUsecase struct {
	logger         sflogger.Logger
	repo           repository.IDiscountRepository
	authContextSvc common.IAuthContextService
}

func NewDiscountUsecase(c contract.IContainer) *DiscountUsecase {
	return &DiscountUsecase{
		logger:         c.GetLogger(),
		repo:           c.GetDiscountRepo(),
		authContextSvc: c.GetAuthContextTransientService(),
	}
}

func (u *DiscountUsecase) CreateDiscountCommand(params *discount.CreateDiscountCommand) (any, error) {
	u.logger.Info("CreateDiscountCommand called", map[string]interface{}{
		"code":   *params.Code,
		"siteId": *params.SiteID,
	})

	// Check for existing discount code in the same site
	existingDiscount, err := u.repo.GetByCode(*params.Code)
	if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		if err == nil && existingDiscount.SiteID == *params.SiteID {
			return nil, errors.New("کد تخفیف تکراری است")
		}
	}

	// Validate expiry date is in the future
	if params.ExpiryDate.Before(time.Now()) {
		return nil, errors.New("تاریخ انقضا باید در آینده باشد")
	}

	// Get user ID from auth context
	userID, err := u.authContextSvc.GetUserID()
	if err != nil {
		return nil, err
	}

	// Create discount entity
	newDiscount := domain.Discount{
		Code:       *params.Code,
		Quantity:   *params.Quantity,
		Type:       strconv.Itoa(int(*params.Type)),
		Value:      *params.Value,
		ExpiryDate: *params.ExpiryDate,
		SiteID:     *params.SiteID,
		UserID:     userID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsDeleted:  false,
	}

	// Create the discount in the database
	err = u.repo.Create(newDiscount)
	if err != nil {
		return nil, err
	}

	// Fetch the created discount
	createdDiscount, err := u.repo.GetByID(newDiscount.ID)
	if err != nil {
		return nil, err
	}

	return createdDiscount, nil
}

func (u *DiscountUsecase) UpdateDiscountCommand(params *discount.UpdateDiscountCommand) (any, error) {
	u.logger.Info("UpdateDiscountCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing discount
	existingDiscount, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تخفیف یافت نشد")
		}
		return nil, err
	}

	// Check user access
	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}

	userID, err := u.authContextSvc.GetUserID()
	if err != nil {
		return nil, err
	}

	if existingDiscount.UserID != userID && !isAdmin {
		return nil, errors.New("شما به این تخفیف دسترسی ندارید")
	}

	// Validate code uniqueness if changed
	if params.Code != nil && *params.Code != existingDiscount.Code {
		codeDiscount, err := u.repo.GetByCode(*params.Code)
		if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
			if err == nil && codeDiscount.ID != *params.ID {
				return nil, errors.New("کد تخفیف تکراری است")
			}
		}
	}

	// Update fields if provided
	if params.Code != nil {
		existingDiscount.Code = *params.Code
	}

	if params.Quantity != nil {
		existingDiscount.Quantity = *params.Quantity
	}

	if params.Type != nil {
		existingDiscount.Type = strconv.Itoa(int(*params.Type))
	}

	if params.Value != nil {
		existingDiscount.Value = *params.Value
	}

	if params.ExpiryDate != nil {
		// Validate expiry date is in the future
		if params.ExpiryDate.Before(time.Now()) {
			return nil, errors.New("تاریخ انقضا باید در آینده باشد")
		}
		existingDiscount.ExpiryDate = *params.ExpiryDate
	}

	existingDiscount.UpdatedAt = time.Now()

	// Update the discount
	err = u.repo.Update(existingDiscount)
	if err != nil {
		return nil, err
	}

	// Fetch the updated discount
	updatedDiscount, err := u.repo.GetByID(existingDiscount.ID)
	if err != nil {
		return nil, err
	}

	return updatedDiscount, nil
}

func (u *DiscountUsecase) DeleteDiscountCommand(params *discount.DeleteDiscountCommand) (any, error) {
	u.logger.Info("DeleteDiscountCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing discount
	existingDiscount, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تخفیف یافت نشد")
		}
		return nil, err
	}

	// Check user access
	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}

	userID, err := u.authContextSvc.GetUserID()
	if err != nil {
		return nil, err
	}

	if existingDiscount.UserID != userID && !isAdmin {
		return nil, errors.New("شما به این تخفیف دسترسی ندارید")
	}

	// Delete the discount
	err = u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
	}, nil
}

func (u *DiscountUsecase) GetByIdDiscountQuery(params *discount.GetByIdDiscountQuery) (any, error) {
	u.logger.Info("GetByIdDiscountQuery called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get discount by ID
	discount, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تخفیف یافت نشد")
		}
		return nil, err
	}

	// Check user access - anyone can view discounts but logging for audit
	userID, _ := u.authContextSvc.GetUserID()
	if userID > 0 {
		u.logger.Info("Discount accessed by user", map[string]interface{}{
			"discountId": discount.ID,
			"userId":     userID,
		})
	}

	// Parse type string to enum
	typeEnum, _ := strconv.Atoi(discount.Type)

	// Prepare response
	response := map[string]interface{}{
		"id":         discount.ID,
		"code":       discount.Code,
		"quantity":   discount.Quantity,
		"type":       typeEnum,
		"value":      discount.Value,
		"expiryDate": discount.ExpiryDate,
		"siteId":     discount.SiteID,
		"createdAt":  discount.CreatedAt,
		"updatedAt":  discount.UpdatedAt,
	}

	return response, nil
}

func (u *DiscountUsecase) GetAllDiscountQuery(params *discount.GetAllDiscountQuery) (any, error) {
	u.logger.Info("GetAllDiscountQuery called", map[string]interface{}{
		"siteId":   *params.SiteID,
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check site access - simplified in monolithic app
	// In a real implementation, we would check if the user has access to this site

	// Get all discounts for the site
	discounts, count, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// Prepare response items
	var items []map[string]interface{}
	for _, discount := range discounts {
		// Parse type string to enum
		typeEnum, _ := strconv.Atoi(discount.Type)

		item := map[string]interface{}{
			"id":         discount.ID,
			"code":       discount.Code,
			"quantity":   discount.Quantity,
			"type":       typeEnum,
			"value":      discount.Value,
			"expiryDate": discount.ExpiryDate,
			"siteId":     discount.SiteID,
			"createdAt":  discount.CreatedAt,
			"updatedAt":  discount.UpdatedAt,
		}

		items = append(items, item)
	}

	// Return paginated result
	return map[string]interface{}{
		"items":     items,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

func (u *DiscountUsecase) AdminGetAllDiscountQuery(params *discount.AdminGetAllDiscountQuery) (any, error) {
	u.logger.Info("AdminGetAllDiscountQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check admin access
	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		return nil, errors.New("فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}

	// Get all discounts across all sites for admin
	discounts, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// Prepare response items
	var items []map[string]interface{}
	for _, discount := range discounts {
		// Parse type string to enum
		typeEnum, _ := strconv.Atoi(discount.Type)

		item := map[string]interface{}{
			"id":         discount.ID,
			"code":       discount.Code,
			"quantity":   discount.Quantity,
			"type":       typeEnum,
			"value":      discount.Value,
			"expiryDate": discount.ExpiryDate,
			"siteId":     discount.SiteID,
			"createdAt":  discount.CreatedAt,
			"updatedAt":  discount.UpdatedAt,
		}

		items = append(items, item)
	}

	// Return paginated result
	return map[string]interface{}{
		"items":     items,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}
