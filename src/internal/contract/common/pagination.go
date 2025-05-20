package common

type PaginationRequestDto struct {
	Page     int    `form:"page" json:"page" validate:"required,min=1"`
	PageSize int    `form:"pageSize" json:"pageSize" validate:"required,min=1,max=100"`
	Search   string `form:"search" json:"search" validate:"optional_text=0,100"`
	SearchBy string `form:"searchBy" json:"searchBy" validate:"omitempty,oneof=title description"`
	Sort     string `form:"sort" json:"sort" validate:"omitempty,oneof=asc desc"`
	SortBy   string `form:"sortBy" json:"sortBy" validate:"omitempty,oneof=title description"`
}
