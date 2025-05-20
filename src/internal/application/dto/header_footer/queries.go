package header_footer

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// GetByIdHeaderFooterQuery represents a query to get a header/footer by ID
type GetByIdHeaderFooterQuery struct {
	ID     *int64                `json:"id,omitempty" form:"id" validate:"omitempty,gt=0"`
	IDs    []int64               `json:"ids,omitempty" form:"ids" validate:"array_number_optional=0,100,1,0,false"`
	SiteID *int64                `json:"siteId" form:"siteId" validate:"required,gt=0"`
	Type   *HeaderFooterTypeEnum `json:"type,omitempty" form:"type" validate:"enum_optional"`
}

// GetHeaderFooterByDomainOrSiteIdQuery represents a query to get a header/footer by domain or site ID
type GetHeaderFooterByDomainOrSiteIdQuery struct {
	SiteID *int64  `json:"siteId,omitempty" validate:"omitempty,gt=0"`
	Domain *string `json:"domain,omitempty" validate:"domain_optional"`
}

// GetAllHeaderFooterQuery represents a query to get all header/footers with pagination
type GetAllHeaderFooterQuery struct {
	common.PaginationRequestDto
	SiteID *int64                `json:"siteId" validate:"required,gt=0"`
	Type   *HeaderFooterTypeEnum `json:"type,omitempty" validate:"enum_optional"`
}

// AdminGetAllHeaderFooterQuery represents a query for admin to get all header/footers with pagination
type AdminGetAllHeaderFooterQuery struct {
	common.PaginationRequestDto
}
