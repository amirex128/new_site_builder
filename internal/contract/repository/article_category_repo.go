package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	domain2 "github.com/amirex128/new_site_builder/internal/domain"
)

type IArticleCategoryRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain2.ArticleCategory], error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain2.ArticleCategory], error)
	GetAllByParentID(parentID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain2.ArticleCategory], error)
	GetByID(id int64) (*domain2.ArticleCategory, error)
	GetBySlug(slug string) (*domain2.ArticleCategory, error)
	GetBySlugAndSiteID(slug string, siteID int64) (*domain2.ArticleCategory, error)
	Create(category *domain2.ArticleCategory) error
	Update(category *domain2.ArticleCategory) error
	Delete(id int64) error

	// Media relationship methods
	GetCategoryMedia(categoryID int64) ([]domain2.Media, error)
	AddMediaToCategory(categoryID int64, mediaID int64) error
	RemoveMediaFromCategory(categoryID int64, mediaID int64) error
	RemoveAllMediaFromCategory(categoryID int64) error
}
