package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IArticleRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Article, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Article, int64, error)
	GetAllByCategoryID(categoryID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Article, int64, error)
	GetByID(id int64) (domain.Article, error)
	GetBySlug(slug string) (domain.Article, error)
	Create(article domain.Article) error
	Update(article domain.Article) error
	Delete(id int64) error
}
