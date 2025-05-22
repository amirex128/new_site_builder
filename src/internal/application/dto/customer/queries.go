package customer

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
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
	Email *string               `json:"email" nameFa:"ایمیل" form:"email" validate:"required,email"`
	Code  *string               `json:"code" nameFa:"کد" form:"code" validate:"required"`
	Type  *enums.VerifyTypeEnum `json:"type" nameFa:"نوع تایید" form:"type" validate:"required,enum"`
}
