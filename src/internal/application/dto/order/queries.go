package order

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// GetAllOrderUserQuery for getting all orders for a user with pagination
type GetAllOrderUserQuery struct {
	common.PaginationRequestDto
	SiteID *int64 `json:"siteId" validate:"required"`
}

// GetAllOrderCustomerQuery for getting all orders for a customer with pagination
type GetAllOrderCustomerQuery struct {
	common.PaginationRequestDto
	// No additional fields needed
}

// AdminGetAllOrderUserQuery for admin to get all orders with pagination
type AdminGetAllOrderUserQuery struct {
	common.PaginationRequestDto
	// No additional fields needed
}

// GetOrderUserDetailsQuery for getting order details for a user
type GetOrderUserDetailsQuery struct {
	OrderID *int64 `json:"orderId" validate:"required"`
}

// GetOrderCustomerDetailsQuery for getting order details for a customer
type GetOrderCustomerDetailsQuery struct {
	OrderID *int64 `json:"orderId" validate:"required"`
}
