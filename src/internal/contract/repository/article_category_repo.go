package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IArticleCategoryRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ArticleCategory], error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ArticleCategory], error)
	GetAllByParentID(parentID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ArticleCategory], error)
	GetByID(id int64) (domain.ArticleCategory, error)
	GetBySlug(slug string) (domain.ArticleCategory, error)
	GetBySlugAndSiteID(slug string, siteID int64) (domain.ArticleCategory, error)
	Create(category domain.ArticleCategory) error
	Update(category domain.ArticleCategory) error
	Delete(id int64) error

	// Media relationship methods
	GetCategoryMedia(categoryID int64) ([]domain.Media, error)
	AddMediaToCategory(categoryID int64, mediaID int64) error
	RemoveMediaFromCategory(categoryID int64, mediaID int64) error
	RemoveAllMediaFromCategory(categoryID int64) error
}
