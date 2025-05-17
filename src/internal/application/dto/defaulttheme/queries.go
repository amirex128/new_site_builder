package defaulttheme

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// GetByIdDefaultThemeQuery represents a query to get a default theme by ID
type GetByIdDefaultThemeQuery struct {
	ID *int64 `json:"id" validate:"required,gt=0" error:"required=شناسه الزامی است|gt=شناسه باید بزرگتر از 0 باشد"`
}

// GetAllDefaultThemeQuery represents a query to get all default themes with pagination
type GetAllDefaultThemeQuery struct {
	common.PaginationRequestDto
}
