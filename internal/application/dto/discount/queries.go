package discount

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
)

// AdminGetAllDiscountQuery for admin listing of all discounts with pagination
type AdminGetAllDiscountQuery struct {
	common.PaginationRequestDto
}

// GetAllDiscountQuery for listing discounts by site ID with pagination
type GetAllDiscountQuery struct {
	common.PaginationRequestDto
	SiteID *int64 `json:"siteId" form:"siteId" nameFa:"شناسه سایت" validate:"required"`
}

// GetByIdDiscountQuery for retrieving a single discount by ID
type GetByIdDiscountQuery struct {
	ID *int64 `json:"id" form:"id" nameFa:"شناسه" validate:"required,gt=0"`
}
