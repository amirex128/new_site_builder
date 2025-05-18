package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IArticleCategoryRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.BlogCategory, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.BlogCategory, int64, error)
	GetAllByParentID(parentID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.BlogCategory, int64, error)
	GetByID(id int64) (domain.BlogCategory, error)
	GetBySlug(slug string) (domain.BlogCategory, error)
	GetBySlugAndSiteID(slug string, siteID int64) (domain.BlogCategory, error)
	Create(category domain.BlogCategory) error
	Update(category domain.BlogCategory) error
	Delete(id int64) error

	// Media relationship methods
	GetCategoryMedia(categoryID int64) ([]domain.Media, error)
	AddMediaToCategory(categoryID int64, mediaID int64) error
	RemoveMediaFromCategory(categoryID int64, mediaID int64) error
	RemoveAllMediaFromCategory(categoryID int64) error
}
