package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
)

type IProductRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (common.PaginationResponseDto[domain.Product], int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (common.PaginationResponseDto[domain.Product], int64, error)
	GetAllByCategoryID(categoryID int64, paginationRequestDto common.PaginationRequestDto) (common.PaginationResponseDto[domain.Product], int64, error)
	GetByID(id int64) (domain.Product, error)
	GetBySlug(slug string) (domain.Product, error)
	Create(product domain.Product) error
	Update(product domain.Product) error
	Delete(id int64) error
	GetAllByFilterAndSort(siteID int64, filters map[enums.ProductFilterEnum][]string, sort *string, paginationRequestDto common.PaginationRequestDto) (common.PaginationResponseDto[domain.Product], int64, error)
}
