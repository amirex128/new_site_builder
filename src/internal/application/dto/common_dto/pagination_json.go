package common_dto

type PaginationJson struct {
	Items           []interface{} `json:"items" nameFa:"آیتم ها"`
	PageNumber      int           `json:"pageNumber" nameFa:"شماره صفحه"`
	TotalPages      int           `json:"totalPages" nameFa:"تعداد صفحات"`
	TotalCount      int64         `json:"totalCount" nameFa:"تعداد کل"`
	HasPreviousPage bool          `json:"hasPreviousPage" nameFa:"دارای صفحه قبلی"`
	HasNextPage     bool          `json:"hasNextPage" nameFa:"دارای صفحه بعدی"`
	Search          string        `form:"search" json:"search" nameFa:"جستجو" validate:"optional_text=0 100"`
}
