package common

type PaginationRequestDto struct {
	Page     int    `json:"page" validate:"required,min=1"`
	PageSize int    `json:"pageSize" validate:"required,min=1,max=100"`
	Search   string `json:"search" validate:"omitempty,max=100"`
	SearchBy string `json:"searchBy" validate:"omitempty,oneof=title description"`
	Sort     string `json:"sort" validate:"omitempty,oneof=asc desc"`
	SortBy   string `json:"sortBy" validate:"omitempty,oneof=title description"`
}
