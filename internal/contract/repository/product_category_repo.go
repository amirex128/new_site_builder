package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
)

type IProductCategoryRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ProductCategory], error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ProductCategory], error)
	GetAllByParentID(parentID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ProductCategory], error)
	GetByID(id int64) (*domain.ProductCategory, error)
	GetBySlug(slug string) (*domain.ProductCategory, error)
	Create(category *domain.ProductCategory) error
	Update(category *domain.ProductCategory) error
	Delete(id int64) error
}
