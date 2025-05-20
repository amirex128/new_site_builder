package payment

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// AdminGetAllGatewayQuery for admin to get all payment gateways with pagination
type AdminGetAllGatewayQuery struct {
	common.PaginationRequestDto
}

// GetByIdGatewayQuery for getting gateway details by ID
type GetByIdGatewayQuery struct {
	ID *int64 `json:"id" form:"id" validate:"required,gt=0"`
}

// AdminGetAllPaymentQuery for admin to get all payments with pagination
type AdminGetAllPaymentQuery struct {
	common.PaginationRequestDto
}
