package customer_ticket

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
)

// GetByIdCustomerTicketQuery represents a query to get a customer ticket by ID
type GetByIdCustomerTicketQuery struct {
	ID *int64 `json:"id" nameFa:"شناسه" form:"id" validate:"required,gt=0"`
}

// GetAllCustomerTicketQuery represents a query to get all customer tickets with pagination
type GetAllCustomerTicketQuery struct {
	common.PaginationRequestDto
}

// AdminGetAllCustomerTicketQuery represents a query for admin to get all customer tickets with pagination
type AdminGetAllCustomerTicketQuery struct {
	common.PaginationRequestDto
}
