package unitpriceusecase

import (
	"fmt"
	"strconv"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/unit_price"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
)

type UnitPriceUsecase struct {
	logger        sflogger.Logger
	unitPriceRepo repository.IUnitPriceRepository
	userRepo      repository.IUserRepository
	authContext   common.IAuthContextService
}

func NewUnitPriceUsecase(c contract.IContainer) *UnitPriceUsecase {
	return &UnitPriceUsecase{
		logger:        c.GetLogger(),
		unitPriceRepo: c.GetUnitPriceRepo(),
		userRepo:      c.GetUserRepo(),
		authContext:   c.GetAuthContextService(),
	}
}

// UpdateUnitPriceCommand updates a unit price
// Based on UpdateUnitPriceCommand.cs
func (u *UnitPriceUsecase) UpdateUnitPriceCommand(params *unit_price.UpdateUnitPriceCommand) (any, error) {
	// Check admin access
	isAdmin, err := u.authContext.IsAdmin()
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, fmt.Errorf("only admins can update unit prices")
	}

	// Get the existing unit price
	existingUnitPrice, err := u.unitPriceRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if params.Name != nil {
		existingUnitPrice.Name = *params.Name
	}

	if params.Price != nil {
		existingUnitPrice.Price = int64(*params.Price)
	}

	if params.HasDay != nil {
		existingUnitPrice.HasDay = *params.HasDay
	}

	if params.DiscountType != nil {
		existingUnitPrice.DiscountType = strconv.Itoa(int(*params.DiscountType))
	}

	if params.Discount != nil {
		discount := int64(*params.Discount)
		existingUnitPrice.Discount = &discount
	}

	// Update the unit price
	err = u.unitPriceRepo.Update(existingUnitPrice)
	if err != nil {
		return nil, err
	}

	return existingUnitPrice, nil
}

// CalculateUnitPriceQuery calculates the price for unit prices
// Based on CalculateUnitPriceQuery.cs
func (u *UnitPriceUsecase) CalculateUnitPriceQuery(params *unit_price.CalculateUnitPriceQuery) (any, error) {
	// Get the current user ID
	userID, err := u.authContext.GetUserID()
	if err != nil {
		return nil, err
	}
	if userID == 0 {
		return nil, fmt.Errorf("user not authenticated")
	}

	// Get all unit prices by name
	var unitPriceNames []string
	for _, up := range params.UnitPrices {
		unitPriceNames = append(unitPriceNames, strconv.Itoa(int(*up.UnitPriceName)))
	}

	// In a real implementation, we would need to add a method to get unit prices by names
	// For now, we'll get all unit prices and filter them
	allUnitPrices, _, err := u.unitPriceRepo.GetAll(common.PaginationRequestDto{Page: 1, PageSize: 100})
	if err != nil {
		return nil, err
	}

	// Get the current user
	currentUser, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Use the current user for storage credit calculations if needed
	_ = currentUser

	var data []map[string]interface{}
	for _, unitPriceParam := range params.UnitPrices {
		// Find the matching unit price
		var matchingUnitPrice *unit_price.UnitPriceNameEnum
		for _, up := range allUnitPrices {
			if up.Name == strconv.Itoa(int(*unitPriceParam.UnitPriceName)) {
				matchingUnitPrice = unitPriceParam.UnitPriceName

				// Calculate price based on unit price type
				if *unitPriceParam.UnitPriceName == unit_price.StorageMbCredits && unitPriceParam.UnitPriceDay != nil {
					// For storage, we need to consider days and existing credits
					var finalPrice int64 = up.Price * int64(*unitPriceParam.UnitPriceCount) * int64(*unitPriceParam.UnitPriceDay)
					var discountAmount int64 = 0

					// Apply discount if available
					if up.Discount != nil && *up.Discount > 0 {
						if up.DiscountType == strconv.Itoa(int(unit_price.Fixed)) {
							discountAmount = *up.Discount
							finalPrice = finalPrice - discountAmount
						} else if up.DiscountType == strconv.Itoa(int(unit_price.Percentage)) {
							discountAmount = (finalPrice * (*up.Discount)) / 100
							finalPrice = finalPrice - discountAmount
						}
					}

					// Ensure price doesn't go below zero
					if finalPrice < 0 {
						finalPrice = 0
					}

					data = append(data, map[string]interface{}{
						"UnitPriceName":  *unitPriceParam.UnitPriceName,
						"UnitPriceCount": *unitPriceParam.UnitPriceCount,
						"TotalPrice":     finalPrice + discountAmount,
						"FinalPrice":     finalPrice,
						"DiscountAmount": discountAmount,
					})
				} else {
					// For other unit prices
					var finalPrice int64 = up.Price * int64(*unitPriceParam.UnitPriceCount)
					var discountAmount int64 = 0

					// Apply discount if available
					if up.Discount != nil && *up.Discount > 0 {
						if up.DiscountType == strconv.Itoa(int(unit_price.Fixed)) {
							discountAmount = *up.Discount
							finalPrice = finalPrice - discountAmount
						} else if up.DiscountType == strconv.Itoa(int(unit_price.Percentage)) {
							discountAmount = (finalPrice * (*up.Discount)) / 100
							finalPrice = finalPrice - discountAmount
						}
					}

					// Ensure price doesn't go below zero
					if finalPrice < 0 {
						finalPrice = 0
					}

					data = append(data, map[string]interface{}{
						"UnitPriceName":  *unitPriceParam.UnitPriceName,
						"UnitPriceCount": *unitPriceParam.UnitPriceCount,
						"TotalPrice":     finalPrice + discountAmount,
						"FinalPrice":     finalPrice,
						"DiscountAmount": discountAmount,
					})
				}

				break
			}
		}

		if matchingUnitPrice == nil {
			return nil, fmt.Errorf("unit price not found for name: %d", *unitPriceParam.UnitPriceName)
		}
	}

	return data, nil
}

// GetAllUnitPriceQuery gets all unit prices with pagination
// Based on GetAllUnitPriceQuery.cs
func (u *UnitPriceUsecase) GetAllUnitPriceQuery(params *unit_price.GetAllUnitPriceQuery) (any, error) {
	// Check admin access
	isAdmin, err := u.authContext.IsAdmin()
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, fmt.Errorf("only admins can get all unit prices")
	}

	// Get all unit prices
	unitPrices, count, err := u.unitPriceRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": unitPrices,
		"total": count,
	}, nil
}
