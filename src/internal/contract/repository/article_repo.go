package repository

import (
	common_contract "go-boilerplate/src/internal/contract/common"
	"go-boilerplate/src/internal/domain"
)

type IArticleRepository interface {
	GetAll(paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Article, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Article, int64, error)
	GetAllByCategoryID(categoryID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Article, int64, error)
	GetByID(id int64) (domain.Article, error)
	GetBySlug(slug string) (domain.Article, error)
	Create(article domain.Article) error
	Update(article domain.Article) error
	Delete(id int64) error
}
