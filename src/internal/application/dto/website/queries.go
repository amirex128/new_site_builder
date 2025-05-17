package website

import (
	"github.com/amirex128/new_site_builder/src/internal/application/dto/product"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// GetArticlesByCategorySlugQuery for retrieving articles by product_category slug
type GetArticlesByCategorySlugQuery struct {
	common.PaginationRequestDto
	Slug   *string `json:"slug" validate:"required" error:"required=نامک الزامی است"`
	SiteID *int64  `json:"siteId,omitempty" validate:"omitempty" error:""`
	Domain *string `json:"domain" validate:"required,max=200" error:"required=دامنه الزامی است|max=دامنه نمی‌تواند بیشتر از 200 کاراکتر باشد"`
}

// GetByDomainHeaderFooterQuery for retrieving header/footer by domain
type GetByDomainHeaderFooterQuery struct {
	Domain *string `json:"domain" validate:"required,max=200" error:"required=دامنه الزامی است|max=دامنه نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	SiteID *int64  `json:"siteId,omitempty" validate:"omitempty" error:""`
}

// GetByDomainPageQuery for retrieving page by domain and path
type GetByDomainPageQuery struct {
	Domain *string `json:"domain" validate:"required,max=200" error:"required=دامنه الزامی است|max=دامنه نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	SiteID *int64  `json:"siteId,omitempty" validate:"omitempty" error:""`
	Path   *string `json:"path" validate:"required,pattern=^[a-zA-Z0-9\\-._~!$&'()*+,;=:@%/]*$" error:"required=آدرس الزامی است|pattern=آدرس نامعتبر است"`
}

// GetFiltersSortArticleQuery for retrieving articles with filtering and sorting
type GetFiltersSortArticleQuery struct {
	common.PaginationRequestDto
	SelectedFilters map[ArticleFilterEnum][]string `json:"selectedFilters,omitempty" validate:"omitempty" error:""`
	SelectedSort    *ArticleSortEnum               `json:"selectedSort,omitempty" validate:"omitempty" error:""`
	SiteID          *int64                         `json:"siteId,omitempty" validate:"omitempty" error:""`
	Domain          *string                        `json:"domain" validate:"required,max=200" error:"required=دامنه الزامی است|max=دامنه نمی‌تواند بیشتر از 200 کاراکتر باشد"`
}

// GetFiltersSortProductQuery for retrieving products with filtering and sorting
type GetFiltersSortProductQuery struct {
	common.PaginationRequestDto
	SelectedFilters map[product.ProductFilterEnum][]string `json:"selectedFilters,omitempty" validate:"omitempty" error:""`
	SelectedSort    *product.ProductSortEnum               `json:"selectedSort,omitempty" validate:"omitempty" error:""`
	SiteID          *int64                                 `json:"siteId,omitempty" validate:"omitempty" error:""`
	Domain          *string                                `json:"domain" validate:"required,max=200" error:"required=دامنه الزامی است|max=دامنه نمی‌تواند بیشتر از 200 کاراکتر باشد"`
}

// GetProductsByCategorySlugQuery for retrieving products by product_category slug
type GetProductsByCategorySlugQuery struct {
	common.PaginationRequestDto
	Slug   *string `json:"slug" validate:"required" error:"required=نامک الزامی است"`
	SiteID *int64  `json:"siteId,omitempty" validate:"omitempty" error:""`
	Domain *string `json:"domain" validate:"required,max=200" error:"required=دامنه الزامی است|max=دامنه نمی‌تواند بیشتر از 200 کاراکتر باشد"`
}

// GetSingleArticleBySlugQuery for retrieving a single article by slug
type GetSingleArticleBySlugQuery struct {
	Slug   *string `json:"slug" validate:"required" error:"required=نامک الزامی است"`
	SiteID *int64  `json:"siteId,omitempty" validate:"omitempty" error:""`
	Domain *string `json:"domain" validate:"required,max=200" error:"required=دامنه الزامی است|max=دامنه نمی‌تواند بیشتر از 200 کاراکتر باشد"`
}

// GetSingleProductBySlugQuery for retrieving a single product by slug
type GetSingleProductBySlugQuery struct {
	Slug   *string `json:"slug" validate:"required" error:"required=نامک الزامی است"`
	SiteID *int64  `json:"siteId,omitempty" validate:"omitempty" error:""`
	Domain *string `json:"domain" validate:"required,max=200" error:"required=دامنه الزامی است|max=دامنه نمی‌تواند بیشتر از 200 کاراکتر باشد"`
}

// ProductSearchListQuery for product search listing
type ProductSearchListQuery struct {
	common.PaginationRequestDto
	Domain *string `json:"domain" validate:"required,max=200" error:"required=دامنه الزامی است|max=دامنه نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	SiteID *int64  `json:"siteId,omitempty" validate:"omitempty" error:""`
}
