package ticket

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// GetByIdTicketQuery represents a query to get a ticket by ID
type GetByIdTicketQuery struct {
	ID *int64 `json:"id" form:"id" validate:"required,gt=0" nameFa:"شناسه"`
}

// GetAllTicketQuery represents a query to get all tickets with pagination
type GetAllTicketQuery struct {
	common.PaginationRequestDto
}

// AdminGetAllTicketQuery represents a query for admin to get all tickets with pagination
type AdminGetAllTicketQuery struct {
	common.PaginationRequestDto
}
