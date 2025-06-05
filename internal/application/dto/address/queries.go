package address

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
)

// GetByIdAddressQuery represents a query to get an address by ID
type GetByIdAddressQuery struct {
	ID *int64 `json:"id" nameFa:"شناسه" form:"id" validate:"required,gt=0"`
}

// GetAllAddressQuery represents a query to get all addresses with pagination
type GetAllAddressQuery struct {
}

// AdminGetAllAddressQuery represents a query for admin to get all addresses with pagination
type AdminGetAllAddressQuery struct {
	common.PaginationRequestDto
}

// GetAllCityQuery represents a query to get all cities with pagination
type GetAllCityQuery struct {
	common.PaginationRequestDto
}

// GetAllProvinceQuery represents a query to get all provinces with pagination
type GetAllProvinceQuery struct {
	common.PaginationRequestDto
}
