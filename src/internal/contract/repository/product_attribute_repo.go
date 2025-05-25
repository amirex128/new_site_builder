package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IProductAttributeRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ProductAttribute], error)
	GetAllByProductID(productID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ProductAttribute], error)
	GetByID(id int64) (*domain.ProductAttribute, error)
	Create(attribute *domain.ProductAttribute) error
	Update(attribute *domain.ProductAttribute) error
	Delete(id int64) error
}
