package address

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// GetByIdAddressQuery represents a query to get an address by ID
type GetByIdAddressQuery struct {
	ID *int64 `json:"id" validate:"required,gt=0" error:"required=شناسه آدرس الزامی است|gt=شناسه آدرس باید بزرگتر از 0 باشد"`
}

// GetAllAddressQuery represents a query to get all addresses with pagination
type GetAllAddressQuery struct {
	common.PaginationRequestDto
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
