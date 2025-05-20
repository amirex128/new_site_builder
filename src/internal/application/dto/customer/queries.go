package customer

import (
	"github.com/amirex128/new_site_builder/src/internal/application/dto/user"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// AdminGetAllCustomerQuery represents a query to get all customers with pagination for admin
type AdminGetAllCustomerQuery struct {
	common.PaginationRequestDto
}

// GetProfileCustomerQuery represents a query to get customer profile
type GetProfileCustomerQuery struct {
}

// VerifyCustomerQuery represents a query to verify customer
type VerifyCustomerQuery struct {
	Email *string              `json:"email" form:"email" validate:"required,email"`
	Code  *int                 `json:"code" form:"code" validate:"required"`
	Type  *user.VerifyTypeEnum `json:"type" form:"type" validate:"required,enum"`
}
