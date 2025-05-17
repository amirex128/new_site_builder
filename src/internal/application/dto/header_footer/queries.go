package header_footer

import (
	"github.com/amirex128/new_site_builder/src/internal/application/dto/site"
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// GetByIdHeaderFooterQuery represents a query to get a header/footer by ID
type GetByIdHeaderFooterQuery struct {
	ID     *int64                `json:"id,omitempty" validate:"omitempty,gt=0" error:"gt=شناسه هدر/فوتر باید بزرگتر از 0 باشد"`
	IDs    []int64               `json:"ids,omitempty" validate:"omitempty,dive,gt=0" error:"gt=شناسه‌های هدر/فوتر باید بزرگتر از 0 باشند"`
	SiteID *int64                `json:"siteId" validate:"required,gt=0" error:"required=شناسه سایت الزامی است|gt=شناسه سایت باید بزرگتر از 0 باشد"`
	Type   *HeaderFooterTypeEnum `json:"type,omitempty" validate:"omitempty" error:""`
}

// GetHeaderFooterByDomainOrSiteIdQuery represents a query to get a header/footer by domain or site ID
type GetHeaderFooterByDomainOrSiteIdQuery struct {
	SiteID *int64  `json:"siteId,omitempty" validate:"omitempty,gt=0" error:"gt=شناسه سایت باید بزرگتر از 0 باشد"`
	Domain *string `json:"domain,omitempty" validate:"omitempty,max=200" error:"max=دامنه نمی‌تواند بیشتر از 200 کاراکتر باشد"`
}

// GetAllHeaderFooterQuery represents a query to get all header/footers with pagination
type GetAllHeaderFooterQuery struct {
	common.PaginationRequestDto
	SiteID *int64                `json:"siteId" validate:"required,gt=0" error:"required=شناسه سایت الزامی است|gt=شناسه سایت باید بزرگتر از 0 باشد"`
	Type   *HeaderFooterTypeEnum `json:"type,omitempty" validate:"omitempty" error:""`
}

// AdminGetAllHeaderFooterQuery represents a query for admin to get all header/footers with pagination
type AdminGetAllHeaderFooterQuery struct {
	common.PaginationRequestDto
}
