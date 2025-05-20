package website

import (
	"github.com/amirex128/new_site_builder/src/internal/application/dto/product"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// GetArticlesByCategorySlugQuery for retrieving articles by product_category slug
type GetArticlesByCategorySlugQuery struct {
	common.PaginationRequestDto
	Slug   *string `json:"slug" form:"slug" validate:"required,slug"`
	SiteID *int64  `json:"siteId,omitempty" form:"siteId" validate:"omitempty"`
	Domain *string `json:"domain" form:"domain" validate:"required,domain"`
}

// GetByDomainHeaderFooterQuery for retrieving header/footer by domain
type GetByDomainHeaderFooterQuery struct {
	Domain *string `json:"domain" form:"domain" validate:"required,domain"`
	SiteID *int64  `json:"siteId,omitempty" form:"siteId" validate:"omitempty"`
}

// GetByDomainPageQuery for retrieving page by domain and path
type GetByDomainPageQuery struct {
	Domain *string `json:"domain" form:"domain" validate:"required,domain"`
	SiteID *int64  `json:"siteId,omitempty" form:"siteId" validate:"omitempty"`
	Path   *string `json:"path" form:"path" validate:"required_text=1,200"`
}

// GetFiltersSortArticleQuery for retrieving articles with filtering and sorting
type GetFiltersSortArticleQuery struct {
	common.PaginationRequestDto
	SelectedFilters map[ArticleFilterEnum][]string `json:"selectedFilters,omitempty" form:"selectedFilters" validate:"enum_string_map_optional"`
	SelectedSort    *ArticleSortEnum               `json:"selectedSort,omitempty" form:"selectedSort" validate:"enum_optional"`
	SiteID          *int64                         `json:"siteId,omitempty" form:"siteId" validate:"omitempty"`
	Domain          *string                        `json:"domain" form:"domain" validate:"required,domain"`
}

// GetFiltersSortProductQuery for retrieving products with filtering and sorting
type GetFiltersSortProductQuery struct {
	common.PaginationRequestDto
	SelectedFilters map[product.ProductFilterEnum][]string `json:"selectedFilters,omitempty" form:"selectedFilters" validate:"enum_string_map_optional"`
	SelectedSort    *product.ProductSortEnum               `json:"selectedSort,omitempty" form:"selectedSort" validate:"enum_optional"`
	SiteID          *int64                                 `json:"siteId,omitempty" form:"siteId" validate:"omitempty"`
	Domain          *string                                `json:"domain" form:"domain" validate:"required,domain"`
}

// GetProductsByCategorySlugQuery for retrieving products by product_category slug
type GetProductsByCategorySlugQuery struct {
	common.PaginationRequestDto
	Slug   *string `json:"slug" form:"slug" validate:"required,slug"`
	SiteID *int64  `json:"siteId,omitempty" form:"siteId" validate:"omitempty"`
	Domain *string `json:"domain" form:"domain" validate:"required,domain"`
}

// GetSingleArticleBySlugQuery for retrieving a single article by slug
type GetSingleArticleBySlugQuery struct {
	Slug   *string `json:"slug" form:"slug" validate:"required,slug"`
	SiteID *int64  `json:"siteId,omitempty" form:"siteId" validate:"omitempty"`
	Domain *string `json:"domain" form:"domain" validate:"required,domain"`
}

// GetSingleProductBySlugQuery for retrieving a single article by slug
type GetSingleProductBySlugQuery struct {
	Slug   *string `json:"slug" form:"slug" validate:"required,slug"`
	SiteID *int64  `json:"siteId,omitempty" form:"siteId" validate:"omitempty"`
	Domain *string `json:"domain" form:"domain" validate:"required,domain"`
}

// ProductSearchListQuery for article search listing
type ProductSearchListQuery struct {
	common.PaginationRequestDto
	Domain *string `json:"domain" form:"domain" validate:"required,domain"`
	SiteID *int64  `json:"siteId,omitempty" form:"siteId" validate:"omitempty"`
}
