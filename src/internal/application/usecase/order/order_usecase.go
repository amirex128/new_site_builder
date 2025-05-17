package orderusecase

import (
	"fmt"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/order"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
)

type OrderUsecase struct {
	logger sflogger.Logger
	repo   repository.IOrderRepository
}

func NewOrderUsecase(c contract.IContainer) *OrderUsecase {
	return &OrderUsecase{
		logger: c.GetLogger(),
		repo:   c.GetOrderRepo(),
	}
}

func (u *OrderUsecase) CreateOrderRequestCommand(params *order.CreateOrderRequestCommand) (any, error) {
	// Implementation for creating an order request
	fmt.Println(params)

	// TODO: Implement proper order creation logic
	// This might involve:
	// 1. Getting the current user's basket
	// 2. Converting the basket to an order
	// 3. Setting up payment via the specified gateway
	// 4. Redirecting to payment URL

	return map[string]interface{}{
		"success":    true,
		"paymentUrl": "https://example.com/payment",
	}, nil
}

func (u *OrderUsecase) GetAllOrderCustomerQuery(params *order.GetAllOrderCustomerQuery) (any, error) {
	// Implementation to get all orders for a customer
	fmt.Println(params)

	// TODO: Get customer ID from authentication context
	customerID := int64(1) // Placeholder

	result, count, err := u.repo.GetAllByCustomerID(customerID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}

func (u *OrderUsecase) GetOrderCustomerDetailsQuery(params *order.GetOrderCustomerDetailsQuery) (any, error) {
	// Implementation to get order details for a customer
	fmt.Println(params)

	result, err := u.repo.GetByID(*params.OrderID)
	if err != nil {
		return nil, err
	}

	// TODO: Verify the order belongs to the current customer

	return result, nil
}

func (u *OrderUsecase) GetAllOrderUserQuery(params *order.GetAllOrderUserQuery) (any, error) {
	// Implementation to get all orders for a user by site ID
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

func (u *OrderUsecase) GetOrderUserDetailsQuery(params *order.GetOrderUserDetailsQuery) (any, error) {
	// Implementation to get order details for a user
	fmt.Println(params)

	result, err := u.repo.GetByID(*params.OrderID)
	if err != nil {
		return nil, err
	}

	// TODO: Verify the order belongs to the user's site

	return result, nil
}

func (u *OrderUsecase) AdminGetAllOrderUserQuery(params *order.AdminGetAllOrderUserQuery) (any, error) {
	// Implementation to get all orders for admin
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
