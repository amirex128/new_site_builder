package repository

import (
	common_contract "go-boilerplate/src/internal/contract/common"
	"go-boilerplate/src/internal/domain"
)

type IProductCategoryRepository interface {
	GetAll(paginationRequestDto common_contract.PaginationRequestDto) ([]domain.ProductCategory, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.ProductCategory, int64, error)
	GetAllByParentID(parentID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.ProductCategory, int64, error)
	GetByID(id int64) (domain.ProductCategory, error)
	GetBySlug(slug string) (domain.ProductCategory, error)
	Create(category domain.ProductCategory) error
	Update(category domain.ProductCategory) error
	Delete(id int64) error
}
