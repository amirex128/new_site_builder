package defaulttheme

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
)

// GetByIdDefaultThemeQuery represents a query to get a default theme by ID
type GetByIdDefaultThemeQuery struct {
	ID *int64 `json:"id" form:"id" nameFa:"شناسه" validate:"required,gt=0"`
}

// GetAllDefaultThemeQuery represents a query to get all default themes with pagination
type GetAllDefaultThemeQuery struct {
	common.PaginationRequestDto
}
