package page

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// GetByIdPageQuery represents a query to get page(s) by ID
type GetByIdPageQuery struct {
	ID     *int64  `json:"id,omitempty" form:"id" validate:"omitempty,gt=0"`
	IDs    []int64 `json:"ids,omitempty" form:"ids" validate:"array_number_optional=0,100,1,0,false"`
	SiteID *int64  `json:"siteId" form:"siteId" validate:"required,gt=0"`
}

// GetByPathPageQuery represents a query to get page(s) by path
type GetByPathPageQuery struct {
	Paths  []string `json:"paths" form:"paths" validate:"array_string=1,100,1,200"`
	SiteID *int64   `json:"siteId" form:"siteId" validate:"required,gt=0"`
}

// GetAllPageQuery represents a query to get all pages with pagination
type GetAllPageQuery struct {
	common.PaginationRequestDto
	SiteID *int64 `json:"siteId" form:"siteId" validate:"required,gt=0"`
}

// AdminGetAllPageQuery represents a query for admin to get all pages with pagination
type AdminGetAllPageQuery struct {
	common.PaginationRequestDto
}
