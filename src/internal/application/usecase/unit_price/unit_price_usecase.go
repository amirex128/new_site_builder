package unitpriceusecase

import (
	"fmt"
	"strconv"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/unit_price"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
)

type UnitPriceUsecase struct {
	logger sflogger.Logger
	repo   repository.IUnitPriceRepository
}

func NewUnitPriceUsecase(c contract.IContainer) *UnitPriceUsecase {
	return &UnitPriceUsecase{
		logger: c.GetLogger(),
		repo:   c.GetUnitPriceRepo(),
	}
}

func (u *UnitPriceUsecase) UpdateUnitPriceCommand(params *unit_price.UpdateUnitPriceCommand) (any, error) {
	// Implementation for updating a unit price
	fmt.Println(params)

	existingUnitPrice, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

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

	err = u.repo.Update(existingUnitPrice)
	if err != nil {
		return nil, err
	}

	return existingUnitPrice, nil
}

func (u *UnitPriceUsecase) CalculateUnitPriceQuery(params *unit_price.CalculateUnitPriceQuery) (any, error) {
	// Implementation to calculate unit prices
	fmt.Println(params)

	var totalPrice int64 = 0
	var items []map[string]interface{}

	for _, unitPrice := range params.UnitPrices {
		// Get the unit price from the database
		up, err := u.repo.GetByName(strconv.Itoa(int(*unitPrice.UnitPriceName)))
		if err != nil {
			return nil, err
		}

		// Calculate the base price
		itemPrice := up.Price * int64(*unitPrice.UnitPriceCount)

		// For storage, we need to multiply by days
		if up.HasDay && unitPrice.UnitPriceDay != nil {
			itemPrice = itemPrice * int64(*unitPrice.UnitPriceDay)
		}

		// Apply discount if available
		var discountAmount int64 = 0
		if up.Discount != nil && *up.Discount > 0 {
			if up.DiscountType == strconv.Itoa(int(unit_price.Fixed)) {
				discountAmount = *up.Discount
			} else if up.DiscountType == strconv.Itoa(int(unit_price.Percentage)) {
				discountAmount = (itemPrice * (*up.Discount)) / 100
			}
			itemPrice = itemPrice - discountAmount
		}

		// Ensure price doesn't go below zero
		if itemPrice < 0 {
			itemPrice = 0
		}

		totalPrice += itemPrice

		items = append(items, map[string]interface{}{
			"unitPrice":      up,
			"count":          *unitPrice.UnitPriceCount,
			"days":           unitPrice.UnitPriceDay,
			"originalPrice":  up.Price * int64(*unitPrice.UnitPriceCount),
			"discountAmount": discountAmount,
			"finalPrice":     itemPrice,
		})
	}

	return map[string]interface{}{
		"items":      items,
		"totalPrice": totalPrice,
	}, nil
}

func (u *UnitPriceUsecase) GetAllUnitPriceQuery(params *unit_price.GetAllUnitPriceQuery) (any, error) {
	// Implementation to get all unit prices
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
