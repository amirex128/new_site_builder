package common_dto

type PaginationJson struct {
	Items           []interface{} `json:"items"`
	PageNumber      int           `json:"pageNumber"`
	TotalPages      int           `json:"totalPages"`
	TotalCount      int64         `json:"totalCount"`
	HasPreviousPage bool          `json:"hasPreviousPage"`
	HasNextPage     bool          `json:"hasNextPage"`
}
