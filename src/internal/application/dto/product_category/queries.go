package product_category

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// AdminGetAllCategoryQuery for admin listing of all article categories with pagination
type AdminGetAllCategoryQuery struct {
	common.PaginationRequestDto
}

// GetAllCategoryQuery for listing article categories by site ID with pagination
type GetAllCategoryQuery struct {
	common.PaginationRequestDto
	SiteID *int64 `json:"siteId" validate:"required"`
}

// GetByIdCategoryQuery for retrieving a single article product_category by ID
type GetByIdCategoryQuery struct {
	ID *int64 `json:"id" validate:"required,gt=0"`
}
