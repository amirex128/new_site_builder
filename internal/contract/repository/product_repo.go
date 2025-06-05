package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
	"github.com/amirex128/new_site_builder/internal/domain/enums"
)

type IProductRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Product], error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Product], error)
	GetAllByCategoryID(categoryID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Product], error)
	GetByID(id int64) (*domain.Product, error)
	GetBySlug(slug string) (*domain.Product, error)
	Create(product *domain.Product) error
	Update(product *domain.Product) error
	Delete(id int64) error
	GetAllByFilterAndSort(siteID int64, filters map[enums.ProductFilterEnum][]string, sort *string, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Product], error)
}
