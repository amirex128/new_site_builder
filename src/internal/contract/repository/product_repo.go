package repository

import (
	common_contract "go-boilerplate/src/internal/contract/common"
	"go-boilerplate/src/internal/domain"
)

type IProductRepository interface {
	GetAll(paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Product, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Product, int64, error)
	GetAllByCategoryID(categoryID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Product, int64, error)
	GetByID(id int64) (domain.Product, error)
	GetBySlug(slug string) (domain.Product, error)
	Create(product domain.Product) error
	Update(product domain.Product) error
	Delete(id int64) error
}
