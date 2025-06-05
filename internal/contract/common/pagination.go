package common

type PaginationRequestDto struct {
	Page     int    `form:"page" json:"page" nameFa:"صفحه" validate:"required,min=1"`
	PageSize int    `form:"pageSize" json:"pageSize" nameFa:"تعداد صفحه" validate:"required,min=1,max=100"`
	Search   string `form:"search" json:"search" nameFa:"جستجو" validate:"optional_text=0,100"`
	SearchBy string `form:"searchBy" json:"searchBy" nameFa:"جستجو براساس" validate:"omitempty,oneof=title description"`
	Sort     string `form:"sort" json:"sort" nameFa:"مرتب سازی" validate:"omitempty,oneof=asc desc"`
	SortBy   string `form:"sortBy" json:"sortBy" nameFa:"مرتب سازی براساس" validate:"omitempty,oneof=title description"`
}

type PaginationResponseDto[T any] struct {
	Items           []T   `json:"items" nameFa:"آیتم ها"`
	PageNumber      int   `json:"pageNumber" nameFa:"شماره صفحه"`
	TotalPages      int   `json:"totalPages" nameFa:"تعداد صفحات"`
	TotalCount      int64 `json:"totalCount" nameFa:"تعداد کل"`
	HasPreviousPage bool  `json:"hasPreviousPage" nameFa:"دارای صفحه قبلی"`
	HasNextPage     bool  `json:"hasNextPage" nameFa:"دارای صفحه بعدی"`
}
