package website

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
)

// GetArticlesByCategorySlugQuery for retrieving articles by product_category slug
type GetArticlesByCategorySlugQuery struct {
	common.PaginationRequestDto
	Slug   *string `json:"slug" form:"slug" validate:"required,slug" nameFa:"اسلاگ"`
	SiteID *int64  `json:"siteId,omitempty" form:"siteId" validate:"omitempty" nameFa:"شناسه سایت"`
	Domain *string `json:"domain" form:"domain" validate:"required,domain" nameFa:"دامنه"`
}

// GetByDomainHeaderFooterQuery for retrieving header/footer by domain
type GetByDomainHeaderFooterQuery struct {
	Domain *string `json:"domain" form:"domain" validate:"required,domain" nameFa:"دامنه"`
	SiteID *int64  `json:"siteId,omitempty" form:"siteId" validate:"omitempty" nameFa:"شناسه سایت"`
}

// GetByDomainPageQuery for retrieving page by domain and path
type GetByDomainPageQuery struct {
	Domain *string `json:"domain" form:"domain" validate:"required,domain" nameFa:"دامنه"`
	SiteID *int64  `json:"siteId,omitempty" form:"siteId" validate:"omitempty" nameFa:"شناسه سایت"`
	Path   *string `json:"path" form:"path" validate:"required_text=1 200" nameFa:"مسیر"`
}

// GetFiltersSortArticleQuery for retrieving articles with filtering and sorting
type GetFiltersSortArticleQuery struct {
	common.PaginationRequestDto
	SelectedFilters map[enums.ArticleFilterEnum][]string `json:"selectedFilters,omitempty" form:"selectedFilters" validate:"enum_string_map_optional" nameFa:"فیلترهای انتخاب شده"`
	SelectedSort    *enums.ArticleSortEnum               `json:"selectedSort,omitempty" form:"selectedSort" validate:"enum_optional" nameFa:"مرتب سازی انتخاب شده"`
	SiteID          *int64                               `json:"siteId,omitempty" form:"siteId" validate:"omitempty" nameFa:"شناسه سایت"`
	Domain          *string                              `json:"domain" form:"domain" validate:"required,domain" nameFa:"دامنه"`
}

// GetFiltersSortProductQuery for retrieving products with filtering and sorting
type GetFiltersSortProductQuery struct {
	common.PaginationRequestDto
	SelectedFilters map[enums.ProductFilterEnum][]string `json:"selectedFilters,omitempty" form:"selectedFilters" validate:"enum_string_map_optional" nameFa:"فیلترهای انتخاب شده"`
	SelectedSort    *enums.ProductSortEnum               `json:"selectedSort,omitempty" form:"selectedSort" validate:"enum_optional" nameFa:"مرتب سازی انتخاب شده"`
	SiteID          *int64                               `json:"siteId,omitempty" form:"siteId" validate:"omitempty" nameFa:"شناسه سایت"`
	Domain          *string                              `json:"domain" form:"domain" validate:"required,domain" nameFa:"دامنه"`
}

// GetProductsByCategorySlugQuery for retrieving products by product_category slug
type GetProductsByCategorySlugQuery struct {
	common.PaginationRequestDto
	Slug   *string `json:"slug" form:"slug" validate:"required,slug" nameFa:"اسلاگ"`
	SiteID *int64  `json:"siteId,omitempty" form:"siteId" validate:"omitempty" nameFa:"شناسه سایت"`
	Domain *string `json:"domain" form:"domain" validate:"required,domain" nameFa:"دامنه"`
}

// GetSingleArticleBySlugQuery for retrieving a single article by slug
type GetSingleArticleBySlugQuery struct {
	Slug   *string `json:"slug" form:"slug" validate:"required,slug" nameFa:"اسلاگ"`
	SiteID *int64  `json:"siteId,omitempty" form:"siteId" validate:"omitempty" nameFa:"شناسه سایت"`
	Domain *string `json:"domain" form:"domain" validate:"required,domain" nameFa:"دامنه"`
}

// GetSingleProductBySlugQuery for retrieving a single article by slug
type GetSingleProductBySlugQuery struct {
	Slug   *string `json:"slug" form:"slug" validate:"required,slug" nameFa:"اسلاگ"`
	SiteID *int64  `json:"siteId,omitempty" form:"siteId" validate:"omitempty" nameFa:"شناسه سایت"`
	Domain *string `json:"domain" form:"domain" validate:"required,domain" nameFa:"دامنه"`
}

// ProductSearchListQuery for article search listing
type ProductSearchListQuery struct {
	common.PaginationRequestDto
	Domain *string `json:"domain" form:"domain" validate:"required,domain" nameFa:"دامنه"`
	SiteID *int64  `json:"siteId,omitempty" form:"siteId" validate:"omitempty" nameFa:"شناسه سایت"`
}
