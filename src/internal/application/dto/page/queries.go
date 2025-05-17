package page

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// GetByIdPageQuery represents a query to get page(s) by ID
type GetByIdPageQuery struct {
	ID     *int64  `json:"id,omitempty" validate:"omitempty,gt=0" error:"gt=شناسه صفحه باید بزرگتر از 0 باشد"`
	IDs    []int64 `json:"ids,omitempty" validate:"omitempty,dive,gt=0" error:"gt=شناسه‌های صفحه باید بزرگتر از 0 باشند"`
	SiteID *int64  `json:"siteId" validate:"required,gt=0" error:"required=شناسه سایت الزامی است|gt=شناسه سایت باید بزرگتر از 0 باشد"`
}

// GetByPathPageQuery represents a query to get page(s) by path
type GetByPathPageQuery struct {
	Paths  []string `json:"paths" validate:"required,dive,pattern=^[a-z0-9-]+$" error:"required=مسیرها الزامی هستند|pattern=مسیر باید یک اسلاگ معتبر باشد"`
	SiteID *int64   `json:"siteId" validate:"required,gt=0" error:"required=شناسه سایت الزامی است|gt=شناسه سایت باید بزرگتر از 0 باشد"`
}

// GetAllPageQuery represents a query to get all pages with pagination
type GetAllPageQuery struct {
	common.PaginationRequestDto
	SiteID *int64 `json:"siteId" validate:"required,gt=0" error:"required=شناسه سایت الزامی است|gt=شناسه سایت باید بزرگتر از 0 باشد"`
}

// AdminGetAllPageQuery represents a query for admin to get all pages with pagination
type AdminGetAllPageQuery struct {
	common.PaginationRequestDto
}
