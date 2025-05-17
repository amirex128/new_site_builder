package discountusecase

import (
	"fmt"
	"strconv"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/discount"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type DiscountUsecase struct {
	logger sflogger.Logger
	repo   repository.IDiscountRepository
}

func NewDiscountUsecase(c contract.IContainer) *DiscountUsecase {
	return &DiscountUsecase{
		logger: c.GetLogger(),
		repo:   c.GetDiscountRepo(),
	}
}

func (u *DiscountUsecase) CreateDiscountCommand(params *discount.CreateDiscountCommand) (any, error) {
	// Implementation for creating a discount
	fmt.Println(params)

	newDiscount := domain.Discount{
		Code:       *params.Code,
		Quantity:   *params.Quantity,
		Type:       strconv.Itoa(int(*params.Type)),
		Value:      *params.Value,
		ExpiryDate: *params.ExpiryDate,
		SiteID:     *params.SiteID,
		UserID:     1, // Should come from auth context
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsDeleted:  false,
	}

	err := u.repo.Create(newDiscount)
	if err != nil {
		return nil, err
	}

	return newDiscount, nil
}

func (u *DiscountUsecase) UpdateDiscountCommand(params *discount.UpdateDiscountCommand) (any, error) {
	// Implementation for updating a discount
	fmt.Println(params)

	existingDiscount, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

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
		existingDiscount.ExpiryDate = *params.ExpiryDate
	}

	existingDiscount.UpdatedAt = time.Now()

	err = u.repo.Update(existingDiscount)
	if err != nil {
		return nil, err
	}

	return existingDiscount, nil
}

func (u *DiscountUsecase) DeleteDiscountCommand(params *discount.DeleteDiscountCommand) (any, error) {
	// Implementation for deleting a discount
	fmt.Println(params)

	err := u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": params.ID,
	}, nil
}

func (u *DiscountUsecase) GetByIdDiscountQuery(params *discount.GetByIdDiscountQuery) (any, error) {
	// Implementation to get discount by ID
	fmt.Println(params)

	result, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *DiscountUsecase) GetAllDiscountQuery(params *discount.GetAllDiscountQuery) (any, error) {
	// Implementation to get all discounts by site ID
	fmt.Println(params)

	result, count, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}

func (u *DiscountUsecase) AdminGetAllDiscountQuery(params *discount.AdminGetAllDiscountQuery) (any, error) {
	// Implementation to get all discounts for admin
	fmt.Println(params)

	result, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}
