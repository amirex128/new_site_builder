package repository

import (
	common_contract "go-boilerplate/src/internal/contract/common"
	"go-boilerplate/src/internal/domain"
)

type IBlogCategoryRepository interface {
	GetAll(paginationRequestDto common_contract.PaginationRequestDto) ([]domain.BlogCategory, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.BlogCategory, int64, error)
	GetAllByParentID(parentID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.BlogCategory, int64, error)
	GetByID(id int64) (domain.BlogCategory, error)
	GetBySlug(slug string) (domain.BlogCategory, error)
	Create(category domain.BlogCategory) error
	Update(category domain.BlogCategory) error
	Delete(id int64) error
}
