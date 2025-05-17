package basketusecase

import (
	"fmt"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/basket"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
)

type BasketUsecase struct {
	logger sflogger.Logger
	repo   repository.IBasketRepository
}

func NewBasketUsecase(c contract.IContainer) *BasketUsecase {
	return &BasketUsecase{
		logger: c.GetLogger(),
		repo:   c.GetBasketRepo(),
	}
}

func (u *BasketUsecase) UpdateBasketCommand(params *basket.UpdateBasketCommand) (any, error) {
	// Implementation for updating a basket
	fmt.Println(params)

	// TODO: Implement proper basket update logic
	// This might involve:
	// 1. Finding the active basket for the current user/customer
	// 2. Updating the basket items
	// 3. Saving the updated basket

	return map[string]interface{}{
		"success": true,
	}, nil
}

func (u *BasketUsecase) GetBasketQuery(params *basket.GetBasketQuery) (any, error) {
	// Implementation to get a basket by site ID
	fmt.Println(params)

	// TODO: Implement proper basket retrieval logic
	// This would typically:
	// 1. Find the customer ID from the authenticated user
	// 2. Get the active basket for that customer and site

	// Placeholder response
	return map[string]interface{}{
		"items": []interface{}{},
		"total": 0,
	}, nil
}

func (u *BasketUsecase) GetAllBasketUserQuery(params *basket.GetAllBasketUserQuery) (any, error) {
	// Implementation to get all baskets for a user by site ID
	fmt.Println(params)

	// TODO: Get user ID from authentication context
	userID := int64(1) // Placeholder

	result, count, err := u.repo.GetAllByCustomerID(userID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}

func (u *BasketUsecase) AdminGetAllBasketUserQuery(params *basket.AdminGetAllBasketUserQuery) (any, error) {
	// Implementation to get all baskets for admin
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
