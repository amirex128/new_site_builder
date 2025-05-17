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
	Email *string              `json:"email" validate:"required,email" error:"required=ایمیل الزامی است|email=ایمیل باید معتبر باشد"`
	Code  *int                 `json:"code" validate:"required" error:"required=کد تایید الزامی است"`
	Type  *user.VerifyTypeEnum `json:"type" validate:"required" error:"required=نوع تایید الزامی است"`
}
