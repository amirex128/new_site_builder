package article

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// GetByIdArticleQuery for retrieving a single article by ID
type GetByIdArticleQuery struct {
	ID *int64 `json:"id" form:"id" validate:"required"`
}

// GetSingleArticleQuery for retrieving a single article by slug
type GetSingleArticleQuery struct {
	Slug   *string `json:"slug" form:"slug" validate:"slug"`
	SiteID *int64  `json:"siteId" form:"siteId" validate:"required"`
}

// AdminGetAllArticleQuery for admin listing of all articles with pagination
type AdminGetAllArticleQuery struct {
	common.PaginationRequestDto
}

// GetAllArticleQuery for listing articles by site ID with pagination
type GetAllArticleQuery struct {
	common.PaginationRequestDto
	SiteID *int64 `json:"siteId" form:"siteId" validate:"required"`
}

// GetArticleByCategoryQuery for retrieving articles by product_category with pagination
type GetArticleByCategoryQuery struct {
	common.PaginationRequestDto
	Slug   *string `json:"slug" form:"slug" validate:"slug"`
	SiteID *int64  `json:"siteId" form:"siteId" validate:"required"`
}

// GetByFiltersSortArticleQuery for retrieving articles with filtering and sorting
type GetByFiltersSortArticleQuery struct {
	common.PaginationRequestDto
	SelectedFilters map[ArticleFilterEnum][]string `json:"selectedFilters,omitempty" form:"selectedFilters" validate:"enum_string_map_optional"`
	SelectedSort    *ArticleSortEnum               `json:"selectedSort,omitempty" form:"selectedSort" validate:"enum_optional"`
	SiteID          *int64                         `json:"siteId" form:"siteId" validate:"required"`
}
