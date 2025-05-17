package discount

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// AdminGetAllDiscountQuery for admin listing of all discounts with pagination
type AdminGetAllDiscountQuery struct {
	common.PaginationRequestDto
}

// GetAllDiscountQuery for listing discounts by site ID with pagination
type GetAllDiscountQuery struct {
	common.PaginationRequestDto
	SiteID *int64 `json:"siteId" validate:"required" error:"required=سایت الزامی است"`
}

// GetByIdDiscountQuery for retrieving a single discount by ID
type GetByIdDiscountQuery struct {
	ID *int64 `json:"id" validate:"required,gt=0" error:"required=شناسه الزامی است|gt=شناسه باید بزرگتر از 0 باشد"`
}
