package article_category

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// GetByIdCategoryQuery for retrieving a single product_category by ID
type GetByIdCategoryQuery struct {
	ID *int64 `json:"id" validate:"required"`
}

// AdminGetAllCategoryQuery for admin listing of all categories with pagination
type AdminGetAllCategoryQuery struct {
	common.PaginationRequestDto
}

// GetAllCategoryQuery for listing categories by site ID with pagination
type GetAllCategoryQuery struct {
	common.PaginationRequestDto
	SiteID *int64 `json:"siteId" validate:"required"`
}
