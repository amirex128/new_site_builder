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
	Create(category domain.BlogCategory) error
	Update(category domain.BlogCategory) error
	Delete(id int64) error
}
