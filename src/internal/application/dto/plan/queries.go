package plan

import "github.com/amirex128/new_site_builder/src/internal/contract/common"

// GetAllPlanQuery represents a query to get all plans with pagination
type GetAllPlanQuery struct {
	common.PaginationRequestDto
}

// GetByIdPlanQuery represents a query to get a plan by ID
type GetByIdPlanQuery struct {
	ID *int64 `json:"id" validate:"required,gt=0" error:"required=شناسه الزامی است|gt=شناسه باید بزرگتر از 0 باشد"`
}
