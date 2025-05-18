package article

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// ArticleSortEnum matches the .NET ArticleSortEnum
type ArticleSortEnum int

const (
	TitleAZ ArticleSortEnum = iota
	TitleZA
	RecentlyAdded
	RecentlyUpdated
	MostVisited
	LeastVisited
	MostRated
	LeastRated
	MostReviewed
	LeastReviewed
)

// ArticleFilterEnum matches the .NET ArticleFilterEnum
type ArticleFilterEnum int

const (
	RateRange ArticleFilterEnum = iota
	ReviewRange
	VisitedRange
	AddedRange
	UpdatedRange
	CategoryIds
	ArticleIds
	Badges
)

// GetByIdArticleQuery for retrieving a single article by ID
type GetByIdArticleQuery struct {
	ID *int64 `json:"id" validate:"required"`
}

// GetSingleArticleQuery for retrieving a single article by slug
type GetSingleArticleQuery struct {
	Slug   *string `json:"slug" validate:"slug"`
	SiteID *int64  `json:"siteId" validate:"required"`
}

// AdminGetAllArticleQuery for admin listing of all articles with pagination
type AdminGetAllArticleQuery struct {
	common.PaginationRequestDto
}

// GetAllArticleQuery for listing articles by site ID with pagination
type GetAllArticleQuery struct {
	common.PaginationRequestDto
	SiteID *int64 `json:"siteId" validate:"required"`
}

// GetArticleByCategoryQuery for retrieving articles by product_category with pagination
type GetArticleByCategoryQuery struct {
	common.PaginationRequestDto
	Slug   *string `json:"slug" validate:"slug"`
	SiteID *int64  `json:"siteId" validate:"required"`
}

// GetByFiltersSortArticleQuery for retrieving articles with filtering and sorting
type GetByFiltersSortArticleQuery struct {
	common.PaginationRequestDto
	SelectedFilters map[ArticleFilterEnum][]string `json:"selectedFilters,omitempty" validate:"enum_string_map_optional"`
	SelectedSort    *ArticleSortEnum               `json:"selectedSort,omitempty" validate:"enum_optional"`
	SiteID          *int64                         `json:"siteId" validate:"required"`
}
