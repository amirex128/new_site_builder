package site

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// GetByIdSiteQuery represents a query to get a site by ID
type GetByIdSiteQuery struct {
	ID *int64 `json:"id" validate:"required,gt=0"`
}

// GetByDomainSiteQuery represents a query to get a site by domain
type GetByDomainSiteQuery struct {
	Domain *string `json:"domain" validate:"required,domain"`
}

// GetAllSiteQuery represents a query to get all sites with pagination
type GetAllSiteQuery struct {
	common.PaginationRequestDto
}

// AdminGetAllSiteQuery represents a query for admin to get all sites with pagination
type AdminGetAllSiteQuery struct {
	common.PaginationRequestDto
}
